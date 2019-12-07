package spt

// http://iquilezles.org/www/articles/distfunctions/distfunctions.htm

import (
	"encoding/gob"
	"math"
)

func init() {
	gob.Register(SDFExtrude{})
	gob.Register(SDFRevolve{})
	gob.Register(SDFSphere{})
	gob.Register(SDFCube{})
	gob.Register(SDFTorus{})
	gob.Register(SDFCone{})
	gob.Register(SDFRounded{})
	gob.Register(SDFHollow{})
	gob.Register(SDFElongate{})
	gob.Register(SDFRepeat{})
}

type SDF3 interface {
	SDF() func(Vec3) float64
	Sphere() (Vec3, float64)
}

func SDF3Normal(sdf SDF3, p Vec3) Vec3 {
	distance := sdf.SDF()
	step := 0.000001
	gradient := Vec3{
		distance(Vec3{p.X + step, p.Y, p.Z}) - distance(Vec3{p.X - step, p.Y, p.Z}),
		distance(Vec3{p.X, p.Y + step, p.Z}) - distance(Vec3{p.X, p.Y - step, p.Z}),
		distance(Vec3{p.X, p.Y, p.Z + step}) - distance(Vec3{p.X, p.Y, p.Z - step}),
	}
	return gradient.Unit()
}

type SDFExtrude struct {
	H float64
	SDF2
}

func (s SDFExtrude) SDF() func(Vec3) float64 {
	sdf := s.SDF2.SDF()
	return func(pos Vec3) float64 {
		d := sdf(V2(pos.X, pos.Y))
		w := V2(d, abs(pos.Z)-s.H)
		return min(max(w.X, w.Y), 0.0) + len2(max2(w, Zero2))
	}
}

func (s SDFExtrude) Sphere() (Vec3, float64) {
	center, radius := s.SDF2.Circle()
	return Zero3.Add(V3(center.X, center.Y, 0)), sqrt(radius*radius + s.H*s.H)
}

func Extrude(h float64, sdf SDF2) SDF3 {
	return SDFExtrude{h / 2, sdf}
}

type SDFRevolve struct {
	O float64
	SDF2
}

func (s SDFRevolve) SDF() func(Vec3) float64 {
	sdf := s.SDF2.SDF()
	return func(pos Vec3) float64 {
		return sdf(V2(len2(V2(pos.X, pos.Z))-s.O, pos.Y))
	}
}

func (s SDFRevolve) Sphere() (Vec3, float64) {
	center, radius := s.SDF2.Circle()
	return Zero3.Add(V3(center.X, center.Y, 0)), radius
}

func Revolve(o float64, sdf SDF2) SDF3 {
	return SDFRevolve{o, sdf}
}

type SDFSphere struct {
	R float64
}

func (s SDFSphere) SDF() func(Vec3) float64 {
	return func(pos Vec3) float64 {
		return len3(pos) - s.R
	}
}

func (s SDFSphere) Sphere() (Vec3, float64) {
	return Zero3, s.R
}

func Sphere(r float64) SDF3 {
	return SDFSphere{r}
}

type SDFCube struct {
	X, Y, Z float64
}

func (s SDFCube) SDF() func(Vec3) float64 {
	box := V3(s.X, s.Y, s.Z)
	return func(pos Vec3) float64 {
		q := sub3(abs3(pos), box)
		return len3(max3(q, Zero3)) + min(max(q.X, max(q.Y, q.Z)), 0.0)
	}
}

func (s SDFCube) Sphere() (Vec3, float64) {
	return Zero3, len3(V3(s.X, s.Y, s.Z))
}

func Cube(x, y, z float64) SDF3 {
	return SDFCube{x / 2, y / 2, z / 2}
}

func CubeR(x, y, z, r float64) SDF3 {
	d := r * 2
	return Round(r, Cube(x-d, y-d, z-d))
}

func Cylinder(h, r float64) SDF3 {
	return Extrude(h, Circle(r))
}

func CylinderR(h, r, ro float64) SDF3 {
	return Round(ro, Cylinder(h-ro*2, r-ro))
}

type SDFTorus struct {
	V Vec2
}

func (s SDFTorus) SDF() func(Vec3) float64 {
	return func(pos Vec3) float64 {
		q := V2(len2(V2(pos.X, pos.Z))-s.V.X, pos.Y)
		return len2(q) - s.V.Y
	}
}

func (s SDFTorus) Sphere() (Vec3, float64) {
	return Zero3, s.V.X + s.V.Y
}

func Torus(x, y float64) SDF3 {
	w := x - y
	return SDFTorus{Vec2{x - w/2, w / 2}}
}

type SDFCone struct {
	X, Y, H, R float64
}

func (s SDFCone) SDF() func(Vec3) float64 {
	return func(pos Vec3) float64 {
		q := V2(len2(V2(pos.X, pos.Y)), pos.Z)
		d1 := -pos.Z - s.H
		d2 := max(dot2(q, V2(s.X, s.Y)), pos.Z)
		return len2(max2(V2(d1, d2), Zero2)) + min(max(d1, d2), 0.0)
	}
}

func (s SDFCone) Sphere() (Vec3, float64) {
	return Zero3, sqrt(s.H*s.H + s.R*s.R)
}

func Cone(h, r float64) SDF3 {
	rad := math.Atan(h / r)
	return Translate(V3(0, 0, h/2), SDFCone{math.Sin(rad), math.Cos(rad), h, r})
}

func TriPrism(h, w float64) SDF3 {
	return Translate(V3(0, 0, -h/2), Rotate(V3(1, 0, 0), -90, Extrude(w, Triangle(
		V2(0, h),
		V2(-w/2, 0),
		V2(w/2, 0),
	))))
}

func Pyramid(h, w float64) SDF3 {
	prism := TriPrism(h, w)
	return Intersection(prism,
		Rotate(V3(0, 0, 1), 90, prism),
	)
}

type SDFRounded struct {
	Radius float64
	SDF3
}

func (s SDFRounded) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	return func(pos Vec3) float64 {
		return sdf(pos) - s.Radius
	}
}

func (s SDFRounded) Sphere() (Vec3, float64) {
	center, radius := s.SDF3.Sphere()
	return center, radius + s.Radius
}

func Round(radius float64, sdf SDF3) SDF3 {
	return SDFRounded{radius, sdf}
}

type SDFHollow struct {
	Thickness float64
	SDF3
}

func (s SDFHollow) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	return func(pos Vec3) float64 {
		return abs(sdf(pos)) - s.Thickness
	}
}

func (s SDFHollow) Sphere() (Vec3, float64) {
	center, radius := s.SDF3.Sphere()
	return center, radius
}

func Hollow(thickness float64, sdf SDF3) SDF3 {
	return SDFHollow{thickness, sdf}
}

type SDFElongate struct {
	H Vec3
	SDF3
}

func (s SDFElongate) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	return func(pos Vec3) float64 {
		return sdf(sub3(pos, clamp3(pos, neg3(s.H), s.H)))
	}
}

func (s SDFElongate) Sphere() (Vec3, float64) {
	center, radius := s.SDF3.Sphere()
	return center, radius + s.H.Length()
}

func Elongate(v Vec3, sdf SDF3) SDF3 {
	return SDFElongate{v, sdf}
}

type SDFRepeat struct {
	Step Vec3
	SDF3
}

func (s SDFRepeat) SDF() func(Vec3) float64 {
	sdf := s.SDF3.SDF()
	return func(pos Vec3) float64 {
		hstep := scale3(s.Step, 0.5)
		return sdf(sub3(mod3(add3(pos, hstep), s.Step), hstep))
	}
}

func (s SDFRepeat) Sphere() (Vec3, float64) {
	center, _ := s.SDF3.Sphere()
	return center, math.MaxFloat64
}

func Repeat(step Vec3, sdf SDF3) SDF3 {
	return SDFRepeat{step, sdf}
}
