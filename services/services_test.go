package services

import (
	"testing"
)

func TestServerImplementsWriter(t *testing.T) {
	var server Server
	var _ WriterServer = &server
}
