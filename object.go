package spt

type Thing struct {
	Mat Material
	SDF3
	sdf    func(Vec3) float64
	center Vec3
	radius float64
}

func Object(mat Material, sdf SDF3) Thing {
	return Thing{mat, sdf, nil, Zero3, 0.0}
}

func (o *Thing) Material() Material {
	return o.Mat
}

func (o *Thing) Normal(pos Vec3) Vec3 {
	return SDF3Normal(o.SDF3, pos)
}

func (o *Thing) Distance(pos Vec3) float64 {
	return o.sdf(pos)
}

func (o *Thing) Prepare() {
	o.sdf = o.SDF3.SDF()
	o.center, o.radius = o.SDF3.Sphere()
}

func (o *Thing) SDF() func(Vec3) float64 {
	return o.sdf
}

func (o *Thing) Sphere() (Vec3, float64) {
	return o.center, o.radius
}

func (o *Thing) BoundingDistance(pos Vec3) float64 {
	return pos.Sub(o.center).Length() - o.radius
}
