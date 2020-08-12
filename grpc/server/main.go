package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/mrturkmencom/wg/config"

	"google.golang.org/grpc/reflection"

	proto "github.com/mrturkmencom/wg/proto"
	"github.com/mrturkmencom/wg/vpn"
)

func main() {

	configuration, err := config.InitializeConfig()
	if err != nil {
		panic("Configuration initialization error: " + err.Error())
	}
	port := strconv.FormatUint(uint64(configuration.GrpcConfig.Domain.Port), 10)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	wgServer, err := wg.InitServer(configuration)
	if err != nil {
		return
	}
	opts, err := wg.SecureConn(configuration.GrpcConfig.Tls)
	if err != nil {
		log.Fatalf("failed to retrieve secure options %s", err.Error())
	}

	gRPCEndpoint := wgServer.AddAuth(opts...)

	reflection.Register(gRPCEndpoint)
	proto.RegisterWireguardServer(gRPCEndpoint, wgServer)

	fmt.Printf("wireguard gRPC server is running at port %s...\n", port)
	if err := gRPCEndpoint.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
