package spt

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"image"
	"log"
	"net/rpc"
	"sync"
	"time"
)

type Renderer interface {
	Render(Scene) (Raster, error)
}

func RenderSave(out string, scene Scene, renderers []Renderer) {
	for pass := range Render(scene, renderers) {
		SavePNG(pass, out)
	}
}

func Render(scene Scene, renderers []Renderer) chan image.Image {

	if len(renderers) == 0 {
		renderers = []Renderer{NewLocalRenderer()}
	}

	rasters := make(chan Raster, len(renderers))
	jobs := make(chan struct{}, len(renderers)*2)
	ready := make(chan Renderer, len(renderers))

	for _, r := range renderers {
		ready <- r
	}

	var group sync.WaitGroup

	group.Add(1)
	go func(scene Scene) {
		for range jobs {
			renderer := <-ready
			go func() {
				raster, err := renderer.Render(scene)
				if err == nil {
					ready <- renderer
					rasters <- raster
				} else {
					log.Println("renderer", renderer, err)
					jobs <- struct{}{}
					time.Sleep(5 * time.Second)
					ready <- renderer
				}
			}()
		}
		group.Done()
	}(scene)

	scene.Raster = make(Raster, scene.Width*scene.Height)

	group.Add(1)
	go func() {
		for pass := 1; pass <= scene.Passes || scene.Passes == 0; pass++ {
			jobs <- struct{}{}
		}
		group.Done()
	}()

	frames := make(chan image.Image, 1)

	go func() {
		for pass := 1; pass <= scene.Passes || scene.Passes == 0; pass++ {
			log.Println("pass", pass, "of", scene.Passes)
			scene.Merge(<-rasters)
			copy := scene
			frames <- &copy
		}
		close(jobs)
		group.Wait()
		close(frames)
	}()

	return frames
}

type LocalRenderer struct{}

func (r LocalRenderer) Render(scene Scene) (Raster, error) {
	raster := scene.Render()
	return raster, nil
}

func NewLocalRenderer() Renderer {
	return LocalRenderer{}
}

type RPCRenderer struct {
	addr string
}

func NewRPCRenderer(address string) Renderer {
	return RPCRenderer{address}
}

func (r RPCRenderer) Render(scene Scene) (Raster, error) {

	var (
		slave  *rpc.Client
		cr     CompressedRaster
		raster Raster
		err    error
	)

	if slave, err = rpc.Dial("tcp", r.addr); err == nil {
		defer slave.Close()
		err = slave.Call("RenderRPC.Render", scene, &cr)

		if err == nil {
			raster = make(Raster, scene.Width*scene.Height)
			fr := flate.NewReader(bytes.NewReader(cr.Buf))
			binary.Read(fr, binary.LittleEndian, raster)
		}
	}

	return raster, err
}
