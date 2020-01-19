package spt

// GLSL-like functions to make translation easier. Wait... I know what
// would help here! Generics! We should ask the Go team for those. /s

import (
	"math"
)

var (
	Zero2 = Vec2{}
	Zero3 = Vec3{}
	X2    = Vec2{X: 1}
	Y2    = Vec2{Y: 1}
	X3    = Vec3{X: 1}
	Y3    = Vec3{Y: 1}
	Z3    = Vec3{Z: 1}
)

func V2(x, y float64) Vec2 {
	return Vec2{x, y}
}

func V3(x, y, z float64) Vec3 {
	return Vec3{x, y, z}
}

func len2(v Vec2) float64 {
	return v.Length()
}

func len3(v Vec3) float64 {
	return v.Length()
}

func dot2(vA, vB Vec2) float64 {
	return vA.Dot(vB)
}

func dot3(vA, vB Vec3) float64 {
	return vA.Dot(vB)
}

func tif(f bool, a, b float64) float64 {
	if f {
		return a
	}
	return b
}

func sqrt(a float64) float64 {
	return math.Sqrt(a)
}

func abs(a float64) float64 {
	return math.Abs(a)
}

func abs2(v Vec2) Vec2 {
	return v.Abs()
}

func abs3(v Vec3) Vec3 {
	return v.Abs()
}

func pow(a, n float64) float64 {
	return math.Pow(a, n)
}

func max(a, b float64) float64 {
	return math.Max(a, b)
}

func min(a, b float64) float64 {
	return math.Min(a, b)
}

func max2(vA, vB Vec2) Vec2 {
	return vA.Max(vB)
}

func min2(vA, vB Vec2) Vec2 {
	return vA.Min(vB)
}

func max3(vA, vB Vec3) Vec3 {
	return vA.Max(vB)
}

func min3(vA, vB Vec3) Vec3 {
	return vA.Min(vB)
}

func add2(vA, vB Vec2) Vec2 {
	return vA.Add(vB)
}

func sub2(vA, vB Vec2) Vec2 {
	return vA.Sub(vB)
}

func add3(vA, vB Vec3) Vec3 {
	return vA.Add(vB)
}

func sub3(vA, vB Vec3) Vec3 {
	return vA.Sub(vB)
}

func mod2(a, b Vec2) Vec2 {
	x := a.X - b.X*math.Floor(a.X/b.X)
	y := a.Y - b.Y*math.Floor(a.Y/b.Y)
	return Vec2{x, y}
}

func mod3(a, b Vec3) Vec3 {
	x := a.X - b.X*math.Floor(a.X/b.X)
	y := a.Y - b.Y*math.Floor(a.Y/b.Y)
	z := a.Z - b.Z*math.Floor(a.Z/b.Z)
	return Vec3{x, y, z}
}

func scale2(v Vec2, f float64) Vec2 {
	return v.Scale(f)
}

func scale3(v Vec3, f float64) Vec3 {
	return v.Scale(f)
}

func mul2(v, vB Vec2) Vec2 {
	return v.Mul(vB)
}

func mul3(v, vB Vec3) Vec3 {
	return v.Mul(vB)
}

func div2(v, vB Vec2) Vec2 {
	return v.Div(vB)
}

func div3(v, vB Vec3) Vec3 {
	return v.Div(vB)
}

func clamp(val, low, high float64) float64 {
	return math.Min(high, math.Max(low, val))
}

func clamp2(v, l, h Vec2) Vec2 {
	return Vec2{X: clamp(v.X, l.X, h.X), Y: clamp(v.Y, l.Y, h.Y)}
}

func clamp3(v, l, h Vec3) Vec3 {
	return Vec3{X: clamp(v.X, l.X, h.X), Y: clamp(v.Y, l.Y, h.Y), Z: clamp(v.Z, l.Z, h.Z)}
}

func round3(v Vec3) Vec3 {
	return Vec3{X: math.Round(v.X), Y: math.Round(v.Y), Z: math.Round(v.Z)}
}

func neg2(v Vec2) Vec2 {
	return Vec2{X: -v.X, Y: -v.Y}
}

func neg3(v Vec3) Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}

func sign(val float64) float64 {
	if val < 0.0 {
		return -1.0
	}
	if val > 0.0 {
		return 1.0
	}
	return 0.0
}
