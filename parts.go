package spt

func WorkBench(size float64) Thing {
	return Object(
		Matt(Color{0.16, 0.12, 0.09}),
		//Matt(Color{0.15, 0.15, 0.15}),
		Translate(V3(0, 0, -5), Cube(size, size, 10)),
	)
}

func GearWheel() SDF3 {

	teeth := []SDF3{}
	for i := 0; i < 18; i++ {
		teeth = append(teeth,
			Rotate(
				V3(0, 0, 1), float64(i)*20.0,
				Translate(
					V3(400, 0, 0),
					//Extrude(100, Trapezoid(V2(0, 0), V2(100, 0), 50, 40)),
					Distort(V3(1, 0.4, 1), Intersection(
						Cylinder(200, 110),
						Cube(200, 200, 200),
					)),
				),
			),
		)
	}

	return Difference(
		Union(append([]SDF3{Cylinder(200, 420)}, teeth...)...),
		Cylinder(400, 200),
	)
}
