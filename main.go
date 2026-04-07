package main

import (
	"image"
	"log"

	"rpg/src/characters"

	"github.com/hajimehoshi/ebiten/v2"
)

func CheckCollisionHorizontal(sprite *characters.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - 16.0
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *characters.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(image.Rect(int(sprite.X), int(sprite.Y), int(sprite.X)+16, int(sprite.Y)+16)) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - 16.0
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 200
}

func main() {
	ebiten.SetWindowSize(640, 400)
	ebiten.SetWindowTitle("Hello, World!")
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// NewImageFromFile return *ebiten.Image, then go's image.Image, and then an error
	// we use the ebiten image to actually render the thing, image.Image is the image data, and then error if not found
	// in the run game we insert the player image we got from new iamge from file into the actual RunGame struct
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
