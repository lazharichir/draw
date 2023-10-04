package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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
func (lr *LandRegistry) DeleteLease(ctx context.Context, id string) error {
	query := `DELETE FROM leases WHERE id = $1`
	_, err := lr.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed DeleteLease: %w", err)
	}
	return nil
}

func (lr *LandRegistry) SaveLease(ctx context.Context, lease core.Lease) error {
	query := `
		INSERT INTO "leases" ("id", "leaseholder_id", "canvas_id", "tl_x", "tl_y", "br_x", "br_y", "width", "height", "status", "start", "end", "price", "metadata", "updated_at", "updated_by", "created_at", "created_by")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		ON CONFLICT ("id") DO UPDATE SET
			"status" = excluded."status",
			"start" = excluded."start",
			"end" = excluded."end",
			"price" = excluded."price",
			"metadata" = excluded."metadata",
			"updated_at" = excluded."updated_at",
			"updated_by" = excluded."updated_by"
	`
	_, err := lr.db.ExecContext(
		ctx,
		query,
		lease.ID,
		lease.LeaseholderID,
		lease.CanvasID,
		lease.Area.Min.X,
		lease.Area.Min.Y,
		lease.Area.Max.X,
		lease.Area.Max.Y,
		lease.Area.Width(),
		lease.Area.Height(),
		lease.Status,
		lease.Start,
		lease.End,
		lease.Price,
		lease.Metadata,
		lease.UpdatedAt,
		lease.UpdatedBy,
		lease.CreatedAt,
		lease.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to save lease: %w", err)
	}
	return nil
}

func (lr *LandRegistry) GetLease(ctx context.Context, leaseID string) (*core.Lease, error) {
	leases, err := lr.GetLeasesByID(ctx, leaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lease: %w", err)
	}
	if len(leases) == 0 {
		return nil, nil
	}

	return &leases[0], nil
}

func (lr *LandRegistry) GetLeasesByID(ctx context.Context, ids ...string) ([]core.Lease, error) {
	leases := []core.Lease{}
	if len(ids) == 0 {
		return leases, nil
	}

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		SELECT 
			"id",
			"leaseholder_id",
			"canvas_id",
			"tl_x",
			"tl_y",
			"br_x",
			"br_y",
			"status",
			"start",
			"end",
			"price",
			"metadata",
			"updated_at",
			"updated_by",
			"created_at",
			"created_by"
		FROM "leases"
		WHERE "id" IN (%s)
	`, strings.Join(placeholders, ", "))

	rows, err := lr.db.QueryContext(ctx, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return leases, nil
		}
		return nil, fmt.Errorf("failed to get leases: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lease core.Lease
		err := rows.Scan(
			&lease.ID,
			&lease.LeaseholderID,
			&lease.CanvasID,
			&lease.Area.Min.X,
			&lease.Area.Min.Y,
			&lease.Area.Max.X,
			&lease.Area.Max.Y,
			&lease.Status,
			&lease.Start,
			&lease.End,
			&lease.Price,
			&lease.Metadata,
			&lease.UpdatedAt,
			&lease.UpdatedBy,
			&lease.CreatedAt,
			&lease.CreatedBy,
		)
		lease.Start = lease.Start.UTC()
		lease.End = lease.End.UTC()
		lease.UpdatedAt = lease.UpdatedAt.UTC()
		lease.CreatedAt = lease.CreatedAt.UTC()
		if err != nil {
			return nil, fmt.Errorf("failed to scan lease: %w", err)
		}
		leases = append(leases, lease)
	}

	return leases, nil
}

func (lr *LandRegistry) GetLeasesByPoint(ctx context.Context, canvasID int64, point core.Point) ([]core.Lease, error) {
	query := `
		SELECT id
		FROM leases
		WHERE
			canvas_id = $1
			AND (tl_x <= $2 AND br_x > $2) 
			AND (br_y >= $3 AND tl_y <= $3) 
	`
	args := []any{
		canvasID,
		point.X,
		point.Y,
	}

	rows, err := lr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed GetLeasesByPoint: %w", err)
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("failed GetLeasesByPoint scan: %w", err)
		}
		ids = append(ids, id)
	}

	return lr.GetLeasesByID(ctx, ids...)
}

func (lr *LandRegistry) GetLeasesByArea(ctx context.Context, canvasID int64, area core.Area) ([]core.Lease, error) {
	query := `
		SELECT id
		FROM leases
		WHERE 
			canvas_id = $1 AND
			(
				(tl_x <= $2 AND br_x > $2) OR
				(tl_x < $3 AND br_x >= $3)
			) AND
			(
				(br_y >= $4 AND tl_y <= $4) OR
				(br_y > $5 AND tl_y <= $5)
			)
	`
	args := []any{
		canvasID,
		area.Max.X,
		area.Min.X,
		area.Max.Y,
		area.Min.Y,
	}

	rows, err := lr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed GetLeasesByArea: %w", err)
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("failed GetLeasesByArea scan: %w", err)
		}
		ids = append(ids, id)
	}
	return lr.GetLeasesByID(ctx, ids...)
}

func (lr *LandRegistry) CanDrawPixel(ctx context.Context, canvasID int64, drawerID int64, pixel core.Pixel) (bool, error) {
	leases, err := lr.GetLeasesByPoint(ctx, canvasID, pixel.Point)
	if err != nil {
		return false, fmt.Errorf("CanDrawPixel: %w", err)
	}

	// if no leases, then the pixel is free to draw
	if len(leases) == 0 {
		return true, nil
	}

	// if there are active leases, then the pixel is free to draw if the drawer owns it
	for _, lease := range leases {
		pixelInLease := lease.Area.ContainsPoint(pixel.Point)
		if !pixelInLease {
			continue
		}

		if !lease.IsActiveAt(time.Now()) {
			continue
		}

		if lease.LeaseholderID == drawerID {
			return true, nil //fmt.Errorf("CanDrawPixel: %d cannot draw in %s", drawerID, pixel.Point.String())
		}
	}

	return false, nil
}

func (lr *LandRegistry) CanDrawInArea(ctx context.Context, canvasID int64, drawerID int64, area core.Area) (bool, error) {
	leases, err := lr.GetLeasesByArea(ctx, canvasID, area)
	if err != nil {
		return false, fmt.Errorf("CanDrawPixel: %w", err)
	}

	// if no leases, then the area is free to draw
	if len(leases) == 0 {
		return true, nil
	}

	// if there are active leases, then the area is free to draw if the drawer owns all of its pixels
	now := time.Now()
	for _, lease := range leases {
		// Ignore the lease if it is not active at the current time.
		if !lease.IsActiveAt(now) {
			continue
		}

		// Ignore the lease if it does not intersect the given area.
		if !lease.Area.IntersectsArea(area) {
			continue
		}

		// Not allowed if one of the relevant leases is not owned by the drawer.
		if lease.LeaseholderID == drawerID {
			return true, nil //fmt.Errorf("CanDrawInArea: %d cannot draw in %s", drawerID, area.String())
		}
	}

	return false, nil
}
