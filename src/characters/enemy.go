package characters

import "rpg/components"

type Enemy struct {
	*Sprite
	FollowsPlayer   bool
	CombatComponent *components.EnemyCombat
}
