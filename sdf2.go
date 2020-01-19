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
	gob.Register(SDFStadium{})
	gob.Register(SDFParabola{})
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

type SDFStadium struct {
	H, R1, R2 float64
}

func (s SDFStadium) SDF() func(Vec2) float64 {
	return func(p Vec2) float64 {
		p.X = abs(p.X)
		b := (s.R1 - s.R2) / s.H
		a := sqrt(1.0 - b*b)
		k := dot2(p, V2(-b, a))
		if k < 0.0 {
			return len2(p) - s.R1
		}
		if k > a*s.H {
			return len2(sub2(p, V2(0.0, s.H))) - s.R2
		}
		return dot2(p, V2(a, b)) - s.R1
	}
}

func (s SDFStadium) Circle() (Vec2, float64) {
	return Zero2, s.H + s.R1 + s.R2
}

func Stadium(h, r1, r2 float64) SDF2 {
	return SDFStadium{h, r1, r2}
}

type SDFParabola struct {
	M, H float64
}

func (s SDFParabola) SDF() func(Vec2) float64 {
	return func(pos Vec2) float64 {

		pos.X = abs(pos.X)
		m := s.M

		// capped at height
		if pos.Y > s.H {
			l := sqrt(s.H / m)
			a := Vec2{-l, s.H}
			b := Vec2{l, s.H}
			pa := sub2(pos, a)
			ba := sub2(b, a)
			h := clamp(dot2(pa, ba)/dot2(ba, ba), 0.0, 1.0)
			return len2(sub2(pa, scale2(ba, h)))
		}

		p := (2.0*m*pos.Y - 1.0) / (6.0 * m * m)
		q := abs(pos.X) / (4.0 * m * m)
		h := q*q - p*p*p
		r := sqrt(abs(h))
		var x float64
		if h > 0 {
			x = pow(q+r, 1.0/3.0) - pow(abs(q-r), 1.0/3.0)*sign(r-q)
		} else {
			x = 2.0 * math.Cos(math.Atan2(r, q)/3.0) * sqrt(p)
		}
		y := m * x * x
		return len2(sub2(pos, Vec2{x, y})) * sign(pos.X-x)
	}
}

func (s SDFParabola) Circle() (Vec2, float64) {
	x := sqrt(s.H / s.M)
	r := sqrt(x*x + s.H*s.H)
	return Zero2, r
}

// width on x-axis at height on y-axis
func Parabola(w, h float64) SDF2 {
	w = w / 2
	return SDFParabola{h / (w * w), h}
}
