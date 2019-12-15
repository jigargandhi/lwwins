package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jigargandhi/lwwins/address"
	log "github.com/sirupsen/logrus"
)

func main() {
	nodeid := flag.Int("id", 1, "nodeid")
	token := flag.String("token", "", "")
	flag.Parse()

	register := address.Make(*nodeid, *token)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	register.Start()
	log.Info("Starting lwwins service")
	<-stop
}
