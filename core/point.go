package core

import (
	"fmt"
	"math"
)

func Pt(x, y int64) Point {
	return NewPoint(x, y)
}

func NewPoint(x, y int64) Point {
	return Point{X: x, Y: y}
}

type Point struct {
	X int64
	Y int64
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func (p Point) DistanceTo(other Point) float64 {
	dx := float64(p.X - other.X)
	dy := float64(p.Y - other.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func (p Point) Translate(dx, dy int64) Point {
	return Point{X: p.X + dx, Y: p.Y + dy}
}

func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p Point) Subtract(other Point) Point {
	return Point{X: p.X - other.X, Y: p.Y - other.Y}
}

func (p Point) Multiply(scalar int64) Point {
	return Point{X: p.X * scalar, Y: p.Y * scalar}
}

func (p Point) Divide(scalar int64) Point {
	return Point{X: p.X / scalar, Y: p.Y / scalar}
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) IsOrigin() bool {
	return p.X == 0 && p.Y == 0
}
