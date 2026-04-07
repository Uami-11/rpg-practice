package characters

import (
	"rpg/animations"
	"rpg/components"
)

type PlayerState uint8

const (
	Down PlayerState = iota
	Up
	Left
	Right
)

type Player struct {
	*Sprite
	Health          uint
	Animations      map[PlayerState]*animations.Animation
	CombatComponent *components.BasicCombat
}

func (player *Player) ActiveAnimation(Dx, Dy int) *animations.Animation {
	if Dx > 0 {
		return player.Animations[Right]
	}
	if Dx < 0 {
		return player.Animations[Left]
	}
	if Dy < 0 {
		return player.Animations[Up]
	}
	if Dy > 0 {
		return player.Animations[Down]
	}

	return nil
}
