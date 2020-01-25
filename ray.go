package spt

import (
	"math"
)

type Ray struct {
	Origin    Vec3
	Direction Vec3
	rnd       Random
}

type Hit struct {
	Thing    *Thing
	Position Vec3
	Normal   Vec3
}

func (r Ray) PathTrace(scene *Scene, depth int, bypass *Thing) (Color, int, float64) {

	var (
		shadow      Ray
		color       Color
		scolor      Color
		attenuation Color
		scattered   bool
		alpha       float64
		bounces     int
		invisible   bool
	)

	alpha = 1.0

	if thing, hit := r.march(scene, bypass); thing != nil {
		_, invisible = thing.Material().(Invisible)

		if depth == 0 && invisible {
			alpha = scene.ShadowH
			if r.directLight(scene, hit).Brightness() > 0.0 {
				// an invisible material under direct light for a primary ray cannot catch shadows
				alpha = 0.0
			}
		}

		if alpha > 0.0 && depth < scene.Bounces {

			if shadow, attenuation, scattered = thing.Material().Scatter(r, thing, hit, depth); scattered {
				scolor, bounces, _ = shadow.PathTrace(scene, depth+1, thing)

				// invisible surfaces pass through ambient lighting for non-primary ray hits
				if depth > 0 && invisible {
					attenuation = White
					scolor = scene.Ambient
				}

				color = color.Add(attenuation.Mul(scolor))

				// invisible surfaces use nested shadow ray colors only to weight their own shadow alpha,
				// giving a soft-shadow prenumbra effect similar to real shadows on normal materials
				if depth == 0 && invisible {
					alpha = math.Min(scene.ShadowH, math.Max(0.0, (alpha-(scolor.Brightness()/scene.ShadowD))))
					color = attenuation
				}
			}

			if light, is := thing.Material().Light(); is {
				color = color.Add(light)
			}

			bounces++
		}

		return color, bounces, alpha
	}

	if depth > 0 {
		return scene.Ambient, 0, 1.0
	}

	return Naught, 0, 0.0
}

// ray marching by sphere tracing
func (r Ray) march(scene *Scene, bypass *Thing) (*Thing, Vec3) {

	pos := r.Origin

	// refracted rays bypass one object, then act normally
	for bypass != nil {
		dist := bypass.Distance(pos)
		if dist > 0 {
			break
		}
		dist = math.Max(math.Abs(dist), scene.Threshold)
		pos = pos.Add(r.Direction.Scale(dist))
	}

	// shadow acne
	pos = pos.Add(r.Direction.Scale(scene.Threshold * 10))

	// find all possible targets
	var targets []*Thing
	for i := range scene.Stuff {
		t := &scene.Stuff[i]
		center, radius := t.Sphere()

		// since objects have a bounding sphere, behave like a non-SDF
		// ray tracer and do a line-sphere intersection test to quickly
		// rule them in or out
		to := r.Origin.Sub(center)
		b := to.Dot(r.Direction)
		c := to.Dot(to) - radius*radius
		d := b*b - c
		if d > 0 {
			d = sqrt(d)
			t1 := -b - d
			if t1 > scene.Threshold {
				targets = append(targets, t)
				continue
			}
			t2 := -b + d
			if t2 > scene.Threshold {
				targets = append(targets, t)
				continue
			}
		}
	}

	if len(targets) > 0 {

		// now behave like a path tracer and evaluate SDFs
		for pos.Sub(r.Origin).Length() < scene.Horizon {

			near := (*Thing)(nil)
			dist := 0.0

			for _, t := range targets {
				if near != nil {
					bd := t.BoundingDistance(pos)
					if bd > dist {
						continue
					}
				}
				d := t.Distance(pos)
				if d < dist || near == nil {
					near = t
					dist = d
					continue
				}
			}

			if dist < scene.Threshold {
				return near, pos
			}

			pos = pos.Add(r.Direction.Scale(dist))
		}
	}

	return nil, Z3
}

// direct sample all lights
func (r Ray) directLight(scene *Scene, pos Vec3) Color {
	var color Color
	for i := range scene.Stuff {
		t := &scene.Stuff[i]
		if light, is := t.Material().Light(); is {
			center, radius := t.Sphere()
			center = center.Add(pickVec3(r.rnd).Scale(radius * scene.ShadowR))
			lr := Ray{pos, center.Sub(pos).Unit(), r.rnd}
			if lt, _ := lr.march(scene, nil); lt == t {
				color = color.Add(light)
			}
		}
	}
	return color
}
