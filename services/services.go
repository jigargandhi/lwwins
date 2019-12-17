//go:generate protoc -I ../services --go_out=plugins=grpc:../services ../services/services.proto

package services

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
)

// Server holds an implementation of grpc server
type Server struct {
}

// Update value
func (server *Server) Update(context context.Context, value *SetValue) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// Query is interface implementation
func (server *Server) Query(context context.Context, empty *empty.Empty) (*SetValue, error) {

	return &SetValue{Value: 1}, nil

}

// Merge dones merge
func (server *Server) Merge(context context.Context, request *MergeRequest) (*empty.Empty, error) {
	return nil, nil
}

func (server *Server) new() WriterServer {
	return server
}
