package deepcopy

import (
	"fmt"
	"testing"
)

type TestS struct {
	Dict  map[int]int `json:"-"`
	Array []int
	Value int
}

func TestCopyJsonObject(t *testing.T) {
	x := TestS{
		Dict:  map[int]int{1: 2, 3: 4},
		Array: []int{1, 2, 3, 4},
		Value: 10,
	}
	y := CopyJsonObject(x)
	fmt.Printf("%+v", y)
}
