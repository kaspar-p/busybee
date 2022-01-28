package utils_test

import (
	"testing"

	"github.com/kaspar-p/bee/src/test"
	. "github.com/kaspar-p/bee/src/utils"
)

func TestWrapStringInCodeBlock(t *testing.T) {
	t.Parallel()

	s := "some string"
	actual := WrapStringInCodeBlock(s)
	test.Equals(t, "```some string```", actual)
}

func TestRemoveStringFromSlice(t *testing.T) {
	t.Parallel()

	type Input struct {
		Slice []string
		Elem  string
		Exp   []string
	}

	tests := []Input{
		{Slice: []string{"a", "b", "c"}, Elem: "a", Exp: []string{"b", "c"}},
		{Slice: []string{"a", "b", "c"}, Elem: "d", Exp: []string{"a", "b", "c"}},
		{Slice: []string{"a"}, Elem: "a", Exp: []string{}},
		{Slice: []string{"a", "b", "a", "a"}, Elem: "a", Exp: []string{"b", "a", "a"}},
	}

	for _, oneTest := range tests {
		actual := RemoveStringFromSlice(oneTest.Slice, oneTest.Elem)
		test.Equals(t, actual, oneTest.Exp)
	}
}

func TestStringInSlice(t *testing.T) {
	t.Parallel()

	type Input struct {
		Slice    []string
		Elem     string
		ExpValue bool
		ExpIndex int
	}

	tests := []Input{
		{Slice: []string{"a", "b", "c"}, Elem: "a", ExpValue: true, ExpIndex: 0},
		{Slice: []string{"a", "b", "c"}, Elem: "d", ExpValue: false, ExpIndex: 0},
		{Slice: []string{"a"}, Elem: "a", ExpValue: true, ExpIndex: 0},
		{Slice: []string{"c", "b", "a", "a"}, Elem: "a", ExpValue: true, ExpIndex: 2},
	}

	for _, oneTest := range tests {
		actualValue, actualIndex := StringInSlice(oneTest.Slice, oneTest.Elem)
		test.Equals(t, actualValue, oneTest.ExpValue)
		test.Equals(t, actualIndex, oneTest.ExpIndex)
	}
}
