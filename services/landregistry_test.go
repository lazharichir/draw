package services_test

import (
	"context"
	"image/color"
	"testing"
	"time"

	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/services"
	"github.com/lazharichir/draw/storage"
	"github.com/lazharichir/draw/utils"
	"github.com/stretchr/testify/assert"
)

var lr *services.LandRegistry

func init() {
	// Create a new LandRegistry instance.
	db := storage.NewPG()
	lr = services.NewLandRegistry(db)
}

func TestLandRegistry_SaveLease(t *testing.T) {
	// Create a test lease.
	now := time.Now().UTC()
	lease := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 123,
		CanvasID:      456,
		Area: core.Area{
			Min: core.Point{X: 0, Y: 0},
			Max: core.Point{X: 100, Y: 100},
		},
		Status:    "active",
		Start:     now,
		End:       now.Add(time.Hour),
		Price:     1000,
		Metadata:  core.Metadata{"foo": "bar"},
		UpdatedAt: now,
		UpdatedBy: 789,
		CreatedAt: now,
		CreatedBy: 456,
	}

	// Save the lease.
	err := lr.SaveLease(context.Background(), lease)
	assert.NoError(t, err)

	// Retrieve the lease from the database.
	retrievedLease, err := lr.GetLease(context.Background(), lease.ID)
	assert.NoError(t, err)
	assert.Equal(t, lease, *retrievedLease)

	// Update the lease.
	lease.Status = "inactive"
	lease.UpdatedAt = time.Now().UTC()
	err = lr.SaveLease(context.Background(), lease)
	assert.NoError(t, err)

	// Retrieve the updated lease from the database.
	retrievedUpdatedLease, err := lr.GetLease(context.Background(), lease.ID)
	assert.NoError(t, err)
	assert.Equal(t, lease, *retrievedUpdatedLease)

	// Delete the lease.
	err = lr.DeleteLease(context.Background(), lease.ID)
	assert.NoError(t, err)

	// Retrieve the deleted lease from the database.
	retrievedDeletedLease, err := lr.GetLease(context.Background(), lease.ID)
	assert.NoError(t, err)
	assert.Nil(t, retrievedDeletedLease)
}

func TestLandRegistry_GetLeasesByPoint(t *testing.T) {
	// Create a new LandRegistry instance.
	db := storage.NewPG()
	lr := services.NewLandRegistry(db)

	// Create test leases.
	now := time.Now().UTC()
	lease1 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 123,
		CanvasID:      0,
		Area:          core.NewArea(core.NewPoint(0, 0), core.NewPoint(100, 100)),
		Status:        "active",
		Start:         now,
		End:           now.Add(time.Hour),
		Price:         1000,
		Metadata:      core.Metadata{"foo": "bar"},
		UpdatedAt:     now,
		UpdatedBy:     789,
		CreatedAt:     now,
		CreatedBy:     456,
	}
	lease2 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 456,
		CanvasID:      0,
		Area:          core.NewArea(core.NewPoint(50, 50), core.NewPoint(150, 150)),
		Status:        "active",
		Start:         now,
		End:           now.Add(time.Hour),
		Price:         2000,
		Metadata:      core.Metadata{"baz": "qux"},
		UpdatedAt:     now,
		UpdatedBy:     123,
		CreatedAt:     now,
		CreatedBy:     789,
	}

	// Save the test leases.
	err := lr.SaveLease(context.Background(), lease1)
	assert.NoError(t, err)
	err = lr.SaveLease(context.Background(), lease2)
	assert.NoError(t, err)

	// Retrieve the leases by point.
	leases, err := lr.GetLeasesByPoint(context.Background(), 0, core.Point{X: 75, Y: 75})
	assert.NoError(t, err)
	assert.Len(t, leases, 2)

	leases, err = lr.GetLeasesByPoint(context.Background(), 0, core.Point{X: 25, Y: 25})
	assert.NoError(t, err)
	assert.Len(t, leases, 1)
	assert.Equal(t, lease1.ID, leases[0].ID)

	leases, err = lr.GetLeasesByPoint(context.Background(), 0, core.Point{X: 175, Y: 175})
	assert.NoError(t, err)
	assert.Len(t, leases, 0)

	// Delete the test leases.
	err = lr.DeleteLease(context.Background(), lease1.ID)
	assert.NoError(t, err)
	err = lr.DeleteLease(context.Background(), lease2.ID)
	assert.NoError(t, err)
}

func TestLandRegistry_GetLeasesByArea(t *testing.T) {
	// Create a new LandRegistry instance.
	db := storage.NewPG()
	lr := services.NewLandRegistry(db)

	// Create test leases.
	now := time.Now().UTC()
	lease1 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 123,
		CanvasID:      0,
		Area:          core.NewArea(core.NewPoint(0, 0), core.NewPoint(100, 100)),
		Status:        "active",
		Start:         now,
		End:           now.Add(time.Hour),
		Price:         1000,
		Metadata:      core.Metadata{"foo": "bar"},
		UpdatedAt:     now,
		UpdatedBy:     789,
		CreatedAt:     now,
		CreatedBy:     456,
	}
	lease2 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 456,
		CanvasID:      0,
		Area:          core.NewArea(core.NewPoint(50, 50), core.NewPoint(150, 150)),
		Status:        "active",
		Start:         now,
		End:           now.Add(time.Hour),
		Price:         2000,
		Metadata:      core.Metadata{"baz": "qux"},
		UpdatedAt:     now,
		UpdatedBy:     123,
		CreatedAt:     now,
		CreatedBy:     789,
	}

	// Save the test leases.
	err := lr.SaveLease(context.Background(), lease1)
	assert.NoError(t, err)
	err = lr.SaveLease(context.Background(), lease2)
	assert.NoError(t, err)

	// Retrieve the leases by area.
	leases, err := lr.GetLeasesByArea(context.Background(), 0, core.NewArea(core.NewPoint(25, 25), core.NewPoint(75, 75)))
	assert.NoError(t, err)
	assert.Len(t, leases, 2)

	leases, err = lr.GetLeasesByArea(context.Background(), 0, core.NewArea(core.NewPoint(125, 125), core.NewPoint(175, 175)))
	assert.NoError(t, err)
	assert.Len(t, leases, 1)

	// Delete the test leases.
	err = lr.DeleteLease(context.Background(), lease1.ID)
	assert.NoError(t, err)
	err = lr.DeleteLease(context.Background(), lease2.ID)
	assert.NoError(t, err)
}

func TestLandRegistry_CanDrawPixel(t *testing.T) {
	// Create a new LandRegistry instance.
	db := storage.NewPG()
	lr := services.NewLandRegistry(db)

	// Create test leases.
	now := time.Now().UTC()
	lease1 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 123,
		CanvasID:      0,
		Area: core.Area{
			Min: core.Point{X: 0, Y: 0},
			Max: core.Point{X: 100, Y: 100},
		},
		Status:    "active",
		Start:     now,
		End:       now.Add(time.Hour),
		Price:     1000,
		Metadata:  core.Metadata{"foo": "bar"},
		UpdatedAt: now,
		UpdatedBy: 789,
		CreatedAt: now,
		CreatedBy: 456,
	}
	lease2 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 456,
		CanvasID:      0,
		Area: core.Area{
			Min: core.Point{X: 101, Y: 101},
			Max: core.Point{X: 300, Y: 300},
		},
		Status:    "active",
		Start:     now,
		End:       now.Add(time.Hour),
		Price:     2000,
		Metadata:  core.Metadata{"baz": "qux"},
		UpdatedAt: now,
		UpdatedBy: 123,
		CreatedAt: now,
		CreatedBy: 789,
	}

	// Save the test leases.
	err := lr.SaveLease(context.Background(), lease1)
	assert.NoError(t, err)
	err = lr.SaveLease(context.Background(), lease2)
	assert.NoError(t, err)

	// Test that a drawer can draw a pixel in a lease they own.
	canDraw, err := lr.CanDrawPixel(context.Background(), lease1.CanvasID, lease1.LeaseholderID, core.Pixel{
		Point: core.Point{X: 25, Y: 25},
		RGBA:  color.RGBA{0, 0, 0, 0},
	})
	assert.NoError(t, err)
	assert.True(t, canDraw)

	// Test that a drawer cannot draw a pixel in a lease they do not own.
	canDraw, err = lr.CanDrawPixel(context.Background(), lease1.CanvasID, lease2.LeaseholderID, core.Pixel{
		Point: core.Point{X: 75, Y: 75},
		RGBA:  color.RGBA{0, 0, 0, 0},
	})
	assert.NoError(t, err)
	assert.False(t, canDraw)

	// Test that a drawer cannot draw a pixel outside of any leases.
	canDraw, err = lr.CanDrawPixel(context.Background(), lease1.CanvasID, lease1.LeaseholderID, core.Pixel{
		Point: core.Point{X: 200, Y: 200},
		RGBA:  color.RGBA{0, 0, 0, 0},
	})
	assert.NoError(t, err)
	assert.False(t, canDraw)

	// Delete the test leases.
	err = lr.DeleteLease(context.Background(), lease1.ID)
	assert.NoError(t, err)
	err = lr.DeleteLease(context.Background(), lease2.ID)
	assert.NoError(t, err)
}

func TestLandRegistry_CanDrawInArea(t *testing.T) {
	// Create a new LandRegistry instance.
	db := storage.NewPG()
	lr := services.NewLandRegistry(db)

	// Create test leases.
	now := time.Now().UTC()
	lease1 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 123,
		CanvasID:      0,
		Area:          core.NewArea(core.NewPoint(0, 0), core.NewPoint(100, 100)),
		Status:        "active",
		Start:         now,
		End:           now.Add(time.Hour),
		Price:         1000,
		Metadata:      core.Metadata{"foo": "bar"},
		UpdatedAt:     now,
		UpdatedBy:     789,
		CreatedAt:     now,
		CreatedBy:     456,
	}
	lease2 := core.Lease{
		ID:            utils.NewLeaseID(),
		LeaseholderID: 456,
		CanvasID:      0,
		Area:          core.NewArea(core.NewPoint(200, 200), core.NewPoint(300, 300)),
		Status:        "active",
		Start:         now,
		End:           now.Add(time.Hour),
		Price:         2000,
		Metadata:      core.Metadata{"baz": "qux"},
		UpdatedAt:     now,
		UpdatedBy:     123,
		CreatedAt:     now,
		CreatedBy:     789,
	}

	// Save the test leases.
	err := lr.SaveLease(context.Background(), lease1)
	assert.NoError(t, err)
	err = lr.SaveLease(context.Background(), lease2)
	assert.NoError(t, err)

	// Test that a drawer can draw in an area they own.
	canDraw, err := lr.CanDrawInArea(context.Background(), 0, lease1.LeaseholderID, core.NewArea(core.NewPoint(25, 25), core.NewPoint(75, 75)))
	assert.NoError(t, err)
	assert.True(t, canDraw)

	// Test that a drawer cannot draw in an area they do not own.
	canDraw, err = lr.CanDrawInArea(context.Background(), 0, lease2.LeaseholderID, core.NewArea(core.NewPoint(25, 25), core.NewPoint(75, 75)))
	assert.NoError(t, err)
	assert.False(t, canDraw)

	// Test that a drawer cannot draw in an area outside of any leases.
	canDraw, err = lr.CanDrawInArea(context.Background(), 0, lease2.LeaseholderID, core.NewArea(core.NewPoint(200, 200), core.NewPoint(250, 250)))
	assert.NoError(t, err)
	assert.True(t, canDraw)

	// Delete the test leases.
	err = lr.DeleteLease(context.Background(), lease1.ID)
	assert.NoError(t, err)
	err = lr.DeleteLease(context.Background(), lease2.ID)
	assert.NoError(t, err)
}
