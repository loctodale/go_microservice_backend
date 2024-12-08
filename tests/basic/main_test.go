package basic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddOne(t *testing.T) {
	//var (
	//	input  = 1
	//	output = 2
	//)
	//
	//actual := AddOne(input)
	//if actual != output {
	//	t.Error("Expected", output, "Got", actual)
	//}
	assert.Equal(t, AddOne(1), 2, "fail")
}
