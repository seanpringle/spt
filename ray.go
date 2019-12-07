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

func (r Ray) PathTrace(scene *Scene, depth int, bypass *Thing) Color {

	if thing, hit := r.march(scene, bypass); thing != nil {
		color := Nought

		if depth < scene.Bounces {

			if shadow, attenuation, does := thing.Material().Scatter(r, thing, hit); does {
				color = color.Add(attenuation.Mul(shadow.PathTrace(scene, depth+1, thing)))
			}

			if light, is := thing.Material().Light(); is {
				color = color.Add(light)
			}
		}

		return color
	}

	return scene.Ambient
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
