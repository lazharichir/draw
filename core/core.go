package core

// Create a function that takes a point and returns the square tile it belongs to.
// e.g., GetTileAreaFromPoint(Point{X: 50, Y: 50}, 1024) returns tile area (0, 0) -> (1024, 1024)
// e.g., GetTileAreaFromPoint(Point{X: 1600, Y: 1700}, 1024) returns tile area (1024, 1024) -> (2048, 2048)
func GetTileAreaFromPoint(pt Point, side int64) Area {
	// Calculate the top-left corner of the tile that contains the point.
	tileX := (pt.X / side) * side
	tileY := (pt.Y / side) * side

	// Calculate the bottom-right corner of the tile.
	tileBottomRightX := tileX + side
	tileBottomRightY := tileY + side

	// Create and return the tile area.
	min := Pt(tileX, tileY)
	max := Pt(tileBottomRightX, tileBottomRightY)
	area := NewArea(min, max)
	return area
}
