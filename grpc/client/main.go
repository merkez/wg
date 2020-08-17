package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
	wg "github.com/mrturkmencom/wg/proto"

	"google.golang.org/grpc"
)

type Creds struct {
	Token    string
	Insecure bool
}

func (c Creds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": string(c.Token),
	}, nil
}

func (c Creds) RequireTransportSecurity() bool {
	return !c.Insecure
}

func main() {
	var conn *grpc.ClientConn
	// wg is AUTH_KEY from vpn/auth.go
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"wg": "deneme",
	})

	tokenString, err := token.SignedString([]byte("test"))
	if err != nil {
		fmt.Println("Error creating the token")
	}

	authCreds := Creds{Token: tokenString}
	dialOpts := []grpc.DialOption{}
	authCreds.Insecure = true
	dialOpts = append(dialOpts,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(authCreds))

	conn, err = grpc.Dial(":5353", dialOpts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := wg.NewWireguardClient(conn)
	//privatekey generation example
	privKeyResp, err := client.GenPrivateKey(context.Background(), &wg.PrivKeyReq{PrivateKeyName: "random_privatekey"})
	if err != nil {
		fmt.Println(fmt.Sprintf("Error happened in creating private key %v", err))
		panic(err)
	}
	fmt.Println(privKeyResp.Message)

	// publickey generation example
	publicKeyResp, err := client.GenPublicKey(context.Background(), &wg.PubKeyReq{PrivKeyName: "random_privatekey", PubKeyName: "random_publickey"})
	if err != nil {
		fmt.Println(fmt.Sprintf("Error happened in creating public key %s", err.Error()))
		//panic(err)
	}
	if publicKeyResp != nil {
		fmt.Println(publicKeyResp.Message)
	}

	privateKey, err := client.GetPrivateKey(context.Background(), &wg.PrivKeyReq{PrivateKeyName: "random_privatekey"})
	if err != nil {
		fmt.Println(fmt.Sprintf("Get content of private key error %s", err.Error()))
		//panic(err)
	}
	if privateKey != nil {
		fmt.Println(privateKey.Message)
	}

	publicKey, err := client.GetPublicKey(context.Background(), &wg.PubKeyReq{PubKeyName: "random_publickey"})
	if err != nil {
		fmt.Println(fmt.Sprintf("Get content of public key error %s", err.Error()))
		//panic(err)
	}
	if publicKey != nil {
		fmt.Println(publicKey.Message)
	}

	// insert content of privatekey in to initInterface
	interfaceGenResp, err := client.InitializeI(context.Background(), &wg.IReq{
		Address:    "10.0.2.1/24",
		ListenPort: 4000,
		SaveConfig: true,
		PrivateKey: privateKey.Message,
		Eth:        "eth0",
		IName:      "wg1",
	})
	if err != nil {
		fmt.Println(fmt.Sprintf(" Initializing interface error %v", err.Error()))

	}
	if interfaceGenResp != nil {
		fmt.Println(interfaceGenResp.Message)
	}

	nicInfoResp, err := client.GetNICInfo(context.Background(), &wg.NICInfoReq{Interface: "wg1"})
	if err != nil {
		fmt.Println(fmt.Sprintf("Getting information of interface error %s", err.Error()))
		panic(err)
	}
	if nicInfoResp != nil {
		fmt.Println(nicInfoResp.Message)
	}

	downI, err := client.ManageNIC(context.Background(), &wg.ManageNICReq{Cmd: "down", Nic: "wg1"})
	if err != nil {
		fmt.Println(fmt.Sprintf("down interface is failed %s", err.Error()))
		panic(err)
	}
	fmt.Println(downI.Message)

	//resp, err := client.GenPublicKey(context.Background(), &wg.PrivKeyReq{})
}
