package core

type Lease struct {
	ID           string
	LeasholderID int
	TopLeft      Point
	BottomRight  Point
	Width        int
	Height       int
}
