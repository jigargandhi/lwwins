package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jigargandhi/lwwins/services"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	serverAddr := flag.String("server_addr", "0.0.0.0:3333", "")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	fmt.Println(*serverAddr)
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("unable to dial because %v ", err)
	}

	defer conn.Close()
	client := services.NewWriterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = client.Update(ctx, &services.SetValue{Value: 3})
	val, _ := client.Query(ctx, &empty.Empty{})
	log.Printf("Value %d", val)
	if err != nil {
		log.Fatalf("%v has error %v", client, err)
	}
	log.Println("Done")
}
