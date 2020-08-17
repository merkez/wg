package wg

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/mrturkmencom/wg/config"
)

// add gRPC connection
// tests
// parse configuration
// add append functionality to conf

const (
	// wireguard should be installed before hand
	wgManageBin = "wg"
	wgQuickBin  = "wg-quick"
)

var (
	// todo: fix configuration variables
	configuration, _ = config.InitializeConfig()
)

type Interface struct {
	address    string // subnet
	saveConfig bool
	listenPort uint32
	privateKey string
	eth        string
	iName      string
}

type Peer struct {
	publicKey  string
	allowedIPs string
	endPoint   string
}

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
func generatePublicKey(ctx context.Context, privateKeyName, publicKeyName string) error {
	directory := configuration.WgInterface.Dir
	cmd := "wg pubkey < " + directory + privateKeyName

	out, err := exec.CommandContext(ctx, "bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmd)
	}

	if err := writeToFile(directory+publicKeyName, string(out)); err != nil {
		return err
	}
	return nil
}

// wg-quick up wg0
// wg0 configuration file should be exists at /etc/wireguard/
// or the place where docker is mounted
func upDown(ctx context.Context, nic, cmd string) (string, error) {
	_, err := WireGuardCmd(ctx, wgQuickBin, cmd, nic)
	if err != nil {
		return "Error: ", err
	}
	return "Interface " + nic + " is " + cmd, nil
}

//wg genkey > privatekey
func generatePrivateKey(ctx context.Context, privateKeyName string) (string, error) {
	out, err := WireGuardCmd(ctx, wgManageBin, "genkey")
	if err != nil {
		return "Error on running wg bin, unable to generate private key", fmt.Errorf("GeneratePrivateKey error %v", err)
	}

	if err := writeToFile(privateKeyName, string(out)); err != nil {
		return "WriteToFile Error ", err
	}
	return string(out), nil
}

// getContent returns content of privateKey or publicKey depending on keyName
func getContent(keyName string) (string, error) {
	out, err := ioutil.ReadFile(configuration.WgInterface.Dir + keyName)
	if err != nil {
		return "", fmt.Errorf("could not read the file %s err: %v", keyName, err)
	}
	return string(out), nil
}

// will generate configuration file regarding to wireguard interface
func genInterfaceConf(i Interface, confPath string) (string, error) {
	wgConf := fmt.Sprintf(
		`[Interface]
Address = %s
ListenPort = %d
SaveConfig = %v
PrivateKey = %s
PostUp = iptables -A FORWARD -i %s -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE
PostDown = iptables -D FORWARD -i %s -j ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE`, i.address, i.listenPort, i.saveConfig, i.privateKey,
		i.iName, i.eth, i.iName, i.eth)

	if err := writeToFile(confPath+i.iName+".conf", wgConf); err != nil {
		return "GenInterface Error:  ", err
	}
	return i.iName + " configuration saved to " + configuration.WgInterface.Dir, nil
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
