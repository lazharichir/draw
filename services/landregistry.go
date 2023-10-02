package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lazharichir/draw/core"
)

var ErrCannotDrawInArea = func(drawerID int, topLeft, bottomRight core.Point) error {
	return fmt.Errorf("drawer %d cannot draw in area tl%v br%v", drawerID, topLeft, bottomRight)
}

type LandRegistry struct {
	db *sql.DB
}

func NewLandRegistry(db *sql.DB) *LandRegistry {
	return &LandRegistry{db}
}

func (lr *LandRegistry) SaveLease(ctx context.Context, lease core.Lease) error {
	return nil
}

func (lr *LandRegistry) GetLease(ctx context.Context, leaseID string) (*core.Lease, error) {
	return nil, nil
}

func (lr *LandRegistry) GetLeasesByArea(ctx context.Context, leaseID string) ([]core.Lease, error) {
	return nil, nil
}

func (lr *LandRegistry) LockArea(ctx context.Context, ownerID int, topLeft, bottomRight core.Point) error {
	return nil
}

func (lr *LandRegistry) CanDrawPixel(ctx context.Context, drawerID int, pixel core.Pixel) (bool, error) {
	return true, nil
}

func (lr *LandRegistry) CanDrawInArea(ctx context.Context, drawerID int, topLeft, bottomRight core.Point) (bool, error) {
	return true, nil
}
