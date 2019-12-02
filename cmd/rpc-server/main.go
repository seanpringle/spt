package main

import (
	"flag"
	"github.com/seanpringle/spt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	port := flag.Int("p", 34242, "TCP port")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sem := make(chan struct{}, 1)
	stop := make(chan struct{})

	sem <- struct{}{}
	go func() {
		spt.RenderServeRPC(stop, *port)
		<-sem
	}()

	<-sigs
	close(stop)
	sem <- struct{}{}
}
