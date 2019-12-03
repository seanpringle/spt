package spt

type sphere struct {
	c Vec3
	r float64
}

func (bs sphere) distance(p Vec3) float64 {
	return p.Sub(bs.c).Length() - bs.r
}
