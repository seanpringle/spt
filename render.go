package spt

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"log"
	"net/rpc"
	"sync"
	"time"
)

type Renderer interface {
	Render(Scene) (Raster, error)
}

func Render(out string, scene Scene, renderers []Renderer) {

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

	for pass := 1; pass <= scene.Passes || scene.Passes == 0; pass++ {
		log.Println("pass", pass, "of", scene.Passes)

		jobs <- struct{}{}
		scene.Merge(<-rasters)

		buf := new(bytes.Buffer)
		png.Encode(buf, &scene)
		ioutil.WriteFile(out, buf.Bytes(), 0644)
	}

	close(jobs)
	group.Wait()
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
		raster Raster
		err    error
	)

	if slave, err = rpc.Dial("tcp", r.addr); err == nil {
		defer slave.Close()
		err = slave.Call("RenderRPC.Render", scene, &raster)
	}

	return raster, err
}
