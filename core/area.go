package core

import "fmt"

func NewArea(topLeft, bottomRight Point) Area {
	a := Area{TopLeft: topLeft, BottomRight: bottomRight}
	a.MaybeSwapPoints()
	return a
}

func NewAreaWH(topLeft Point, width, height int64) Area {
	bottomRight := Point{X: topLeft.X + width, Y: topLeft.Y + height}
	return NewArea(topLeft, bottomRight)
}

type Area struct {
	TopLeft     Point
	BottomRight Point
}

func (area *Area) MaybeSwapPoints() {
	if area.TopLeft.X > area.BottomRight.X {
		area.TopLeft.X, area.BottomRight.X = area.BottomRight.X, area.TopLeft.X
	}
	if area.TopLeft.Y > area.BottomRight.Y {
		area.TopLeft.Y, area.BottomRight.Y = area.BottomRight.Y, area.TopLeft.Y
	}
}

func (area Area) Equal(other Area) bool {
	return area.TopLeft == other.TopLeft && area.BottomRight == other.BottomRight
}

func (area Area) Height() int64 {
	return Abs(area.TopLeft.Y-area.BottomRight.Y) + 1
}

func (area Area) Width() int64 {
	return Abs(area.TopLeft.X-area.BottomRight.X) + 1
}

func (area Area) IsLandscape() bool {
	return area.Width() > area.Height()
}

func (area Area) IsPortrait() bool {
	return area.Height() > area.Width()
}

func (area Area) ContainsPoint(p Point) bool {
	return p.X >= area.TopLeft.X && p.X <= area.BottomRight.X && p.Y >= area.TopLeft.Y && p.Y <= area.BottomRight.Y
}

func (area Area) ContainsArea(a Area) bool {
	return area.ContainsPoint(a.TopLeft) && area.ContainsPoint(a.BottomRight)
}

func (area Area) IntersectsArea(a Area) bool {
	return area.ContainsPoint(a.TopLeft) || area.ContainsPoint(a.BottomRight)
}

func (area Area) Surface() int64 {
	return area.Height() * area.Width()
}

func (area Area) CountOverlappingPixels(other Area) int64 {
	// Calculate the overlapping area
	overlap := Area{
		TopLeft: Point{
			X: Max(area.TopLeft.X, other.TopLeft.X),
			Y: Max(area.TopLeft.Y, other.TopLeft.Y),
		},
		BottomRight: Point{
			X: Min(area.BottomRight.X, other.BottomRight.X),
			Y: Min(area.BottomRight.Y, other.BottomRight.Y),
		},
	}

	// If there is no overlap, return 0
	if overlap.TopLeft.X >= overlap.BottomRight.X || overlap.TopLeft.Y >= overlap.BottomRight.Y {
		return 0
	}

	// Return the overlapping area's surface size
	return overlap.Surface()
}

func (area Area) String() string {
	return fmt.Sprintf("Area[TopLeft: %s, BottomRight: %s, Height: %d, Width: %d]", area.TopLeft, area.BottomRight, area.Height(), area.Width())
}
