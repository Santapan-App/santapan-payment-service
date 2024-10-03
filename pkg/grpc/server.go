package grpc

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	addr string
}

func NewGrpcServer(addr string) *grpcServer {
	return &grpcServer{
		addr: addr,
	}
}

func (g *grpcServer) Server() (net.Listener, *grpc.Server) {
	// Create a listener on TCP port 9999
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", g.addr))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", g.addr, err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	if os.Getenv("DEBUG") == "True" {
		reflection.Register(grpcServer)
	}

	// Start serving
	go func() {
		log.Printf("grpc Server Starting On Port :%s", g.addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	return lis, grpcServer
}
