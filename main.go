package main

import (
	// "fmt"
	"image"
	"log"

	"rpg/src/characters"
	"rpg/src/environment"
	// "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Player       *characters.Player
	enemies      []*characters.Enemy
	potions      []*characters.Potion
	tilemapJSON  *environment.TilemapJSON
	tilemapImage *ebiten.Image
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.Y += 2
	}

	// for _, potion := range g.potions {
	// 	if g.Player.X < potion.X {
	// 		g.Player.health += potion.AmtHeal
	// 		fmt.Println(g.Player.health)
	// 	}
	// }

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			if enemy.X > g.Player.X {
				enemy.X -= 1
			} else if enemy.X < g.Player.X {
				enemy.X += 1
			}

			if enemy.Y > g.Player.Y {
				enemy.Y -= 1
			} else if enemy.Y < g.Player.Y {
				enemy.Y += 1
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := &ebiten.DrawImageOptions{}

	// loop over the layers

	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22

			srcX *= 16
			srcY *= 16

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImage.SubImage(image.Rect(srcX, srcY, srcX+16, srcY+16)).(*ebiten.Image),
				opts,
			)

			opts.GeoM.Reset()
		}
	}

	opts.GeoM.Translate(g.Player.X, g.Player.Y)
	// drawing our player
	screen.DrawImage(g.Player.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), opts)
	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), opts)
		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()

	for _, sprite := range g.potions {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), opts)
		opts.GeoM.Reset()
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
	PlayerImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/demon.png")
	if err != nil {
		log.Fatal(err)
	}

	PotionImg, _, err := ebitenutil.NewImageFromFile("assets/images/items/LifePot.png")
	if err != nil {
		log.Fatal(err)
	}

	EnemyImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/tilesets/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := environment.NewTilemapJSON("assets/resources/maps/spawn.json")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		Player: &characters.Player{
			Sprite: &characters.Sprite{
				Img: PlayerImg,
				X:   160,
				Y:   100,
			},
			Health: 100,
		},

		enemies: []*characters.Enemy{
			{
				Sprite: &characters.Sprite{
					Img: EnemyImg,
					X:   80,
					Y:   100,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &characters.Sprite{
					Img: EnemyImg,
					X:   80,
					Y:   100,
				},
				FollowsPlayer: false,
			},
		},

		potions: []*characters.Potion{
			{
				Sprite: &characters.Sprite{
					Img: PotionImg,
					X:   20,
					Y:   20,
				},
				AmtHeal: 10,
			},
		},
		tilemapJSON:  tilemapJSON,
		tilemapImage: tilemapImg,
	}
	// in the run game we insert the player image we got from new iamge from file into the actual RunGame struct
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
