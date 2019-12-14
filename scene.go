package spt

import (
	"bytes"
	"encoding/gob"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"runtime"
	"time"
)

func init() {
	gob.Register(color.Alpha16{})
}

type Scene struct {
	Seed      int64   // Optional
	Camera    Camera  // Required
	Stuff     []Thing // Required
	Width     int     // in pixels
	Height    int     // in pixels
	Passes    int     // number of render passes
	Samples   int     // number of jittered samples per pixel
	Bounces   int     // max shadow ray bounces
	Horizon   float64 // max scene distance from 0,0,0 to limit marching rays
	Threshold float64 // distance from SDF considered close enough to be a hit
	Ambient   Color   // color when rays stop before reaching a light
	ShadowH   float64 // shadow alpha upper limit on invisible surfaces (dark center)
	ShadowL   float64 // shadow alpha lower limit on invisible surfaces (prenumbra cut-off)
	Raster    Raster  // summed samples per pixel
}

var _ image.Image = (*Scene)(nil)
var Transparent = color.Transparent

type Pixel struct {
	Color Color
	Rays  int32   // encoding/gob won't send a slice of pixels using a platform-dependent int size
	Alpha float64 // candidate for invisible shadows-only surface
}

type Raster []Pixel

type CompressedRaster struct {
	Buf []byte
}

type Random interface {
	Float64() float64
}

func (scene *Scene) Merge(raster Raster) {
	for y := 0; y < scene.Height; y++ {
		for x := 0; x < scene.Width; x++ {
			spixel := &scene.Raster[y*scene.Width+x]
			rpixel := &raster[y*scene.Width+x]
			spixel.Color = spixel.Color.Add(rpixel.Color)
			spixel.Rays += rpixel.Rays
			spixel.Alpha += rpixel.Alpha
		}
	}
}

func (scene Scene) Render() Raster {

	if scene.Seed == 0 {
		scene.Seed = time.Now().UTC().UnixNano()
	}

	grnd := rand.New(rand.NewSource(scene.Seed))

	for i := range scene.Stuff {
		t := &scene.Stuff[i]
		t.Prepare()
	}

	raster := make(Raster, scene.Width*scene.Height)
	semaphore := make(chan struct{}, runtime.NumCPU())

	for y := 0; y < scene.Height; y++ {
		semaphore <- struct{}{}
		go func(y int) {
			rnd := rand.New(rand.NewSource(grnd.Int63()))
			for x := 0; x < scene.Width; x++ {
				for sample := 0; sample < scene.Samples; sample++ {
					u := rnd.Float64()
					v := rnd.Float64()
					r := scene.Camera.CastRay(x, y, scene.Width, scene.Height, u, v, rnd)
					c, _, p := r.PathTrace(&scene, 0, nil)
					pixel := &raster[y*scene.Width+x]
					pixel.Color = pixel.Color.Add(c)
					pixel.Alpha += p
					pixel.Rays++
				}
			}
			<-semaphore
		}(y)
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		semaphore <- struct{}{}
	}
	return raster
}

func (scene *Scene) ColorModel() color.Model {
	return color.RGBAModel
}

func (scene *Scene) At(x, y int) color.Color {

	pixel := &scene.Raster[y*scene.Width+x]
	// average
	c := pixel.Color.Scale(1.0 / float64(pixel.Rays))
	// gamma correction
	c = Color{R: math.Sqrt(c.R), G: math.Sqrt(c.G), B: math.Sqrt(c.B)}

	if c == Naught {
		return Transparent
	}

	alpha := float64(pixel.Alpha) / float64(pixel.Rays)

	if alpha < scene.ShadowL {
		alpha = 0.0
	}

	if alpha < scene.ShadowH+0.01 {
		u := func(v float64) uint8 {
			return uint8(255.0 * v)
		}
		return color.RGBA{u(c.R), u(c.G), u(c.B), u(alpha)}
	}

	return c
}

func (scene *Scene) Bounds() image.Rectangle {
	return image.Rect(0, 0, scene.Width, scene.Height)
}

func SavePNG(img image.Image, path string) {
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	ioutil.WriteFile(path, buf.Bytes(), 0644)
}
