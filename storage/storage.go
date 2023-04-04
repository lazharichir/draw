package storage

import (
	"database/sql"
	"image/color"

	"github.com/huandu/go-sqlbuilder"
	"github.com/lazharichir/draw/core"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

type PixelStore interface {
	GetPixelsFromTopLeft(canvasID, x, y, z int64) ([]core.Pixel, error)
	DrawPixelRGBA(canvasID, x, y int64, color color.RGBA) error
	ErasePixel(canvasID, x, y int64) error
}

type pgPixelStore struct {
	db  *sql.DB
	log *slog.Logger
}

func NewPGPixelStore(db *sql.DB, log *slog.Logger) PixelStore {
	return &pgPixelStore{db, log}
}

// ErasePixel implements PixelStore
// It deletes a pixel from the database
func (store *pgPixelStore) ErasePixel(canvasID int64, x int64, y int64) error {

	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom("pixels")
	db.Where(
		db.And(
			db.Equal("canvas_id", canvasID),
			db.Equal("x", x),
			db.Equal("y", y),
		),
	)

	query, args := db.Build()

	_, err := store.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// DrawPixelRGBA implements PixelStore
// It upserts a pixel in the database
func (store *pgPixelStore) DrawPixelRGBA(canvasID int64, x int64, y int64, color color.RGBA) error {

	sb := sqlbuilder.PostgreSQL.NewInsertBuilder()
	sb.InsertInto("pixels")
	sb.Cols("canvas_id", "x", "y", "r", "g", "b", "a", "drawn_at", "drawn_by")
	sb.Values(canvasID, x, y, color.R, color.G, color.B, color.A, "NOW()", 0)
	sb.SQL("ON CONFLICT (canvas_id, x, y) DO UPDATE SET r = EXCLUDED.r, g = EXCLUDED.g, b = EXCLUDED.b, a = EXCLUDED.a, drawn_at = EXCLUDED.drawn_at, drawn_by = EXCLUDED.drawn_by")

	query, args := sb.Build()

	_, err := store.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetPixels implements PixelStore
func (store *pgPixelStore) GetPixelsFromTopLeft(canvasID int64, tlX int64, tlY int64, width int64) ([]core.Pixel, error) {

	xFrom := tlX
	xTo := tlX + width
	yFrom := tlY
	yTo := tlY + width

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("x", "y", "r", "g", "b", "a")
	sb.From("pixels")
	sb.Where(
		sb.And(
			sb.Equal("canvas_id", canvasID),
			sb.Between("x", xFrom, xTo),
			sb.Between("y", yFrom, yTo),
		),
	)

	query, args := sb.Build()

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var pixels []core.Pixel

	for rows.Next() {
		var pixel core.Pixel
		err := rows.Scan(&pixel.X, &pixel.Y, &pixel.RGBA.R, &pixel.RGBA.G, &pixel.RGBA.B, &pixel.RGBA.A)
		if err != nil {
			return nil, err
		}
		pixels = append(pixels, pixel)
	}

	return pixels, nil
}

func NewPG() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres dbname=draw sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}
