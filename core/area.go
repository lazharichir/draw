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
	return Abs(area.Min.Y-area.Max.Y) + 1
}

func (area Area) Width() int64 {
	return Abs(area.Min.X-area.Max.X) + 1
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

func (area Area) Surface() int64 {
	return area.Height() * area.Width()
}

func (area Area) CountOverlappingPixels(other Area) int64 {
	is, does := area.Intersect(other)
	if !does {
		return 0
	}
	return is.Surface()

	// // Calculate the overlapping area
	// overlap := NewArea(
	// 	NewPoint(Max(area.Min.X, other.Min.X), Max(area.Min.Y, other.Min.Y)),
	// 	NewPoint(Min(area.Max.X, other.Max.X), Min(area.Max.Y, other.Max.Y)),
	// )

	// // If there is no overlap, return 0
	// if overlap.Min.X >= overlap.Max.X || overlap.Min.Y >= overlap.Max.Y {
	// 	return 0
	// }

	// // Return the overlapping area's surface size
	// return overlap.Surface()
}

func (area Area) String() string {
	return fmt.Sprintf("Area[Min: %s, Max: %s, Height: %d, Width: %d]", area.Min, area.Max, area.Height(), area.Width())
}
