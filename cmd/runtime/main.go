package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jigargandhi/lwwins/address"
	"github.com/jigargandhi/lwwins/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	nodeid := flag.Int("id", 1, "nodeid")
	token := flag.String("token", "", "")
	flag.Parse()

	register := address.Make(*nodeid, *token)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)

	register.Start()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3334))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	serverImpl := services.Make(0, register)
	services.RegisterWriterServer(grpcServer, serverImpl)
	// determine whether to use TLS
	grpcServer.Serve(lis)
	log.Info("Starting lwwins service")
	<-stop
}
