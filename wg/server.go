package wg

import (
	"context"
	"fmt"

	pb "github.com/mrturkmencom/wireguard-setup/wg/proto"
)

type Wireguard struct {
	// add authenticator
	// maybe configuration
}

func (w *Wireguard) InitializeI(ctx context.Context, r *pb.IReq) (*pb.IResp, error) {
	s, err := generatePrivateKey(ctx, r.IName+"_private_key")
	if err != nil {
		return &pb.IResp{}, err
	}
	fmt.Printf("Auto-generated private key from wg : %s\n", s)

	wgI := Interface{
		address:    r.Address,
		listenPort: r.ListenPort,
		privateKey: s,
		eth:        r.Eth,
		saveConfig: r.SaveConfig,
		iName:      r.IName,
	}
	out, err := genInterfaceConf(wgI, "/etc/wireguard")
	if err != nil {
		return &pb.IResp{Message: out}, err
	}

	out, err = upDown(r.IName, "up")
	if err != nil {
		return &pb.IResp{Message: out}, err
	}

	return &pb.IResp{Message: out}, nil
}

func (w *Wireguard) AddPeer(ctx context.Context, r *pb.AddPReq) (*pb.AddPResp, error) {

	out, err := addPeer(r.Nic, r.PublicKey, r.AllowedIPs)
	if err != nil {
		return &pb.AddPResp{Message: out}, err
	}
	return &pb.AddPResp{Message: out}, nil
}

func (w *Wireguard) DelPeer(ctx context.Context, r *pb.DelPReq) (*pb.DelPResp, error) {
	out, err := removePeer(r.PeerPublicKey, r.IpAddress)
	if err != nil {
		return &pb.DelPResp{Message: out}, err
	}

	return &pb.DelPResp{Message: out}, nil
}

func (w *Wireguard) GetNICInfo(ctx context.Context, r *pb.NICInfoReq) (*pb.NICInfoResp, error) {
	out, err := nicInfo(r.Interface)
	if err != nil {
		return &pb.NICInfoResp{Message: string(out)}, err
	}
	return &pb.NICInfoResp{Message: string(out)}, nil
}

func (w *Wireguard) ManageNIC(ctx context.Context, r *pb.ManageNICReq) (*pb.ManageNICResp, error) {
	// todo : add up and down commands
	return &pb.ManageNICResp{}, nil
}

func (w *Wireguard) GetPrivateKey(ctx context.Context, r *pb.PrivKeyReq) (*pb.PrivKeyResp, error) {
	// privatekey
	return &pb.PrivKeyResp{}, nil
}

func (w *Wireguard) GetPublicKey(ctx context.Context, r *pb.PubKeyReq) (*pb.PubKeyResp, error) {
	// todo: get public key
	return &pb.PubKeyResp{}, nil
}
