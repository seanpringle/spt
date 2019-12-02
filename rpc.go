package spt

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

func RenderServeRPC(stop chan struct{}, port int) error {
	if err := rpc.Register(RenderRPC{}); err != nil {
		log.Fatal(err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	log.Println("rpc-server ready")

	run := true
	conns := make(chan net.Conn, 1)
	var group sync.WaitGroup

	group.Add(1)
	go func() {
		for run {
			conn, err := listen.Accept()
			if err != nil {
				log.Println("rpc-server", err)
				continue
			}
			conns <- conn
		}
		group.Done()
	}()

	for run {
		select {
		case conn := <-conns:
			group.Add(1)
			go func() {
				rpc.ServeConn(conn)
				group.Done()
			}()
		case <-stop:
			log.Println("rpc-server stopping...")
			run = false
			if slave, err := rpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port)); err == nil {
				slave.Close()
			}
		}
	}

	group.Wait()
	log.Println("rpc-server stopped")
	return nil
}

type RenderRPC struct{}

func (srv RenderRPC) Render(in Scene, out *Raster) error {
	start := time.Now()
	*out = in.Render()
	log.Println("rpc-server frame", time.Since(start))
	return nil
}
