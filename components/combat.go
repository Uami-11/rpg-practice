// Package components...
package components

type Combat interface {
	Health() int
	AttackPower() int
	Attacking() bool
	Attack() bool
	Update()
	Damage(amount int)
}

type BasicCombat struct {
	health      int
	attackPower int
	attacking   bool
}

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
		return true
	}
	return false
}

func (e *EnemyCombat) Update() {
	e.timeSinceAttack++
}

func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return b.attacking
}

func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func NewBasicCombat(health, attackPow int) *BasicCombat {
	return &BasicCombat{
		health:      health,
		attackPower: attackPow,
		attacking:   false,
	}
}

func NewEnemyCombat(health, attackPow, attackCooldown int) *EnemyCombat {
	return &EnemyCombat{
		BasicCombat:    NewBasicCombat(health, attackPow),
		attackCooldown: attackCooldown,
	}
}

func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Health() int {
	return b.health
}

func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

func (b *BasicCombat) Update() {
}

var (
	_ Combat = (*BasicCombat)(nil)
	_ Combat = (*EnemyCombat)(nil)
)
