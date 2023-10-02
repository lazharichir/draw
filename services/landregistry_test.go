package services_test

import (
	"context"
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
			TopLeft:     core.Point{X: 0, Y: 0},
			BottomRight: core.Point{X: 100, Y: 100},
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
