package spt

import (
	"math"
)

type Vec2 struct {
	X, Y float64
}

// Multiply by scalar
func (v Vec2) Scale(t float64) Vec2 {
	return Vec2{X: v.X * t, Y: v.Y * t}
}

func (v Vec2) Mul(v2 Vec2) Vec2 {
	return Vec2{X: v.X * v2.X, Y: v.Y * v2.Y}
}

func (v Vec2) Div(v2 Vec2) Vec2 {
	return Vec2{X: v.X / v2.X, Y: v.Y / v2.Y}
}

func (v Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{X: v.X - v2.X, Y: v.Y - v2.Y}
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{X: v.X + v2.X, Y: v.Y + v2.Y}
}

func (v Vec2) Abs() Vec2 {
	return Vec2{math.Abs(v.X), math.Abs(v.Y)}
}

func (v Vec2) Min(v1 Vec2) Vec2 {
	return Vec2{X: math.Min(v.X, v1.X), Y: math.Min(v.Y, v1.Y)}
}

func (v Vec2) Max(v1 Vec2) Vec2 {
	return Vec2{X: math.Max(v.X, v1.X), Y: math.Max(v.Y, v1.Y)}
}

func (v Vec2) Unit() Vec2 {
	return v.Scale(1.0 / v.Length())
}

func (v Vec2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v1 Vec2) Dot(v2 Vec2) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}
