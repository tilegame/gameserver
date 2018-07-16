package combat

import (
	"fmt"
	"testing"
)

func TestNewInstance(t *testing.T) {
	x := NewInstance()
	fmt.Println(x)
}

func TestMoveBuffer(t *testing.T) {
	a := NewActor()
	t.Log(a)
	a.buffer.spots[5] = 1
	a.Buffer().ExtendDirectly(5, 30)
	t.Log(a.buffer.spots)
}
