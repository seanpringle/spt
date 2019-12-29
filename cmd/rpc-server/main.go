package main

import (
	"flag"
	"fmt"
	"github.com/seanpringle/spt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	port := flag.Int("p", 34242, "TCP port")
	prof := flag.Int("prof", 0, "pprof port")
	flag.Parse()

	if *prof > 0 {
		go http.ListenAndServe(fmt.Sprintf(":%d", *prof), nil)
	}

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
