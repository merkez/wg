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

	_, err = client.GenPrivateKey(context.Background(), &wg.PrivKeyReq{PrivateKeyName: "test_private_key"})
	if err != nil {
		fmt.Println(fmt.Sprintf("Error happened in creating private key %v", err))
	}
}
