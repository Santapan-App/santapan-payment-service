package grpc

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

type grpcClient struct {
	addr    string
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewGrpcClient creates a new gRPC client with the given address and timeout.
func NewGrpcClient(addr string, timeout time.Duration) (*grpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %v", err)
	}

	return &grpcClient{
		addr:    addr,
		conn:    conn,
		timeout: timeout,
	}, nil
}

// Close closes the connection to the gRPC server.
func (c *grpcClient) Close() {
	if err := c.conn.Close(); err != nil {
		log.Fatalf("failed to close gRPC connection: %v", err)
	}
}
