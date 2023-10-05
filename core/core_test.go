package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTileAreaFromPoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label    string
		input    Point
		side     int64
		expected Area
	}{
		{"", Pt(0, 0), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(1023, 1023), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(1023, 0), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(0, 1023), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(512, 512), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(1600, 1700), 1024, NewArea(Pt(1024, 1024), Pt(2048, 2048))},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			actual := GetTileAreaFromPoint(tc.input, tc.side)
			assert.Equal(t, tc.expected, actual, "GetTileAreaFromPoint(%v, %v) should return %v [%s]", tc.input, tc.side, tc.expected, tc.label)
		})
	}
}
