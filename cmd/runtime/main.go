package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jigargandhi/lwwins/address"
	"github.com/jigargandhi/lwwins/clock"
	"github.com/jigargandhi/lwwins/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func getHostnameOrDefault() string {
	name, err := os.Hostname()
	if err == nil {
		return name
	}

	return "node"
}

func main() {
	os.Hostname()
	nodeid := flag.String("node", getHostnameOrDefault(), "node name defaulted to hostname")
	token := flag.String("key", "", "")
	flag.Parse()

	if *token == "" {
		log.Fatal("Key is required")
	}

	clock := &clock.Loclock{}

	registrar := address.NewRegistrar(*nodeid, *token)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", 3334))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	serverImpl := services.NewServer(clock, 0, registrar)
	services.RegisterWriterServer(grpcServer, serverImpl)
	// determine whether to use TLS
	grpcServer.Serve(grpcListener)
	log.Info("Starting lwwins service")
	<-stop
}
