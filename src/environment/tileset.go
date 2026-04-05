package environment

import (
	"encoding/json"
	"errors"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type DynamicTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (uTiles *UniformTileset) Img(id int) *ebiten.Image {
	id -= uTiles.gid

	srcX := (id) % 22
	srcY := (id) / 22

	srcX *= 16
	srcY *= 16

	return uTiles.img.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image)
}

type DynamicTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (dTileset *DynamicTileset) Img(id int) *ebiten.Image {
	id -= dTileset.gid

	return dTileset.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if strings.Contains(path, "buildings") {
		var dTilesetJSON DynamicTilesetJSON
		err = json.Unmarshal(contents, &dTilesetJSON)
		if err != nil {
			return nil, err
		}

		dTileset := DynamicTileset{}
		dTileset.gid = gid
		dTileset.imgs = make([]*ebiten.Image, 0)

		for _, tileJSON := range dTilesetJSON.Tiles {
			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)
			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}

			dTileset.imgs = append(dTileset.imgs, img)
		}

		return &dTileset, nil

	} else if strings.Contains(path, "TilesetFloor") {
		var uTilesetJSON UniformTilesetJSON
		err = json.Unmarshal(contents, &uTilesetJSON)
		if err != nil {
			return nil, err
		}

		uTileset := UniformTileset{}
		tileJSONPath := uTilesetJSON.Path
		tileJSONPath = filepath.Clean(tileJSONPath)
		tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
		tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
		tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
		tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
		tileJSONPath = filepath.Join("assets/", tileJSONPath)
		img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
		if err != nil {
			return nil, err
		}
		uTileset.img = img
		uTileset.gid = gid

		return &uTileset, nil

	}

	return nil, errors.New("could not identify the type of tileset")
}
