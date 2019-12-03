package spt

import (
	"testing"
)

func partRender(part Thing) {

	scene := Scene{
		Width:     1280,
		Height:    720,
		Passes:    10,
		Samples:   1,
		Bounces:   4,
		Horizon:   10000,
		Threshold: 0.0001,
		Ambient:   Black,
		Sky:       White.Scale(0.05),

		Camera: NewCamera(
			V3(0, -2000, 2000),
			V3(0, 0, 0),
			Z3,
			20,
			Zero3,
			0.0,
		),

		Stuff: []Thing{
			//WorkBench(25000),
			Object(
				Light(White.Scale(4)),
				Translate(V3(-2500, 0, 5000), Sphere(2500)),
			),
			part,
		},
	}

	Render("test.png", scene, nil)
}

func TestGearWheel(t *testing.T) {
	partRender(
		Object(
			Steel,
			GearWheel(),
		),
	)
}
