//go:generate protoc -I ../services --go_out=plugins=grpc:../services ../services/services.proto

package services

import (
	context "context"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jigargandhi/lwwins/address"
	"github.com/jigargandhi/lwwins/lastwriterwins"
	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

// Server holds an implementation of grpc server
type Server struct {
	payload *lastwriterwins.Payload
	address *address.Registrar
}

// Make creates a new instance of server
func Make(val int, address *address.Registrar) *Server {
	holder := &Server{}
	holder.payload = lastwriterwins.New(val)
	holder.address = address
	//server.pa
	return holder
}

// Update value
func (server *Server) Update(context context.Context, value *SetValue) (*empty.Empty, error) {
	var val int
	val = (int)(value.Value)
	server.payload.Update(val)
	log.Printf("A new value %v received updating all", val)
	notifyAll(server.address, val)
	return &empty.Empty{}, nil
}

// Query is interface implementation
func (server *Server) Query(context context.Context, empty *empty.Empty) (*SetValue, error) {
	val := server.payload.Value()
	return &SetValue{Value: (int32)(val)}, nil

}

// Merge dones merge
func (server *Server) Merge(context context.Context, request *MergeRequest) (*empty.Empty, error) {
	val := (int)(request.Value)
	server.payload.Merge(val, request.Timestamp)
	return &empty.Empty{}, nil
}

func (server *Server) new() WriterServer {
	return server
}

func notifyAll(address *address.Registrar, val int) {
	var opt []grpc.DialOption
	opt = append(opt, grpc.WithInsecure())
	opt = append(opt, grpc.WithBlock())

	address.ForAddress(func(addr string) {
		log.Printf("sending value %v to address %v", val, addr)
		notify(opt, addr, val)
	})
}

func notify(opt []grpc.DialOption, addr string, val int) {
	conn, err := grpc.Dial(addr, opt...)
	if err != nil {
		log.Fatalf("unable to dial because %v ", err)
	}

	defer conn.Close()
	client := NewWriterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = client.Merge(ctx, &MergeRequest{Value: (int32)(val), Timestamp: time.Now().UnixNano()})
	if err != nil {
		log.Warnf("Error while updating %v %v", addr, err)
	}

}
