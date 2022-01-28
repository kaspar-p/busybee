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

	type TestCase struct {
		Slice []string
		Elem  string
		Exp   []string
	}

	tests := []TestCase{
		{[]string{"a", "b", "c"}, "a", []string{"b", "c"}},
		{[]string{"a", "b", "c"}, "d", []string{"a", "b", "c"}},
		{[]string{"a"}, "a", []string{}},
		{[]string{"a", "b", "a", "a"}, "a", []string{"b", "a", "a"}},
	}

	for _, oneTest := range tests {
		actual := RemoveStringFromSlice(oneTest.Slice, oneTest.Elem)
		test.Equals(t, actual, oneTest.Exp)
	}
}

func TestStringInSlice(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		Slice    []string
		Elem     string
		ExpValue bool
		ExpIndex int
	}

	tests := []TestCase{
		{[]string{"a", "b", "c"}, "a", true, 0},
		{[]string{"a", "b", "c"}, "d", false, 0},
		{[]string{"a"}, "a", true, 0},
		{[]string{"c", "b", "a", "a"}, "a", true, 2},
	}

	for _, oneTest := range tests {
		actualValue, actualIndex := StringInSlice(oneTest.Slice, oneTest.Elem)
		test.Equals(t, actualValue, oneTest.ExpValue)
		test.Equals(t, actualIndex, oneTest.ExpIndex)
	}
}
