package acrel

import (
	"fmt"
	"testing"
)

func TestFrame_Copy(t *testing.T) {

	a := &Frame{
		Function: 0x84,
		Data:     nil,
	}
	fmt.Println(a)
	b := a.Copy()
	fmt.Println(b)

	b.Function = 0x94
	// a、b是独立的两块内存
	fmt.Println(b)
	fmt.Println(a)
}
