package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mrturkmencom/wg/wg"
	proto "github.com/mrturkmencom/wg/wg/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	wgServer := wg.Wireguard{}

	grpcServer := grpc.NewServer()

	proto.RegisterWireguardServer(grpcServer, &wgServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	fmt.Println("wireguard gRPC server is running ....")
}
