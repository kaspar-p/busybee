package ingest_test

import (
	"testing"

	. "github.com/kaspar-p/busybee/src/ingest"
	"github.com/kaspar-p/busybee/src/test"
)

func TestParseEventTitle(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		String   string
		Expected string
	}

	tests := []TestCase{
		{String: "MAT223H1 LEC0101", Expected: "MAT223"},
		{String: "MAT223H3 LEC0101", Expected: "MAT223"},
		{String: "MAT223H5 LEC0101", Expected: "MAT223"},
		{String: "LIN200Y1 LEC0502", Expected: "LIN200"},
		{String: "LIN200Y3 LEC0502", Expected: "LIN200"},
		{String: "LIN200Y5 LEC0502", Expected: "LIN200"},
		{String: "something else", Expected: "something else"},
	}

	for _, testCase := range tests {
		actual := ParseEventTitle(testCase.String)
		test.Equals(t, testCase.Expected, actual)
	}
}
