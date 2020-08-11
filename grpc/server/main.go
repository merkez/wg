package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/reflection"

	proto "github.com/mrturkmencom/wg/proto"
	"github.com/mrturkmencom/wg/vpn"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	wgServer := wg.Wireguard{}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	proto.RegisterWireguardServer(grpcServer, &wgServer)
	fmt.Println("wireguard gRPC server is running ....")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
