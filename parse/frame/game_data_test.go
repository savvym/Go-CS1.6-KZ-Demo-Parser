package frame

import (
	"fmt"
	"testing"
)

func TestGameDataFrame(t *testing.T) {
	s := message[10]
	if len(s) == 0 {
		fmt.Println("no such message")
	}
	fmt.Println(s)
}
