package spt

import (
	"sync"
	"testing"
	"time"
)

func testScene() Scene {

	steel := Metal(Color{0.4, 0.4, 0.4}, 0.95)
	stainless := Metal(Color{0.4, 0.4, 0.4}, 0.3)
	gold := Metal(Color{0.93, 0.78, 0.31}, 0.0)
	copper := Metal(Color{0.68, 0.45, 0.41}, 0.8)
	brass := Metal(Color{0.80, 0.58, 0.45}, 0.9)

	things := func(material Material, origin Vec3) []Thing {

		at := func(v Vec3, s SDF3) SDF3 {
			return Translate(origin.Add(v), s)
		}

		return append(
			[]Thing{
				Object(material,
					at(V3(0, 0, 500), Sphere(500)),
				),
				Object(material,
					at(V3(0, 2000, 500), Cube(1000, 1000, 1000)),
				),
				Object(material,
					at(V3(-1500, 0, 500), Cylinder(1000, 500)),
				),
				Object(material,
					at(V3(1500, 2000, 500), Cone(1000, 500)),
				),
				Object(material,
					at(V3(1500, 0, 500), Torus(500, 350)),
				),
				Object(material,
					at(V3(-1500, 2000, 500), Pyramid(1000, 1000)),
				),
			}, func() (set []Thing) {
				for i := 3; i <= 8; i++ {
					set = append(set, Object(material,
						at(V3(float64(i-3)*600-1450, -1250, 250), Extrude(500, Polygon(i, 250))),
					))
				}
				return
			}()...,
		)
	}

	stuff := []Thing{
		WorkBench(25000),
		Object(
			Light(White.Scale(4)),
			Translate(V3(-7500, 0, 20000), Sphere(10000)),
		),
		Object(
			brass,
			Translate(V3(-5500, 0, 500),
				Difference(
					Intersection(
						Sphere(500),
						Cube(900, 900, 900),
					),
					Cylinder(1002, 200),
					Rotate(V3(1, 0, 0), 90, Cylinder(1002, 200)),
					Rotate(V3(0, 1, 0), 90, Cylinder(1002, 200)),
				),
			),
		),
		Object(
			brass,
			Translate(V3(5500, 0, 500),
				Union(
					Intersection(
						Sphere(500),
						Cube(900, 900, 900),
					),
					Cylinder(1250, 200),
					Rotate(V3(1, 0, 0), 90, Cylinder(1250, 200)),
					Rotate(V3(0, 1, 0), 90, Cylinder(1250, 200)),
				),
			),
		),
		Object(
			brass,
			Translate(V3(6500, 3500, 500),
				Round(100, Cube(800, 800, 800)),
			),
		),
		Object(
			brass,
			Translate(V3(-6500, 3500, 500),
				Round(100, Cylinder(800, 500)),
			),
		),
	}

	stuff = append(stuff, things(steel, V3(-2500, -2750, 0))...)
	stuff = append(stuff, things(copper, V3(2500, -2750, 0))...)
	stuff = append(stuff, things(stainless, V3(-2500, 2250, 0))...)
	stuff = append(stuff, things(gold, V3(2500, 2250, 0))...)

	for i := 0; i < 9; i++ {
		color := Color{0.5, 0.5, 1.0}
		if i%2 == 0 {
			color = Color{0.5, 1.0, 0.5}
		}
		stuff = append(stuff, Object(
			Glass(color, 1.5),
			Translate(V3(0, float64(i)*1000-3500, 250), Sphere(250)),
		))
	}

	return Scene{
		Width:      1280,
		Height:     720,
		Passes:     10,
		Samples:    1,
		Bounces:    4,
		Horizon:    100000,
		Threshold:  0.0001,
		Ambient:    White.Scale(0.05),
		Background: Transparent,

		Camera: NewCamera(
			V3(0, -8000, 8000),
			V3(0, -1000, 500),
			Z3,
			40,
			Zero3,
			0.0,
		),

		Stuff: stuff,
	}
}

func TestLocal(t *testing.T) {
	Render("test.png", testScene(), []Renderer{NewLocalRenderer()})
}

func TestRPC(t *testing.T) {
	var group sync.WaitGroup
	stop := make(chan struct{})

	group.Add(1)
	go func() {
		RenderServeRPC(stop, 34242)
		group.Done()
	}()

	time.Sleep(time.Second)
	Render("test.png", testScene(), []Renderer{NewRPCRenderer("127.0.0.1:34242")})

	close(stop)
	group.Wait()
}
