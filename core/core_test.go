package core

import (
	"math/rand"
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
		{"zero center point", Pt(0, 0), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(1023, 1023), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(1023, 0), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(0, 1023), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(512, 512), 1024, NewArea(Pt(0, 0), Pt(1024, 1024))},
		{"", Pt(1600, 1700), 1024, NewArea(Pt(1024, 1024), Pt(2048, 2048))},
		{"big point", Pt(122221, 2047), 1024, NewArea(Pt(121856, 1024), Pt(122880, 2048))},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			actual := GetTileAreaFromPoint(tc.input, tc.side)
			assert.Equal(t, tc.expected, actual, "GetTileAreaFromPoint(%v, %v) should return %v [%s]", tc.input, tc.side, tc.expected, tc.label)
		})
	}
}

func TestGetTileAreasFromPoints(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label    string
		side     int64
		input    []Point
		expected []Area
	}{
		{
			"empty slice",
			1024,
			[]Point{},
			[]Area{},
		},
		{
			"single point",
			1024,
			[]Point{Pt(50, 50)},
			[]Area{NewArea(Pt(0, 0), Pt(1024, 1024))},
		},
		{
			"single negative point",
			1024,
			[]Point{Pt(-50, -50)},
			[]Area{NewArea(Pt(-1024, -1024), Pt(0, 0))},
		},
		{
			"partially negative point",
			1024,
			[]Point{Pt(-50, 50)},
			[]Area{NewArea(Pt(-1024, 0), Pt(0, 1024))},
		},
		{
			"multiple points in the same tile",
			1024,
			[]Point{Pt(50, 50), Pt(100, 100)},
			[]Area{NewArea(Pt(0, 0), Pt(1024, 1024))},
		},
		{
			"multiple points in different tiles",
			1024,
			[]Point{Pt(1600, 1700), Pt(50, 50), Pt(2000, 2000), Pt(1025, 1025)},
			[]Area{
				NewArea(Pt(0, 0), Pt(1024, 1024)),
				NewArea(Pt(1024, 1024), Pt(2048, 2048)),
			},
		},
		{
			"multiple points in different tiles, including negative coordinates",
			1024,
			[]Point{Pt(1600, 1700), Pt(50, 50), Pt(2000, 2000), Pt(1025, 1025), Pt(-50, -50)},
			[]Area{
				NewArea(Pt(-1024, -1024), Pt(0, 0)),
				NewArea(Pt(0, 0), Pt(1024, 1024)),
				NewArea(Pt(1024, 1024), Pt(2048, 2048)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			actual := GetTileAreasFromPoints(tc.side, tc.input...)
			assert.Equal(t, tc.expected, actual, "GetTileAreasFromPoints(%v, %v) should return %v [%s]", tc.side, tc.input, tc.expected, tc.label)
		})
	}
}

func BenchmarkGetTileAreasFromPointsSync(b *testing.B) {
	points := []Point{}
	for i := 0; i < 10000; i++ {
		points = append(points, Pt(rand.Int63n(10000), rand.Int63n(10000)))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GetTileAreasFromPoints(1024, points...)
	}
}

func BenchmarkGetTileAreasFromPointsConcurrent(b *testing.B) {
	points := []Point{}
	for i := 0; i < 10000; i++ {
		points = append(points, Pt(rand.Int63n(10000), rand.Int63n(10000)))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getTileAreasFromPointsConcurrently(1024, points...)
	}
}
