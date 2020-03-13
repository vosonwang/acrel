package acrel

import (
	"encoding/binary"
	"fmt"
	"github.com/vosonwang/libcrc"
)

/*
安科瑞 Modbus对接协议
*/

/*基本格式
{{ 命令字（1字节）消息体（可变）校验位（2字节） }}
7b 7b x xxx xx  7d 7d

校验位长度为 2 个字节，Modbus CRC 校验算法 。校验范围为命令字开始（含命令字）到消息体结束。校验位采用小端模式。

注：任何与服务器进行交互的数据都需要按此格式进行编解码。（包括透传，注册等等）
*/

var (
	startDelimiters = byte(0x7b)
	endDelimiters   = byte(0x7d)
)

type Frame struct {
	Function uint8
	Data     []byte
}

// NewFrame converts a packet to a Acrel frame.
func NewFrame(packet []byte) (*Frame, error) {
	// Check the that the packet length.
	if len(packet) < 7 {
		return nil, fmt.Errorf("frame error: packet less than 7 bytes: 0x% x", packet)
	}

	pLen := len(packet)

	// 检查 Delimiters 定界符
	if packet[0]&startDelimiters != startDelimiters || packet[1] != packet[0] ||
		packet[pLen-2]&endDelimiters != endDelimiters || packet[pLen-1] != packet[pLen-2] {
		return nil, fmt.Errorf("定界符错误：0x% x", packet)
	}

	// Check the CRC.
	crcExpect := binary.LittleEndian.Uint16(packet[pLen-4 : pLen-2])
	crcCalc := libcrc.CRCModbus(packet[2 : pLen-4])

	if crcCalc != crcExpect {
		return nil, fmt.Errorf("frame error: CRC (expected 0x%x, got 0x%x)", crcExpect, crcCalc)
	}

	frame := &Frame{
		Function: packet[2],
		Data:     packet[3 : pLen-4],
	}

	return frame, nil
}

func (frame *Frame) Copy() *Frame {
	f := *frame
	return &f
}

// Bytes returns the MODBUS byte stream based on the Frame fields
func (frame *Frame) Bytes() []byte {
	bytes := make([]byte, 3)

	// 添加定界符
	bytes[0] = startDelimiters
	bytes[1] = startDelimiters
	bytes[2] = frame.Function

	bytes = append(bytes, frame.Data...)

	// Calculate the CRC.
	pLen := len(bytes)
	crc := libcrc.CRCModbus(bytes[2:pLen])

	// Add the CRC.
	bytes = append(bytes, []byte{0, 0}...)
	binary.LittleEndian.PutUint16(bytes[pLen:pLen+2], crc)

	bytes = append(bytes, endDelimiters, endDelimiters)
	return bytes
}

// GetFunction returns the Modbus function code.
func (frame *Frame) GetFunction() uint8 {
	return frame.Function
}

// GetData returns the Frame Data byte field.
func (frame *Frame) GetData() []byte {
	return frame.Data
}

// SetData sets the Frame Data byte field and updates the frame length
// accordingly.
func (frame *Frame) SetData(data []byte) {
	frame.Data = data
}
