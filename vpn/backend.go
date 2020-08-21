package wg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/credentials"

	"github.com/mrturkmencom/wg/config"

	"google.golang.org/grpc"

	pb "github.com/mrturkmencom/wg/proto"
)

type wireguard struct {
	auth   Authenticator
	config *config.Config
}

// InitializeI creates interface configuration and make it UP.
func (w *wireguard) InitializeI(ctx context.Context, r *pb.IReq) (*pb.IResp, error) {

	s, err := generatePrivateKey(ctx, w.config.WgInterface.Dir+r.IName+"_priv")
	if err != nil {
		return &pb.IResp{}, err
	}

	wgI := Interface{
		address:    r.Address,
		listenPort: r.ListenPort,
		privateKey: s,
		eth:        r.Eth,
		saveConfig: r.SaveConfig,
		iName:      r.IName,
	}
	out, err := genInterfaceConf(wgI, w.config.WgInterface.Dir)
	if err != nil {
		return &pb.IResp{Message: out}, err
	}

	out, err = upDown(ctx, r.IName, "up")
	if err != nil {
		return &pb.IResp{Message: out}, err
	}
	log.Debug().Str("Address: ", r.Address).
		Uint32("ListenPort: ", r.ListenPort).
		Str("Ethernet I: ", r.Eth).
		Str("PrivateKey: ", r.PrivateKey).
		Bool("SaveConfig", r.SaveConfig).Msgf("Interface %s created and it is up", r.IName)

	return &pb.IResp{Message: out}, nil
}

// AddPeer adds peer to given wireguard interface
func (w *wireguard) AddPeer(ctx context.Context, r *pb.AddPReq) (*pb.AddPResp, error) {

	out, err := addPeer(r.Nic, r.PublicKey, r.AllowedIPs)
	if err != nil {
		return &pb.AddPResp{Message: out}, err
	}
	log.Info().Msgf("Peer with public key: { %s } is added to interface: { %s } from allowed-ips: { %s }", r.PublicKey, r.Nic, r.AllowedIPs)

	return &pb.AddPResp{Message: out}, nil
}

// DelPeer deletes peer from given wireguard interface
func (w *wireguard) DelPeer(ctx context.Context, r *pb.DelPReq) (*pb.DelPResp, error) {
	out, err := removePeer(r.PeerPublicKey, r.IpAddress)
	if err != nil {
		return &pb.DelPResp{Message: out}, err
	}
	log.Info().Msgf("Peer with public key: { %s } is deleted from ip-address: { %s }", r.PeerPublicKey, r.IpAddress)
	return &pb.DelPResp{Message: out}, nil
}

// GetNICInfo returns general information about given wireguard interface
func (w *wireguard) GetNICInfo(ctx context.Context, r *pb.NICInfoReq) (*pb.NICInfoResp, error) {
	out, err := nicInfo(r.Interface)
	if err != nil {
		return &pb.NICInfoResp{Message: string(out)}, err
	}
	log.Debug().Msgf("NIC Information for { %s } is printed ", r.Interface)
	return &pb.NICInfoResp{Message: string(out)}, nil
}

// ManageNIC is managing (up & down) given wireguard interface
func (w *wireguard) ManageNIC(ctx context.Context, r *pb.ManageNICReq) (*pb.ManageNICResp, error) {
	out, err := upDown(ctx, r.Nic, r.Cmd)
	if err != nil {
		return &pb.ManageNICResp{Message: string(out)}, err
	}
	log.Info().Msgf("ManageNIC: interface %s is called to be %s", r.Nic, r.Cmd)
	return &pb.ManageNICResp{Message: out}, nil
}

// wg show <interface-name>
// if interface-name is not provided by user list for all.
func (w *wireguard) ListPeers(ctx context.Context, r *pb.ListPeersReq) (*pb.ListPeersResp, error) {
	out, err := listPeers(ctx, r.Nicname)
	if err != nil {
		log.Printf("Error in listing peers in gRPC %v", err)
		return &pb.ListPeersResp{}, err
	}
	log.Info().Msgf("ListPeers: listing peers for %s interface", r.Nicname)
	return &pb.ListPeersResp{Response: out}, nil
}

// GenPrivateKey generates PrivateKey for wireguard interface
func (w *wireguard) GenPrivateKey(ctx context.Context, r *pb.PrivKeyReq) (*pb.PrivKeyResp, error) {

	_, err := generatePrivateKey(ctx, w.config.WgInterface.Dir+r.PrivateKeyName)
	if err != nil {
		return &pb.PrivKeyResp{}, err
	}
	log.Info().Msgf("GenPrivateKey is called to generate new private key with filename %s", r.PrivateKeyName)
	return &pb.PrivKeyResp{Message: "Private Key is created with name " + w.config.WgInterface.Dir + r.PrivateKeyName}, nil
}

// GenPublicKey generates PublicKey for wireguard interface
func (w *wireguard) GenPublicKey(ctx context.Context, r *pb.PubKeyReq) (*pb.PubKeyResp, error) {
	// check whether private key exists or not, if not generate one
	if _, err := os.Stat(w.config.WgInterface.Dir + r.PrivKeyName); os.IsNotExist(err) {
		fmt.Printf("PrivateKeyFile is not exists, creating one ... %s\n", r.PrivKeyName)
		_, err := generatePrivateKey(ctx, w.config.WgInterface.Dir+r.PrivKeyName)
		if err != nil {
			return &pb.PubKeyResp{Message: "Error"}, fmt.Errorf("error in generation of private key %v", err)
		}
	}

	if err := generatePublicKey(ctx, r.PrivKeyName, r.PubKeyName); err != nil {
		return &pb.PubKeyResp{}, err
	}
	return &pb.PubKeyResp{Message: "Public key is generated with " + w.config.WgInterface.Dir + r.PubKeyName + " name"}, nil
}

// GetPublicKey returns content of given PublicKey
func (w *wireguard) GetPublicKey(ctx context.Context, req *pb.PubKeyReq) (*pb.PubKeyResp, error) {
	//todo: check auth here
	out, err := getContent(req.PubKeyName)
	if err != nil {
		return &pb.PubKeyResp{}, err
	}
	return &pb.PubKeyResp{Message: out}, nil
}

// GetPrivateKey returns content of given PrivateKey
func (w *wireguard) GetPrivateKey(ctx context.Context, req *pb.PrivKeyReq) (*pb.PrivKeyResp, error) {
	//todo: check auth here
	out, err := getContent(req.PrivateKeyName)
	if err != nil {
		return &pb.PrivKeyResp{}, err
	}
	return &pb.PrivKeyResp{Message: out}, nil
}

func GetCreds(conf config.CertConfig) (credentials.TransportCredentials, error) {
	log.Printf("Preparing credentials for RPC")

	certificate, err := tls.LoadX509KeyPair(conf.CertFile, conf.CertKey)
	if err != nil {
		return nil, fmt.Errorf("could not load server key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(conf.CAFile)
	if err != nil {
		return nil, fmt.Errorf("could not read ca certificate: %s", err)
	}
	// CA file for let's encrypt is located under domain conf as `chain.pem`
	// pass chain.pem location
	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("failed to append client certs")
	}

	// Create the TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})
	return creds, nil
}

// SecureConn enables communication over secure channel
func SecureConn(conf config.CertConfig) ([]grpc.ServerOption, error) {
	if conf.Enabled {
		creds, err := GetCreds(conf)

		if err != nil {
			return []grpc.ServerOption{}, errors.New("Error on retrieving certificates: " + err.Error())
		}
		log.Printf("Server is running in secure mode !")
		return []grpc.ServerOption{grpc.Creds(creds)}, nil
	}
	return []grpc.ServerOption{}, nil
}

func InitServer(conf *config.Config) (*wireguard, error) {

	gRPCServer := &wireguard{
		auth:   NewAuthenticator(conf.GrpcConfig.Auth.SKey, conf.GrpcConfig.Auth.AKey),
		config: conf,
	}
	return gRPCServer, nil
}

// AddAuth adds authentication to gRPC server
func (w *wireguard) AddAuth(opts ...grpc.ServerOption) *grpc.Server {
	streamInterceptor := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := w.auth.AuthenticateContext(stream.Context()); err != nil {
			return err
		}
		return handler(srv, stream)
	}

	unaryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := w.auth.AuthenticateContext(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}

	opts = append([]grpc.ServerOption{
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	}, opts...)
	return grpc.NewServer(opts...)

}
