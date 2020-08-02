package wireguard_setup

import (
	"context"
	"io"
	"os"
	"os/exec"
)

/*
[Interface]
Address = 192.168.0.1/24
SaveConfig = true
PostUp = iptables -A FORWARD -i wg1 -j ACCEPT; iptables -t nat -A POSTROUTING -o docker0 -j MASQUERADE
PostDown = iptables -D FORWARD -i wg1 -j ACCEPT; iptables -t nat -D POSTROUTING -o docker0 -j MASQUERADE
ListenPort = 3000  # VPN Port
PrivateKey = kM687Tc/B6hJ2AYSRIfKd1a3sSuOASDT3T5X8mQ5Nk4=

[Peer]
PublicKey = jDLpHy7kazliPIyFEYP+wJBKbqW/nXQLzIzfWjwImh8=
AllowedIPs = 192.168.0.2/32
Endpoint = 37.130.122.223:38683
*/

// add gRPC connection
// tests
// parse configuration

const (
	// wireguard should be installed before hand
	wgManageBin = "wg"
	wgQuickBin  = "wg-quick"
)

type Interface struct {
	address    string // subnet
	saveConfig string
	postUp     string
	postDown   string
	listenPort uint32
	privateKey string
}

type Peer struct {
	publicKey  string
	allowedIPs string
	endPoint   string
}

// VPN server configuration file should be located under /etc/wireguard/....
// in default cases
//func initializeNIC(nicFileName string, p Interface)error {
//	if err := writeToFile(nicFileName,"[Interface]\n "+p.address+"\n"+p.saveConfig+"\n"+p.postUp+"\n"+p.postDown+"\n"+string(p.listenPort)+"\n"+p.privateKey ); err !=nil {
//		return fmt.Errorf("Initializing interface error %v", err )
//	}
//	WireGuardCmd(context.Background(),wgQuickBin, "up", nicFileName)
//
//
//	return nil
//}

// addPeer will add peer to VPN server
// wg set <wireguard-interface-name> <peer-public-key> allowed-ips 192.168.0.2/32
// example <>
func addPeer(nic, publicKey, allowedIPs string) (string, error) {
	_, err := WireGuardCmd(context.Background(), wgManageBin, "set", nic, publicKey, "allowed-ips", allowedIPs)
	if err != nil {
		return "Failed", err
	}
	return "Peer " + publicKey + " successfully added", nil
}

// removePeer will remove peer from VPN server
// wg rm <peer-public-key> allowed-ips 192.168.0.2/32
func removePeer(peerPublicKey, ipAddress string) (string, error) {
	_, err := WireGuardCmd(context.Background(), wgManageBin, "rm", peerPublicKey, "allowed-ips", ipAddress)
	if err != nil {
		return "Error", err
	}
	return "Peer " + peerPublicKey + " deleted !", nil
}

// wg show <name-of-interface>
func nicInfo(nicName string) ([]byte, error) {
	out, err := WireGuardCmd(context.Background(), wgManageBin, "show", nicName)
	if err != nil {
		return []byte("Error: "), err
	}
	return out, nil
}

// all in once
// wg genkey | tee privatekey | wg pubkey > publickey

// wg pubkey < privatekey > publickey
func generatePublicKey(privateKeyName, publicKeyName string) (string, error) {
	out, err := WireGuardCmd(context.Background(), wgManageBin, "pubkey", "<", privateKeyName, ">", publicKeyName)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

//wg genkey > privatekey
func generatePrivateKey(privateKeyName string) (string, error) {
	out, err := WireGuardCmd(context.Background(), wgManageBin, "genkey", ">", privateKeyName)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

// executes given  command from client
func WireGuardCmd(ctx context.Context, cmdBin, cmd string, cmds ...string) ([]byte, error) {
	command := append([]string{cmd}, cmds...)
	c := exec.CommandContext(ctx, cmdBin, command...)
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}
