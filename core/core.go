package core

import (
	"sync"

	"golang.org/x/exp/slices"
)

// Create a function that takes a point and returns the square tile it belongs to.
// e.g., GetTileAreaFromPoint(Point{X: -50, Y: -50}, 1024) returns tile area (-1024, -1024) -> (0, 0)
// e.g., GetTileAreaFromPoint(Point{X: 50, Y: 50}, 1024) returns tile area (0, 0) -> (1024, 1024)
// e.g., GetTileAreaFromPoint(Point{X: 1600, Y: 1700}, 1024) returns tile area (1024, 1024) -> (2048, 2048)
func GetTileAreaFromPoint(pt Point, side int64) Area {
	// Calculate the tile coordinates.
	tileX := pt.X / side
	tileY := pt.Y / side

	// Calculate the minimum and maximum coordinates of the tile.
	minX := tileX * side
	minY := tileY * side
	maxX := minX + side
	maxY := minY + side

	// Adjust the minimum and maximum coordinates for negative coordinates.
	if pt.X < 0 {
		minX -= side
		maxX -= side
	}
	if pt.Y < 0 {
		minY -= side
		maxY -= side
	}

	return Area{Min: Pt(minX, minY), Max: Pt(maxX, maxY)}
}

func GetTileAreasFromPoints(side int64, points ...Point) []Area {
	areas := []Area{}
	if len(points) == 0 {
		return areas
	}

	// Create a map of tile areas to avoid duplicates.
	areasMap := map[Area]bool{}

	for _, pt := range points {
		area := GetTileAreaFromPoint(pt, side)
		areasMap[area] = true
	}

	// Convert the map to a slice.
	for area := range areasMap {
		areas = append(areas, area)
	}

	// Sort the slice so it is deterministic.
	// Sort by Min X, then Min Y, then Max X, then Max Y.
	slices.SortFunc(areas, SortAreasFn)

	return areas
}

func getTileAreasFromPointsConcurrently(side int64, points ...Point) []Area {
	areas := []Area{}
	if len(points) == 0 {
		return areas
	}

	// Create a map of tile areas to avoid duplicates.
	areasMap := make(map[Area]bool)

	// Create a mutex to synchronize access to the areasMap.
	var mutex sync.Mutex

	// Create a wait group to wait for all goroutines to finish.
	var wg sync.WaitGroup

	// Spawn a goroutine for each point to calculate the tile area.
	for _, pt := range points {
		wg.Add(1)
		go func(pt Point) {
			defer wg.Done()
			area := GetTileAreaFromPoint(pt, side)

			// Lock the mutex before accessing the areasMap.
			mutex.Lock()
			areasMap[area] = true
			mutex.Unlock()
		}(pt)
	}

	// Wait for all goroutines to finish.
	wg.Wait()

	// Convert the map to a slice.
	for area := range areasMap {
		areas = append(areas, area)
	}

	// Sort the slice so it is deterministic.
	// Sort by Min X, then Min Y, then Max X, then Max Y.
	slices.SortFunc(areas, SortAreasFn)

	return areas
}
