package spt

import (
	"image/color"
	"math"
)

type Color struct {
	R, G, B float64
}

var (
	White  = Color{1.0, 1.0, 1.0}
	Black  = Color{0.001, 0.001, 0.001}
	Naught = Color{}
	Red    = Color{1.0, 0.0, 0.0}
	Blue   = Color{0.0, 0.0, 1.0}
	Green  = Color{0.0, 1.0, 0.0}
)

func (c Color) Scale(t float64) Color {
	return Color{R: c.R * t, G: c.G * t, B: c.B * t}
}

func (c Color) Mul(c2 Color) Color {
	return Color{R: c.R * c2.R, G: c.G * c2.G, B: c.B * c2.B}
}

func (c Color) Div(c2 Color) Color {
	return Color{R: c.R / c2.R, G: c.G / c2.G, B: c.B / c2.B}
}

func (c Color) Add(c2 Color) Color {
	return Color{R: c.R + c2.R, G: c.G + c2.G, B: c.B + c2.B}
}

func (c Color) Min(c2 Color) Color {
	return Color{math.Min(c.R, c2.R), math.Min(c.G, c2.G), math.Min(c.B, c2.B)}
}

func (c Color) Brightness() float64 {
	return 0.299*c.R + 0.587*c.G + 0.114*c.B
}

func (c Color) Opacity(a float64) color.RGBA {
	u := func(v float64) uint8 {
		return uint8(v * float64(255))
	}
	return color.RGBA{u(c.R), u(c.G), u(c.B), u(a)}
}

func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(math.Max(0, math.Min(float64(0xffff), c.R*float64(0xffff))))
	g := uint32(math.Max(0, math.Min(float64(0xffff), c.G*float64(0xffff))))
	b := uint32(math.Max(0, math.Min(float64(0xffff), c.B*float64(0xffff))))
	return r, g, b, 0xffff
}

func RGB(r, g, b uint8) Color {
	return Color{float64(r) * 255, float64(g) * 255, float64(b) * 255}
}

func Hex(x int) Color {
	r := float64((x>>16)&0xff) / 255
	g := float64((x>>8)&0xff) / 255
	b := float64((x>>0)&0xff) / 255
	return Color{r, g, b}
}
