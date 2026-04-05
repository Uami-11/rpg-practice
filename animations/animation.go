// Package animations..
package animations

type Animation struct {
	First        int
	Last         int
	Step         int // how many indeces to moves per frame
	SpeedInTps   float32
	FrameCounter float32
	frame        int
}

func (anim *Animation) Update() {
	anim.FrameCounter -= 1.0
	if anim.FrameCounter < 0.0 {
		anim.FrameCounter = anim.SpeedInTps
		anim.frame += anim.Step
		if anim.frame > anim.Last {
			anim.frame = anim.First
		}
	}
}

func (anim *Animation) Frame() int {
	return anim.frame
}

func NewAnimation(first, last, step int, speedTps float32) *Animation {
	return &Animation{
		first, last, step, speedTps, speedTps, first,
	}
}
