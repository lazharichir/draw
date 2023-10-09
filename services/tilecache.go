package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lazharichir/draw/core"
	"github.com/lazharichir/draw/utils"
)

// Create a function that generates updated cached tiles

type TileCache struct {
	bucketName string
	s3         *awss3.Client
}

func NewTileCache(s3 *awss3.Client, bucketName string) *TileCache {
	return &TileCache{s3: s3, bucketName: bucketName}
}

func (cache *TileCache) PutTile(ctx context.Context, canvasID int64, tile core.Tile, img image.Image) error {
	x := tile.GetMinX()
	y := tile.GetMinY()
	side := tile.Width
	_, _, _ = x, y, side

	if !tile.IsSquare() {
		return errors.New("tile is not a square")
	}

	// Create a new buffer.
	buf, err := utils.ConvertImageToBytes(img)
	if err != nil {
		return err
	}

	// Put the object.
	putObjectParams := &awss3.PutObjectInput{
		Bucket: &cache.bucketName,
		Key:    utils.Ptr(tile.ObjectNameWithExt("png")),
		Body:   bytes.NewReader(buf),
	}

	fmt.Println("putObjectParams", len(buf), cache.bucketName, tile.ObjectNameWithExt("png"))

	putObjectOutput, err := cache.s3.PutObject(ctx, putObjectParams)
	if err != nil {
		return err
	}

	_ = putObjectOutput
	return nil
}

func (cache *TileCache) GetTile(ctx context.Context, canvasID int64, area core.Area) (image.Image, error) {
	bucket := cache.bucketName
	key := area.ObjectNameWithExt("png")

	getObjectParams := &awss3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	fmt.Println("getObjectParams", getObjectParams)
	getObjectOutput, err := cache.s3.GetObject(ctx, getObjectParams)
	if err != nil {
		return nil, err
	}
	defer getObjectOutput.Body.Close()

	img, _, err := image.Decode(getObjectOutput.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (cache *TileCache) DeleteTile(ctx context.Context, canvasID int64, tile core.Tile) error {
	bucket := cache.bucketName
	key := tile.ObjectNameWithExt("png")

	deleteObjectParams := &awss3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	fmt.Println("deleteObjectParams", deleteObjectParams)

	_, err := cache.s3.DeleteObject(ctx, deleteObjectParams)
	if err != nil {
		return err
	}

	return nil
}
