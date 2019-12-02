package spt

import (
	"math"
)

type Camera struct {
	Origin   Vec3
	U        Vec3
	V        Vec3
	W        Vec3
	M        float64
	Focus    float64
	Aperture float64
}

func NewCamera(lookFrom Vec3, lookAt Vec3, vup Vec3, vfov float64, focus Vec3, aperture float64) Camera {
	c := Camera{}
	c.Origin = lookFrom
	c.W = lookAt.Sub(lookFrom).Unit()
	c.U = vup.Cross(c.W).Unit()
	c.V = c.W.Cross(c.U).Unit()
	c.M = 1 / math.Tan(vfov*math.Pi/360)
	c.Focus = focus.Sub(c.Origin).Length()
	c.Aperture = aperture
	return c
}

func (c Camera) CastRay(imageX, imageY, imageW, imageH int, jitterU, jitterV float64, rnd Random) Ray {

	aspect := float64(imageW) / float64(imageH)
	px := ((float64(imageX)+jitterU-0.5)/(float64(imageW)-1))*2 - 1
	py := ((float64(imageY)+jitterV-0.5)/(float64(imageH)-1))*2 - 1

	direction := c.W.Scale(c.M).
		Add(c.U.Scale(-px * aspect)).
		Add(c.V.Scale(-py)).
		Unit()

	origin := c.Origin

	if c.Aperture > 0 {
		focus := c.Origin.Add(direction.Scale(c.Focus))
		angle := rnd.Float64() * 2 * math.Pi
		radius := rnd.Float64() * c.Aperture

		origin = origin.
			Add(c.U.Scale(math.Cos(angle) * radius)).
			Add(c.V.Scale(math.Sin(angle) * radius))

		direction = focus.Sub(origin).Unit()
	}

	return Ray{origin, direction, rnd}
}
