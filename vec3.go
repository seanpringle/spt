package spt

import (
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Scale(t float64) Vec3 {
	return Vec3{X: v.X * t, Y: v.Y * t, Z: v.Z * t}
}

func (v Vec3) Mul(v2 Vec3) Vec3 {
	return Vec3{X: v.X * v2.X, Y: v.Y * v2.Y, Z: v.Z * v2.Z}
}

func (v Vec3) Div(v2 Vec3) Vec3 {
	return Vec3{X: v.X / v2.X, Y: v.Y / v2.Y, Z: v.Z / v2.Z}
}

func (v Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{X: v.X - v2.X, Y: v.Y - v2.Y, Z: v.Z - v2.Z}
}

func (v Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Unit() Vec3 {
	return v.Scale(1.0 / v.Length())
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) Abs() Vec3 {
	return Vec3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)}
}

func (v Vec3) Min(v1 Vec3) Vec3 {
	return Vec3{X: math.Min(v.X, v1.X), Y: math.Min(v.Y, v1.Y), Z: math.Min(v.Z, v1.Z)}
}

func (v Vec3) Max(v1 Vec3) Vec3 {
	return Vec3{X: math.Max(v.X, v1.X), Y: math.Max(v.Y, v1.Y), Z: math.Max(v.Z, v1.Z)}
}

func (v Vec3) Refract(n Vec3, niOverNt float64) (Vec3, bool) {
	uv := v.Unit()
	un := n.Unit()
	dt := uv.Dot(un)

	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		refracted := uv.Sub(un.Scale(dt)).Scale(niOverNt).Sub(un.Scale(math.Sqrt(discriminant)))
		return refracted, true
	}

	return Zero3, false
}

func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.Scale(2.0 * v.Dot(n)))
}

func (v1 Vec3) Dot(v2 Vec3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{v1.Y*v2.Z - v1.Z*v2.Y, -(v1.X*v2.Z - v1.Z*v2.X), v1.X*v2.Y - v1.Y*v2.X}
}
