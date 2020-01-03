//go:generate protoc -I ../services --go_out=plugins=grpc:../services ../services/services.proto

package services

import (
	context "context"
	fmt "fmt"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jigargandhi/lwwins/address"
	"github.com/jigargandhi/lwwins/clock"
	"github.com/jigargandhi/lwwins/lastwriterwins"
	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
)

// Server holds an implementation of grpc server
type Server struct {
	payload *lastwriterwins.Payload
	address *address.Registrar
	clock   *clock.Loclock
}

// Make creates a new instance of server
func Make(clock *clock.Loclock, val int, address *address.Registrar) *Server {
	holder := &Server{}
	holder.payload = lastwriterwins.New(clock, val)
	holder.address = address
	holder.clock = clock
	go newAddressReceived(address.NewAddress, clock)
	return holder
}

func newAddressReceived(newAddress chan string, clock *clock.Loclock) {
	for address := range newAddress {
		log.Debugf("Received %s", address)
		go syncTime(address, clock)
	}
}

// Update value
func (server *Server) Update(context context.Context, value *SetValue) (*empty.Empty, error) {
	var val int
	val = (int)(value.Value)
	server.clock.Tick()
	server.payload.Update(val)
	log.Debugf("A new value %v received updating all", val)
	notifyAll(server.address, val, server.clock)
	return &empty.Empty{}, nil
}

// Query is interface implementation
func (server *Server) Query(context context.Context, empty *empty.Empty) (*MergeRequest, error) {
	val := server.payload.Value()
	return &MergeRequest{Value: (int32)(val), Timestamp: server.clock.Get()}, nil

}

// Merge dones merge
func (server *Server) Merge(context context.Context, request *MergeRequest) (*empty.Empty, error) {
	val := (int)(request.Value)
	server.clock.Tick()
	server.payload.Merge(val, request.Timestamp)
	server.clock.Update(request.Timestamp)
	return &empty.Empty{}, nil
}

// Sync updates the local time
func (server *Server) Sync(context context.Context, request *SyncRequest) (*empty.Empty, error) {
	log.Debugf("New Time received %d", request.Timestamp)
	server.clock.Tick()
	server.clock.Update(request.Timestamp)
	return &empty.Empty{}, nil
}

func (server *Server) new() WriterServer {
	return server
}

func notifyAll(address *address.Registrar, val int, clock *clock.Loclock) {
	var opt []grpc.DialOption
	opt = append(opt, grpc.WithInsecure())
	opt = append(opt, grpc.WithBlock())

	address.ForAddress(func(addr string) {
		log.Printf("sending value %v to address %v", val, addr)
		notify(opt, addr, val, clock.Get())
	})
}

func notify(opt []grpc.DialOption, addr string, val int, instance uint64) {
	addr = fmt.Sprintf("%s:3334", addr)
	conn, err := grpc.Dial(addr, opt...)
	if err != nil {
		log.Fatalf("unable to dial because %v ", err)
	}

	defer conn.Close()
	client := NewWriterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = client.Merge(ctx, &MergeRequest{Value: (int32)(val), Timestamp: instance})
	if err != nil {
		log.Warnf("Error while updating %v %v", addr, err)
	}
}

func syncTime(addr string, clock *clock.Loclock) {
	log.Debugf("Sending sync request")
	var opt []grpc.DialOption
	opt = append(opt, grpc.WithInsecure())
	opt = append(opt, grpc.WithBlock())
	addr = fmt.Sprintf("%s:3334", addr)
	conn, err := grpc.Dial(addr, opt...)
	if err != nil {
		log.Debugf("received error while syncing time %v", err)
	}
	defer conn.Close()
	client := NewWriterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = client.Sync(ctx, &SyncRequest{Timestamp: clock.Get()})
	if err != nil {
		log.Warnf("Error while syncing time to %v %v", addr, err)
	}
}
