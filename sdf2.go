package spt

// http://iquilezles.org/www/articles/distfunctions2d/distfunctions2d.htm

import (
	"encoding/gob"
	"math"
)

func init() {
	gob.Register(SDFCircle{})
	gob.Register(SDFRectangle{})
	gob.Register(SDFTriangle{})
	gob.Register(SDFPolygon{})
	//	gob.Register(SDFTrapezoid{})
}

type SDF2 interface {
	SDF() func(Vec2) float64
	Circle() (Vec2, float64)
}

type SDFCircle struct {
	Radius float64
}

func (s SDFCircle) SDF() func(Vec2) float64 {
	return func(pos Vec2) float64 {
		return len2(pos) - s.Radius
	}
}

func (s SDFCircle) Circle() (Vec2, float64) {
	return Zero2, s.Radius
}

func Circle(radius float64) SDF2 {
	return SDFCircle{radius}
}

type SDFRectangle struct {
	X, Y float64
}

func (s SDFRectangle) SDF() func(Vec2) float64 {
	return func(pos Vec2) float64 {
		d := sub2(abs2(pos), V2(s.X, s.Y))
		return len2(max2(d, Zero2)) + min(max(d.X, d.Y), 0.0)
	}
}

func (s SDFRectangle) Circle() (Vec2, float64) {
	return Zero2, sqrt(s.X*s.X + s.Y*s.Y)
}

func Rectangle(x, y float64) SDF2 {
	return SDFRectangle{x / 2, y / 2}
}

type SDFTriangle struct {
	P0, P1, P2 Vec2
}

func (s SDFTriangle) SDF() func(Vec2) float64 {
	return func(pos Vec2) float64 {
		e0 := sub2(s.P1, s.P0)
		e1 := sub2(s.P2, s.P1)
		e2 := sub2(s.P0, s.P2)
		v0 := sub2(pos, s.P0)
		v1 := sub2(pos, s.P1)
		v2 := sub2(pos, s.P2)
		pq0 := sub2(v0, scale2(e0, clamp(dot2(v0, e0)/dot2(e0, e0), 0.0, 1.0)))
		pq1 := sub2(v1, scale2(e1, clamp(dot2(v1, e1)/dot2(e1, e1), 0.0, 1.0)))
		pq2 := sub2(v2, scale2(e2, clamp(dot2(v2, e2)/dot2(e2, e2), 0.0, 1.0)))
		s := sign(e0.X*e2.Y - e0.Y*e2.X)
		d := min2(min2(
			V2(dot2(pq0, pq0), s*(v0.X*e0.Y-v0.Y*e0.X)),
			V2(dot2(pq1, pq1), s*(v1.X*e1.Y-v1.Y*e1.X))),
			V2(dot2(pq2, pq2), s*(v2.X*e2.Y-v2.Y*e2.X)))
		return -sqrt(d.X) * sign(d.Y)
	}
}

func (s SDFTriangle) Circle() (Vec2, float64) {
	return Zero2, max(max(len2(s.P0), len2(s.P1)), len2(s.P2))
}

func Triangle(p0, p1, p2 Vec2) SDF2 {
	return SDFTriangle{p0, p1, p2}
}

type SDFPolygon struct {
	N int
	R float64
}

func (s SDFPolygon) SDF() func(Vec2) float64 {
	return func(p Vec2) float64 {
		pi := math.Pi
		n := float64(s.N) / 2.0
		o := pi / 2.0 / n
		a := math.Atan(p.Y / p.X)
		a = tif(p.X < 0, a+pi, a)
		t := math.Round(a/pi*n) / n * pi
		d := math.Round((a+o)/pi*n)/n*pi - o
		f := V2(math.Cos(t), math.Sin(t))

		if abs(dot2(V2(p.X, -p.Y), V2(f.Y, f.X))) < math.Sin(o)*s.R {
			return dot2(p, f) - math.Cos(o)*s.R
		}

		return len2(sub2(p, scale2(V2(math.Cos(d), math.Sin(d)), s.R)))
	}
}

func (s SDFPolygon) Circle() (Vec2, float64) {
	return Zero2, s.R * 2
}

func Polygon(n int, r float64) SDF2 {
	return SDFPolygon{n, r}
}

/*
type SDFTrapezoid struct {
	A, B   Vec2
	RA, RB float64
}

func (s SDFTrapezoid) SDF() func(Vec2) float64 {
	return func(p Vec2) float64 {
		rba := s.RB - s.RA
		baba := dot2(sub2(s.B, s.A), sub2(s.B, s.A))
		papa := dot2(sub2(p, s.A), sub2(p, s.A))
		paba := dot2(sub2(p, s.A), sub2(s.B, s.A)) / baba
		x := sqrt(papa - paba*paba*baba)
		cax := max(0.0, x-tif(paba < 0.5, s.RA, s.RB))
		cay := abs(paba-0.5) - 0.5
		k := rba*rba + baba
		f := clamp((rba*(x-s.RA)+paba*baba)/k, 0.0, 1.0)
		cbx := x - s.RA - f*rba
		cby := paba - f
		ss := tif(cbx < 0.0 && cay < 0.0, -1.0, 1.0)
		return ss * sqrt(min(cax*cax+cay*cay*baba, cbx*cbx+cby*cby*baba))
	}
}

func Trapezoid(a, b Vec2, la, lb float64) SDF2 {
	return SDFTrapezoid{a, b, la, lb}
}
*/
