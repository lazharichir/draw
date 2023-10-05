package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestNewAreaSquare(t *testing.T) {
	{
		// Test that a square area is created correctly.
		min := Pt(0, 0)
		side := int64(1024)
		expected := NewArea(Pt(0, 0), Pt(1024, 1024))
		actual := NewAreaSquare(min, side)
		assert.Equal(t, expected, actual)
	}
	{
		// Test that a square area is created correctly.
		min := Pt(-1024, 0)
		side := int64(1024)
		expected := NewArea(Pt(-1024, 0), Pt(0, 1024))
		actual := NewAreaSquare(min, side)
		assert.Equal(t, expected, actual)
	}
}

func TestNewArea(t *testing.T) {
	topLeft := Point{X: 0, Y: 0}
	bottomRight := Point{X: 10, Y: 10}
	expected := Area{Min: topLeft, Max: bottomRight}
	actual := NewArea(topLeft, bottomRight)
	assert.Equal(t, expected, actual)
}

func TestNewAreaWH(t *testing.T) {
	topLeft := Point{X: 0, Y: 0}
	width := int64(10)
	height := int64(10)
	expected := Area{Min: topLeft, Max: Point{X: 10, Y: 10}}
	actual := NewAreaWH(topLeft, width, height)
	assert.Equal(t, expected, actual)
}

func TestArea_Height(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 0, Y: 10},
	}
	expected := int64(10)
	if actual := area.Height(); actual != expected {
		t.Errorf("Expected height to be %d, but got %d", expected, actual)
	}
}

func TestArea_Width(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 0},
	}
	expected := int64(10)
	if actual := area.Width(); actual != expected {
		t.Errorf("Expected width to be %d, but got %d", expected, actual)
	}
}

func TestArea_IsLandscape(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 0},
	}
	if !area.IsLandscape() {
		t.Error("Expected area to be landscape, but it was not")
	}

	area = Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 0, Y: 10},
	}
	if area.IsLandscape() {
		t.Error("Expected area to not be landscape, but it was")
	}
}

func TestArea_IsPortrait(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 0, Y: 10},
	}
	if !area.IsPortrait() {
		t.Error("Expected area to be portrait, but it was not")
	}

	area = Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 0},
	}
	if area.IsPortrait() {
		t.Error("Expected area to not be portrait, but it was")
	}
}

func TestArea_ContainsPoint(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}
	p := Point{X: 5, Y: 5}
	if !area.ContainsPoint(p) {
		t.Errorf("Expected area to contain point %v, but it did not", p)
	}

	p = Point{X: 15, Y: 15}
	if area.ContainsPoint(p) {
		t.Errorf("Expected area to not contain point %v, but it did", p)
	}
}

func TestArea_ContainsArea(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}
	a := Area{
		Min: Point{X: 2, Y: 2},
		Max: Point{X: 8, Y: 8},
	}
	if !area.ContainsArea(a) {
		t.Errorf("Expected area to contain area %v, but it did not", a)
	}

	a = Area{
		Min: Point{X: 12, Y: 12},
		Max: Point{X: 18, Y: 18},
	}
	if area.ContainsArea(a) {
		t.Errorf("Expected area to not contain area %v, but it did", a)
	}

	a = Area{
		Min: Point{X: 8, Y: 8},
		Max: Point{X: 12, Y: 12},
	}
	if area.ContainsArea(a) {
		t.Errorf("Expected area to not contain area %v, but it did", a)
	}
}

func TestArea_IntersectsArea(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}
	a := Area{
		Min: Point{X: 2, Y: 2},
		Max: Point{X: 8, Y: 8},
	}
	if !area.IntersectsArea(a) {
		t.Errorf("Expected area to intersect area %v, but it did not", a)
	}

	a = Area{
		Min: Point{X: 12, Y: 12},
		Max: Point{X: 18, Y: 18},
	}
	if area.IntersectsArea(a) {
		t.Errorf("Expected area to not intersect area %v, but it did", a)
	}
}

func TestAreaPoints(t *testing.T) {
	{
		// Test that an empty area returns an empty slice of points.
		area := Area{}
		expected := []Point{
			Pt(0, 0),
		}
		actual := area.Points()
		assert.Equal(t, expected, actual)
	}
	{
		area := NewArea(Pt(0, 0), Pt(3, 3))
		expected := []Point{
			Pt(0, 0),
			Pt(0, 1),
			Pt(0, 2),
			Pt(0, 3),
			Pt(1, 0),
			Pt(1, 1),
			Pt(1, 2),
			Pt(1, 3),
			Pt(2, 0),
			Pt(2, 1),
			Pt(2, 2),
			Pt(2, 3),
			Pt(3, 0),
			Pt(3, 1),
			Pt(3, 2),
			Pt(3, 3),
		}
		actual := area.Points()
		assert.Equal(t, expected, actual)
	}
	{
		// one-column area
		area := NewArea(Pt(0, 0), Pt(0, 3))
		expected := []Point{
			Pt(0, 0),
			Pt(0, 1),
			Pt(0, 2),
			Pt(0, 3),
		}
		actual := area.Points()
		assert.Equal(t, expected, actual)
	}
	{
		// one-row area
		area := NewArea(Pt(0, 0), Pt(2, 0))
		expected := []Point{
			Pt(0, 0),
			Pt(1, 0),
			Pt(2, 0),
		}
		actual := area.Points()
		assert.Equal(t, expected, actual)
	}
}

func TestArea_Surface(t *testing.T) {
	area := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}
	expected := int64(121)
	assert.Equal(t, expected, area.Surface())
}

func TestArea_SurfaceOne(t *testing.T) {
	area := Area{
		Min: Point{X: 3, Y: 3},
		Max: Point{X: 3, Y: 3},
	}
	expected := int64(1)
	assert.Equal(t, expected, area.Surface())
}

func TestArea_SurfaceTwo(t *testing.T) {
	area := Area{
		Min: Point{X: 3, Y: 3},
		Max: Point{X: 3, Y: 4},
	}
	expected := int64(2)
	assert.Equal(t, expected, area.Surface())
}

func TestArea_String(t *testing.T) {
	{
		area := Area{
			Min: Point{X: 0, Y: 0},
			Max: Point{X: 10, Y: 10},
		}
		expected := "Area[Min: (0,0), Max: (10,10), Height: 10, Width: 10, Points: 121]"
		assert.Equal(t, expected, area.String())
	}
	{
		area := Area{
			Min: Point{X: 0, Y: 0},
			Max: Point{X: 0, Y: 0},
		}
		expected := "Area[Min: (0,0), Max: (0,0), Height: 0, Width: 0, Points: 1]"
		assert.Equal(t, expected, area.String())
	}
}

func TestArea_Equal(t *testing.T) {
	area1 := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}
	area2 := Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}
	if !area1.Equal(area2) {
		t.Error("Expected areas to be equal, but they were not")
	}

	area1 = Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}
	area2 = Area{
		Min: Point{X: 5, Y: 5},
		Max: Point{X: 15, Y: 15},
	}
	if area1.Equal(area2) {
		t.Error("Expected areas to not be equal, but they were")
	}
}

func TestArea_Intersect(t *testing.T) {
	area1 := NewArea(NewPoint(0, 0), NewPoint(10, 10))
	area2 := NewArea(NewPoint(15, 15), NewPoint(20, 20))
	_, does := area1.Intersect(area2)
	assert.False(t, does)
}

func TestArea_CountOverlappingPixels(t *testing.T) {
	{
		area1 := NewArea(NewPoint(0, 0), NewPoint(10, 10))
		area2 := NewArea(NewPoint(5, 5), NewPoint(15, 15))
		assert.Equal(t, int64(36), area1.CountOverlappingPixels(area2))
	}
	{
		area1 := NewArea(NewPoint(0, 0), NewPoint(10, 10))
		area2 := NewArea(NewPoint(15, 15), NewPoint(20, 20))
		assert.Equal(t, int64(0), area1.CountOverlappingPixels(area2))
	}
}

func TestArea_MaybeSwapPoints(t *testing.T) {
	// Create a new area with the top-left point greater than the bottom-right point.
	area := Area{
		Min: Point{X: 10, Y: 10},
		Max: Point{X: 0, Y: 0},
	}

	// Call the MaybeSwapPoints method.
	area.MaybeSwapPoints()

	// Test that the points have been swapped.
	assert.Equal(t, Point{X: 0, Y: 0}, area.Min)
	assert.Equal(t, Point{X: 10, Y: 10}, area.Max)

	// Create a new area with the top-left point less than the bottom-right point.
	area = Area{
		Min: Point{X: 0, Y: 0},
		Max: Point{X: 10, Y: 10},
	}

	// Call the MaybeSwapPoints method.
	area.MaybeSwapPoints()

	// Test that the points have not been swapped.
	assert.Equal(t, Point{X: 0, Y: 0}, area.Min)
	assert.Equal(t, Point{X: 10, Y: 10}, area.Max)
}

func TestSortAreasFn(t *testing.T) {
	{
		// Test that an empty slice is sorted correctly.
		areas := []Area{}
		expected := []Area{}
		slices.SortFunc(areas, SortAreasFn)
		assert.Equal(t, expected, areas)
	}

	{
		// Test that two simple slices are the same once both are sorted.
		left := []Area{NewArea(Pt(0, 0), Pt(9999, 2))}
		right := []Area{NewArea(Pt(9999, 2), Pt(0, 0))}
		slices.SortFunc(left, SortAreasFn)
		slices.SortFunc(right, SortAreasFn)
		assert.Equal(t, right, left)
	}

	{
		// Test that a slice with multiple elements is sorted correctly.
		areas3 := []Area{
			NewArea(Pt(1024, 1024), Pt(2047, 2047)),
			NewArea(Pt(1024, 0), Pt(2047, 1023)),
			NewArea(Pt(0, 1024), Pt(1023, 2047)),
			NewArea(Pt(0, 0), Pt(1023, 1023)),
		}
		expected3 := []Area{
			NewArea(Pt(0, 0), Pt(1023, 1023)),
			NewArea(Pt(0, 1024), Pt(1023, 2047)),
			NewArea(Pt(1024, 0), Pt(2047, 1023)),
			NewArea(Pt(1024, 1024), Pt(2047, 2047)),
		}
		slices.SortFunc(areas3, SortAreasFn)
		assert.Equal(t, expected3, areas3)
	}
}
