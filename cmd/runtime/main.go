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

func main() {
	nodeid := flag.Int("id", 1, "nodeid")
	token := flag.String("token", "", "")
	flag.Parse()

	clock := &clock.Loclock{}

	registrar := address.NewRegistrar(*nodeid, *token)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	grpc_listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 3334))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	serverImpl := services.NewServer(clock, 0, registrar)
	services.RegisterWriterServer(grpcServer, serverImpl)
	// determine whether to use TLS
	grpcServer.Serve(grpc_listener)
	log.Info("Starting lwwins service")
	<-stop
}
