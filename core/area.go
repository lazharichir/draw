package core

import "fmt"

func NewArea(min, max Point) Area {
	a := Area{Min: min, Max: max}
	return a.Canon()
}

func NewAreaWH(min Point, width, height int64) Area {
	max := Point{X: min.X + width, Y: min.Y + height}
	return NewArea(min, max)
}

func NewAreaSquare(min Point, side int64) Area {
	max := Pt(min.X+side, min.Y+side)
	return NewArea(min, max)
}

type Area struct {
	Min Point
	Max Point
}

func (area Area) Canon() Area {
	area.MaybeSwapPoints()
	return area
}

func (area *Area) MaybeSwapPoints() {
	if area.Max.X < area.Min.X {
		area.Min.X, area.Max.X = area.Max.X, area.Min.X
	}
	if area.Max.Y < area.Min.Y {
		area.Min.Y, area.Max.Y = area.Max.Y, area.Min.Y
	}
}

func (area Area) Equal(other Area) bool {
	return area.Min == other.Min && area.Max == other.Max
}

func (area Area) Height() int64 {
	return Abs(area.Min.Y - area.Max.Y)
}

func (area Area) Width() int64 {
	return Abs(area.Min.X - area.Max.X)
}

func (area Area) IsLandscape() bool {
	return area.Width() > area.Height()
}

func (area Area) IsPortrait() bool {
	return area.Height() > area.Width()
}

func (area Area) ContainsPoint(p Point) bool {
	return p.X >= area.Min.X && p.X <= area.Max.X && p.Y >= area.Min.Y && p.Y <= area.Max.Y
}

func (area Area) ContainsArea(a Area) bool {
	return area.ContainsPoint(a.Min) && area.ContainsPoint(a.Max)
}

func (area Area) IntersectsArea(a Area) bool {
	return area.ContainsPoint(a.Min) || area.ContainsPoint(a.Max)
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Area) Intersect(s Area) (Area, bool) {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r, r.Min.X <= r.Max.X && r.Min.Y <= r.Max.Y
}

func (area Area) Points() []Point {
	points := []Point{}
	for x := area.Min.X; x <= area.Max.X; x++ {
		for y := area.Min.Y; y <= area.Max.Y; y++ {
			points = append(points, Pt(x, y))
		}
	}
	return points
}

func (area Area) Surface() int64 {
	return int64(len(area.Points()))
}

func (area Area) CountOverlappingPixels(other Area) int64 {
	is, does := area.Intersect(other)
	if !does {
		return 0
	}
	return is.Surface()
}

func (area Area) String() string {
	return fmt.Sprintf("Area[Min: %s, Max: %s, Height: %d, Width: %d, Points: %d]", area.Min, area.Max, area.Height(), area.Width(), area.Surface())
}

// SortAreasFn sorts the slice so it is deterministic.
// Sort by Min X, then Min Y, then Max X, then Max Y.
func SortAreasFn(a, b Area) bool {
	if a.Min.X < b.Min.X {
		return true
	}
	if a.Min.Y < b.Min.Y {
		return true
	}
	if a.Max.X < b.Max.X {
		return true
	}
	if a.Max.Y < b.Max.Y {
		return true
	}
	return false
}
