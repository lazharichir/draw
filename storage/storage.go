package storage

import (
	"context"
	"database/sql"
	"fmt"
	"image/color"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/lazharichir/draw/core"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

type PixelStore interface {
	GetLatestPixelsForArea(canvasID int64, topLeft core.Point, bottomRight core.Point, after time.Time) ([]core.Pixel, error)
	GetPixelsFromTopLeft(canvasID, x, y, z int64) ([]core.Pixel, error)
	DrawPixelRGBA(canvasID, x, y int64, color color.RGBA) error
	DrawPixels(canvasID int64, pixels []core.Pixel) error
	ErasePixel(canvasID, x, y int64) error

	//
	SetLastChangedForAreas(ctx context.Context, canvasID int64, side int64, areas ...core.Area) error
	SetLastChangedForPoints(ctx context.Context, canvasID int64, side int64, points ...core.Point) error
	DeleteLastChangedForAreas(ctx context.Context, canvasID int64, side int64, areas ...core.Area) error
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
	return store.DrawPixels(canvasID, []core.Pixel{
		core.NewPixel(x, y, color),
	})
}

// DrawPixelRGBA implements PixelStore
// It upserts a pixel in the database
func (store *pgPixelStore) DrawPixels(canvasID int64, pixels []core.Pixel) error {
	chunks := chunkSlice(pixels, 1000)
	for _, chunk := range chunks {
		if err := store.drawPixelChunk(canvasID, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (store *pgPixelStore) drawPixelChunk(canvasID int64, pixels []core.Pixel) error {
	sb := sqlbuilder.PostgreSQL.NewInsertBuilder()
	sb.InsertInto("pixels")
	sb.Cols("canvas_id", "x", "y", "r", "g", "b", "a", "drawn_at", "drawn_by")

	for _, pixel := range pixels {
		sb.Values(canvasID, pixel.X, pixel.Y, pixel.RGBA.R, pixel.RGBA.G, pixel.RGBA.B, pixel.RGBA.A, "NOW()", 0)
	}

	sb.SQL("ON CONFLICT (canvas_id, x, y) DO UPDATE SET r = EXCLUDED.r, g = EXCLUDED.g, b = EXCLUDED.b, a = EXCLUDED.a, drawn_at = EXCLUDED.drawn_at, drawn_by = EXCLUDED.drawn_by")

	query, args := sb.Build()

	_, err := store.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetPixels implements PixelStore
func (store *pgPixelStore) GetLatestPixelsForArea(canvasID int64, topLeft core.Point, bottomRight core.Point, after time.Time) ([]core.Pixel, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("x", "y", "r", "g", "b", "a")
	sb.From("pixels")
	sb.Where(
		sb.And(
			sb.Equal("canvas_id", canvasID),
			sb.Between("x", topLeft.X, bottomRight.X),
			sb.Between("y", topLeft.Y, bottomRight.Y),
			sb.GreaterThan("drawn_at", after),
		),
	)

	query, args := sb.Build()
	// fmt.Println(`query`, query)
	// fmt.Printf("args %+#v \n", args)

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

// GetPixels implements PixelStore
type tilechange struct {
	CanvasID    int64
	X           int64
	Y           int64
	Side        int64
	LastChanged time.Time
}

func (t tilechange) Area() core.Area {
	min := core.Pt(t.X, t.Y)
	max := core.Pt(t.X+t.Side, t.Y+t.Side)
	return core.NewArea(min, max)
}

func mapTilechangesToAreas(items []tilechange) []core.Area {
	var areas []core.Area
	for _, item := range items {
		areas = append(areas, item.Area())
	}
	return areas
}

func (store *pgPixelStore) FindLastChangedBetween(ctx context.Context, from, to time.Time) ([]core.Area, error) {
	// ensure from < to
	if from.After(to) {
		from, to = to, from
	}

	// SELECT * FROM "tilechanges" WHERE ("last_changed" BETWEEN '2023-10-05 16:14:00+07' AND '2023-10-05 16:14:59+07');

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select(
		"canvas_id",
		"x",
		"y",
		"side",
		"last_changed",
	)
	sb.From("tilechanges")
	sb.Where(sb.Between("last_changed", from, to))

	query, args := sb.Build()
	fmt.Println(`query`, query)
	fmt.Println(`args`, args)

	rows, err := store.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var items []tilechange
	for rows.Next() {
		var item tilechange
		err := rows.Scan(&item.CanvasID, &item.X, &item.Y, &item.Side, &item.LastChanged)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	areas := mapTilechangesToAreas(items)
	slices.SortFunc(areas, core.SortAreasFn)
	return areas, nil
}

func (store *pgPixelStore) DeleteLastChangedForAreas(ctx context.Context, canvasID int64, side int64, areas ...core.Area) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom("tilechanges")
	db.Where(
		db.And(
			db.Equal("canvas_id", canvasID),
			db.Equal("side", side),
		),
	)

	for _, area := range areas {
		db.Or(
			db.And(
				db.Between("x", area.Min.X, area.Max.X),
				db.Between("y", area.Min.Y, area.Max.Y),
			),
		)
	}

	query, args := db.Build()
	fmt.Println(`query`, query)
	fmt.Println(`args`, args)

	_, err := store.db.ExecContext(ctx, query, args...)
	return err
}

func (store *pgPixelStore) SetLastChangedForAreas(ctx context.Context, canvasID int64, side int64, areas ...core.Area) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("tilechanges")
	ib.Cols("canvas_id", "x", "y", "side", "last_changed")

	for _, area := range areas {
		ib.Values(canvasID, area.Min.X, area.Min.Y, side, "NOW()")
	}

	ib.SQL(`
		ON CONFLICT (canvas_id, x, y, side) DO UPDATE SET 
			last_changed = EXCLUDED.last_changed
	`)

	query, args := ib.Build()
	_, err := store.db.ExecContext(ctx, query, args...)
	return err
}

func (store *pgPixelStore) SetLastChangedForPoints(ctx context.Context, canvasID int64, side int64, points ...core.Point) error {
	changedAreas := core.GetTileAreasFromPoints(side, points...)
	if len(changedAreas) == 0 {
		return nil
	}
	return store.SetLastChangedForAreas(ctx, canvasID, side, changedAreas...)
}

func NewPG() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres dbname=draw sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func chunkSlice[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
