package main

import (
	"context"
	"fmt"
	"log"

	wg "github.com/mrturkmencom/wg/wg/proto"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := wg.NewWireguardClient(conn)

	_, err = client.GetPrivateKey(context.Background(), &wg.PrivKeyReq{NicPriv: "test_private_key"})
	if err != nil {
		fmt.Sprintf("Error happened in creating private key %v", err)
		//panic(err)
	}
	//fmt.Println(r.PrivateKey)
	//r, err := client.InitializeI(context.Background(), &wg.IReq{
	//	Address:    "10.100.50.1/24",
	//	ListenPort: 5280,
	//	SaveConfig: true,
	//	Eth:        "ens3",
	//	IName:      "wg35",
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Response : %s", r.Message)

}
