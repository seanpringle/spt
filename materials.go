package spt

import (
	"encoding/gob"
	"math"
)

func init() {
	gob.Register(Nothing{})
	gob.Register(Diffuse{})
	gob.Register(Emitter{})
	gob.Register(Metallic{})
	gob.Register(Dielectric{})
}

var (
	Steel     = Metal(Color{0.4, 0.4, 0.4}, 0.95)
	Stainless = Metal(Color{0.4, 0.4, 0.4}, 0.3)
	Gold      = Metal(Color{0.93, 0.78, 0.31}, 0.0)
	Copper    = Metal(Color{0.68, 0.45, 0.41}, 0.8)
	Brass     = Metal(Color{0.80, 0.58, 0.45}, 0.9)
)

// choose a vec3 less than unit
func pickVec3(rnd Random) Vec3 {
	for {
		v := Vec3{
			rnd.Float64()*2 - 1,
			rnd.Float64()*2 - 1,
			rnd.Float64()*2 - 1,
		}
		if v.Dot(v) < 1.0 {
			return v
		}
	}
}

type Material interface {
	Light() (Color, bool)
	Scatter(Ray, *Thing, Vec3) (Ray, Color, bool)
}

type Nothing struct{}

func (mat Nothing) Light() (Color, bool) {
	return Nought, false
}

func (mat Nothing) Scatter(r Ray, thing *Thing, hit Vec3) (Ray, Color, bool) {
	return Ray{}, Nought, false
}

type Emitter struct {
	Nothing
	Color
}

func (mat Emitter) Light() (Color, bool) {
	return mat.Color, true
}

func Light(c Color) Material {
	return Emitter{Color: c}
}

type Diffuse struct {
	Nothing
	Color
}

func (mat Diffuse) Scatter(r Ray, thing *Thing, hit Vec3) (Ray, Color, bool) {
	normal := thing.Normal(hit)
	redirection := hit.Add(normal).Add(pickVec3(r.rnd)).Sub(hit).Unit()
	bounced := Ray{hit, redirection, r.rnd}
	return bounced, mat.Color, true
}

func Matt(c Color) Material {
	return Diffuse{Color: c}
}

type Metallic struct {
	Nothing
	Color
	Roughness float64
}

func (mat Metallic) Scatter(r Ray, thing *Thing, hit Vec3) (Ray, Color, bool) {
	normal := thing.Normal(hit)
	reflected := r.Direction.Unit().Reflect(normal)

	if mat.Roughness > 0 {
		reflected = reflected.
			Add(normal).
			Add(pickVec3(r.rnd).
				Sub(normal).
				Scale(mat.Roughness)).
			Unit()
	}

	if reflected.Dot(normal) > 0 {
		return Ray{hit, reflected, r.rnd}, mat.Color, true
	}

	return Ray{}, Nought, false
}

func Metal(c Color, roughness float64) Material {
	return Metallic{Color: c, Roughness: roughness}
}

type Dielectric struct {
	Nothing
	Color
	RefractiveIndex float64
}

func schlick(cosine float64, refInd float64) float64 {
	r0 := (1.0 - refInd) / (1.0 + refInd)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow(1.0-cosine, 5)
}

func (mat Dielectric) Scatter(r Ray, thing *Thing, hit Vec3) (Ray, Color, bool) {
	normal := thing.Normal(hit)

	outwardNormal := normal
	niOverNt := 1.0 / mat.RefractiveIndex
	cosine := -r.Direction.Unit().Dot(normal) / r.Direction.Length()

	if r.Direction.Dot(normal) > 0 {
		outwardNormal = normal.Neg()
		niOverNt = mat.RefractiveIndex
		cosine = mat.RefractiveIndex * r.Direction.Unit().Dot(normal) / r.Direction.Length()
	}

	direction := r.Direction.Unit().Reflect(normal)

	if r.rnd.Float64() >= schlick(cosine, mat.RefractiveIndex) {
		if refracted, was := r.Direction.Refract(outwardNormal, niOverNt); was {
			direction = refracted
		}
	}

	return Ray{hit, direction, r.rnd}, mat.Color, true
}

func Glass(color Color, refInd float64) Material {
	return Dielectric{Color: color, RefractiveIndex: refInd}
}
