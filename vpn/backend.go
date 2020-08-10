package wg

import (
	"context"
	"fmt"
	"os"

	pb "github.com/mrturkmencom/wg/proto"
)

type Wireguard struct {
	// add authenticator
	// maybe configuration
	// improve configuration
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
	out, err := genInterfaceConf(wgI, dir.(string))
	if err != nil {
		return &pb.IResp{Message: out}, err
	}

	out, err = upDown(ctx, r.IName, "up")
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
	out, err := upDown(ctx, r.Nic, r.Cmd)
	if err != nil {
		return &pb.ManageNICResp{Message: string(out)}, err
	}
	return &pb.ManageNICResp{Message: out}, nil
}

// wg show <interface-name>
// if interface-name is not provided by user list for all.
func (w *Wireguard) ListPeers(ctx context.Context, r *pb.ListPeersReq) (*pb.ListPeersResp, error) {
	// todo: list peers based on user request
	return &pb.ListPeersResp{}, nil
}
func (w *Wireguard) GenPrivateKey(ctx context.Context, r *pb.PrivKeyReq) (*pb.PrivKeyResp, error) {
	_, err := generatePrivateKey(ctx, r.PrivateKeyName)
	if err != nil {
		return &pb.PrivKeyResp{}, err
	}
	return &pb.PrivKeyResp{Message: "Private Key is created with name " + r.PrivateKeyName}, nil
}

func (w *Wireguard) GenPublicKey(ctx context.Context, r *pb.PubKeyReq) (*pb.PubKeyResp, error) {
	// check whether private key exists or not, if not generate one
	if _, err := os.Stat(dir.(string) + r.PrivKeyName); os.IsNotExist(err) {
		_, err := generatePrivateKey(ctx, r.PrivKeyName)
		if err != nil {
			return &pb.PubKeyResp{Message: "Error"}, fmt.Errorf("error in generation of private key %v", err)
		}
	}
	if err := generatePublicKey(ctx, r.PrivKeyName, r.PubKeyName); err != nil {
		return &pb.PubKeyResp{}, err
	}
	return &pb.PubKeyResp{Message: "Public key is generated with " + r.PubKeyName + " name"}, nil
}

func (w *Wireguard) GetPublicKey(ctx context.Context, req *pb.PubKeyReq) (*pb.PubKeyResp, error) {
	//todo: check auth here
	out, err := getContent(req.PubKeyName)
	if err != nil {
		return &pb.PubKeyResp{}, err
	}
	return &pb.PubKeyResp{Message: out}, nil
}

func (w *Wireguard) GetPrivateKey(ctx context.Context, req *pb.PrivKeyReq) (*pb.PrivKeyResp, error) {
	//todo: check auth here
	out, err := getContent(req.PrivateKeyName)
	if err != nil {
		return &pb.PrivKeyResp{}, err
	}
	return &pb.PrivKeyResp{Message: out}, nil
}
