package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"rpg/animations"
	"rpg/components"
	"rpg/constants"
	"rpg/src/characters"
	"rpg/src/core"
	"rpg/src/environment"
	"rpg/src/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Player            *characters.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*characters.Enemy
	potions           []*characters.Potion
	tilemapJSON       *environment.TilemapJSON
	tilemapImage      *ebiten.Image
	tilesets          []environment.Tileset
	camera            *core.Camera
	colliders         []image.Rectangle
}

func NewGame() *Game {
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

	tilesets, err := tilemapJSON.GenTilesets()
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)

	game := &Game{
		Player: &characters.Player{
			Sprite: &characters.Sprite{
				Img: PlayerImg,
				X:   10,
				Y:   100,
			},
			Health: 3,
			Animations: map[characters.PlayerState]*animations.Animation{
				characters.Up:    animations.NewAnimation(5, 13, 4, 20.0),
				characters.Down:  animations.NewAnimation(4, 12, 4, 20.0),
				characters.Left:  animations.NewAnimation(6, 14, 4, 20.0),
				characters.Right: animations.NewAnimation(7, 15, 4, 20.0),
			},
			CombatComponent: components.NewBasicCombat(3, 1),
		},

		playerSpriteSheet: playerSpriteSheet,

		enemies: []*characters.Enemy{
			{
				Sprite: &characters.Sprite{
					Img: EnemyImg,
					X:   80,
					Y:   100,
				},
				FollowsPlayer:   true,
				CombatComponent: components.NewEnemyCombat(5, 1, 30),
			},
			{
				Sprite: &characters.Sprite{
					Img: EnemyImg,
					X:   80,
					Y:   100,
				},
				FollowsPlayer:   false,
				CombatComponent: components.NewEnemyCombat(5, 1, 30),
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
		tilesets:     tilesets,
		camera:       core.NewCamera(0.0, 0.0),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}
	return game
}

func (g *Game) Update() error {
	g.camera.FollowTarget(g.Player.X+8, g.Player.Y+8, 320, 200)
	g.camera.Constraint(
		float64(g.tilemapJSON.Layers[0].Width)*16,
		float64(g.tilemapJSON.Layers[0].Height)*16,
		320,
		200,
	)

	g.Player.Dx = 0.0
	g.Player.Dy = 0.0

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.Dx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.Dy = 2
	}

	g.Player.X += g.Player.Dx

	CheckCollisionHorizontal(g.Player.Sprite, g.colliders)

	g.Player.Y += g.Player.Dy

	CheckCollisionVertical(g.Player.Sprite, g.colliders)

	activeAnimation := g.Player.ActiveAnimation(int(g.Player.Dx), int(g.Player.Dy))
	if activeAnimation != nil {
		activeAnimation.Update()
	}

	clicked := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	clickX, clickY := ebiten.CursorPosition()
	deadEnemies := make(map[int]struct{})

	g.Player.CombatComponent.Update()
	prect := image.Rect(int(g.Player.X), int(g.Player.Y), int(g.Player.X)+constants.Tilesize, int(g.Player.Y)+constants.Tilesize)
	for index, enemy := range g.enemies {
		enemy.CombatComponent.Update()
		rect := image.Rect(int(enemy.X), int(enemy.Y), int(enemy.X)+constants.Tilesize, int(enemy.Y)+constants.Tilesize)

		if rect.Overlaps(prect) {
			if enemy.CombatComponent.Attack() {
				g.Player.CombatComponent.Damage(enemy.CombatComponent.AttackPower())
				fmt.Printf("Player has %d health remaining\n", g.Player.CombatComponent.Health())

				if g.Player.CombatComponent.Health() <= 0 {
					fmt.Printf("PLAYER DEAD\n")
				}
			}
		}

		// is cursor in rect?
		if clickX > rect.Min.X && clickX < rect.Max.X && clickY > rect.Min.Y && clickY < rect.Max.Y {
			if clicked {
				enemy.CombatComponent.Damage(g.Player.CombatComponent.AttackPower())

				if enemy.CombatComponent.Health() <= 0 {
					deadEnemies[index] = struct{}{}
				}
			}
		}
	}

	if len(deadEnemies) > 0 {
		newEnemies := make([]*characters.Enemy, 0)
		for index, enemy := range g.enemies {
			if _, exists := deadEnemies[index]; !exists {
				newEnemies = append(newEnemies, enemy)
			}
		}
		g.enemies = newEnemies
	}

	// for _, potion := range g.potions {
	// 	if g.Player.X < potion.X {
	// 		g.Player.health += potion.AmtHeal
	// 		fmt.Println(g.Player.health)
	// 	}
	// }

	for _, enemy := range g.enemies {
		enemy.Dx = 0.0
		enemy.Dy = 0.0
		if enemy.FollowsPlayer {
			if enemy.X > g.Player.X {
				enemy.Dx = -1
			} else if enemy.X < g.Player.X {
				enemy.Dx = 1
			}

			if enemy.Y > g.Player.Y {
				enemy.Dy = -1
			} else if enemy.Y < g.Player.Y {
				enemy.Dy = 1
			}
		}

		enemy.X += enemy.Dx
		CheckCollisionHorizontal(enemy.Sprite, g.colliders)
		enemy.Y += enemy.Dy
		CheckCollisionVertical(enemy.Sprite, g.colliders)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := &ebiten.DrawImageOptions{}

	// loop over the layers

	for layerIndex, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := index % layer.Width
			y := index / layer.Width

			x *= 16
			y *= 16

			img := g.tilesets[layerIndex].Img(id)

			opts.GeoM.Translate(float64(x), float64(y))

			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()) + 16))

			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(img, opts)

			opts.GeoM.Reset()
		}

		for _, collider := range g.colliders {
			vector.StrokeRect(screen,
				float32(collider.Min.X)+float32(g.camera.X),
				float32(collider.Min.Y)+float32(g.camera.Y),
				float32(collider.Dx()),
				float32(collider.Dy()),
				1.0,
				color.RGBA{
					255,
					0,
					0,
					255,
				},
				true)
		}
	}

	opts.GeoM.Translate(g.Player.X, g.Player.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)
	// drawing our player
	playerFrame := 0
	activeAnimation := g.Player.ActiveAnimation(int(g.Player.Dx), int(g.Player.Dy))
	if activeAnimation != nil {
		playerFrame = activeAnimation.Frame()
	}
	screen.DrawImage(g.Player.Img.SubImage(g.playerSpriteSheet.Rect(playerFrame)).(*ebiten.Image), opts)
	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), opts)
		opts.GeoM.Reset()
	}

	opts.GeoM.Reset()

	for _, sprite := range g.potions {
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(sprite.Img.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), opts)
		opts.GeoM.Reset()
	}
}
