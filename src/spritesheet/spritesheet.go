// Package spritesheet...
package spritesheet

import "image"

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	Tilesize      int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := (index % s.WidthInTiles) * s.Tilesize
	y := (index / s.HeightInTiles) * s.Tilesize

	return image.Rect(x, y, x+s.Tilesize, y+s.Tilesize)
}

func NewSpriteSheet(width, height, tileSize int) *SpriteSheet {
	return &SpriteSheet{
		WidthInTiles:  width,
		HeightInTiles: height,
		Tilesize:      tileSize,
	}
}
