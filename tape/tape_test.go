package tape

import (
	"testing"
	"time"
)

func TestMins(t *testing.T) {
	type test struct {
		input       time.Duration
		want        string
		description string
	}

	tests := []test{
		{
			input:       time.Minute * 5,
			want:        "  5",
			description: "5 minutes",
		},
		{
			input:       0,
			want:        "  0",
			description: "0 minutes",
		},
		{
			input:       time.Minute * 10,
			want:        " 10",
			description: "10 minutes",
		},
		{
			input:       time.Minute * 999,
			want:        "999",
			description: "999 minutes",
		},
		{
			input:       time.Minute * 9999,
			want:        "999",
			description: "9999 minutes",
		},
	}

	for _, tc := range tests {
		got := mins(tc.input)
		if got != tc.want {
			t.Fatal(
				"description:",
				tc.description,
				"got:",
				got,
				"want:",
				tc.want,
			)
		}
	}
}

func TestDrawSide(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		{
			input: "A.mkv",
			want:  "A   ",
		},
		{
			input: "AA.mkv",
			want:  "AA  ",
		},
		{
			input: "AAAA.mkv",
			want:  "AAAA",
		},
		{
			input: "AAAAA.mkv",
			want:  "ZZZZ",
		},
	}

	for _, tc := range tests {
		got := drawSide(tc.input)
		if got != tc.want {
			t.Fatal("got:", got, "want:", tc.want)
		}
	}
}
