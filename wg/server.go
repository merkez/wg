package wg

import (
	"context"

	pb "github.com/mrturkmencom/wireguard-setup/wg/proto"
)

type wireguard struct {
}

func (w *wireguard) InitializeI(ctx context.Context, r *pb.IReq) (*pb.IResp, error) {
	// todo: initialize interface for wireguard

	return &pb.IResp{}, nil
}

func (w *wireguard) AddPeer(ctx context.Context, r *pb.AddPReq) (*pb.AddPResp, error) {
	// todo : add peer functionality for wg
	return &pb.AddPResp{}, nil
}

func (w *wireguard) DelPeer(ctx context.Context, r *pb.DelPReq) (*pb.DelPResp, error) {
	// todo: add delete operation for wg
	return &pb.DelPResp{}, nil
}

func (w *wireguard) GetNICInfo(ctx context.Context, r *pb.NICInfoReq) (*pb.NICInfoResp, error) {
	// todo: return information about desired wg interface
	return &pb.NICInfoResp{}, nil
}

func (w *wireguard) ManageNIC(ctx context.Context, r *pb.ManageNICReq) (*pb.ManageNICResp, error) {
	// todo : add up and down commands
	return &pb.ManageNICResp{}, nil
}

func (w *wireguard) GetPrivateKey(ctx context.Context, r *pb.PrivKeyReq) (*pb.PrivKeyResp, error) {
	// todo : get private key from gRPC server
	return &pb.PrivKeyResp{}, nil
}

func (w *wireguard) GetPublicKey(ctx context.Context, r *pb.PubKeyReq) (*pb.PubKeyResp, error) {
	// todo: get public key
	return &pb.PubKeyResp{}, nil
}
