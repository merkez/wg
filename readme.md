# wg

Wireguard backed and gRPC wrapped server which is responsible  to create VPN connection through gRPC requests. 
The idea is basically having remote control to gRPC endpoint to be able to setup a VPN connection from your client. 

As initial step, dockerization of wg is dismissed for now, however it will be added. 

## Installation of wireguard

Most of the cases [official installation page](https://www.wireguard.com/install/) is enough to install wireguard however, 
in some cases, the instructions are misleading on official page, hence I am including installation
steps for Debian.  (-in case of error in official installation following steps could be followed -) 

```bash 
$ sudo apt update
$ sudo apt upgrade
$ sudo sh -c "echo 'deb http://deb.debian.org/debian buster-backports main contrib non-free' > /etc/apt/sources.list.d/buster-backports.list"
$ sudo apt update
$ apt search wireguard
$ sudo apt install wireguard
# in some cases command line tools does not  work for wireguard in that case do following 
$ apt-get install wireguard-dkms wireguard-tools linux-headers-$(uname -r)
```

## Available gRPC calls

- GenPrivateKey 
    - Generates private key which is required to initialize wireguard interface.gRPC call requires only name of the file,
    which will have private key in it. 
    
    Example Usage: 
    ````go
    privKeyResp, err := client.GenPrivateKey(context.Background(), &wg.PrivKeyReq{PrivateKeyName: "random_privatekey"})
    if err != nil {
    	fmt.Println(fmt.Sprintf("Error happened in creating private key %v", err))
    	panic(err)
    }
    fmt.Println(privKeyResp.Message)
    ````
    Private key will be availabe in defined configuration directory in config.yml file.
    
- GenPublicKey 
    - Generates pair of private key as public key, in order to use this functionality, it requires
    existing private key name (which is generated in earlier step) then public key name (-which will be generated-)
    
    Example Usage: 
    ````go
    publicKeyResp, err := client.GenPublicKey(context.Background(), &wg.PubKeyReq{PrivKeyName: "random_privatekey", PubKeyName: "random_publickey"})
    if err != nil {
    	fmt.Println(fmt.Sprintf("Error happened in creating public key %s", err.Error()))
    	panic(err)
    }
    if publicKeyResp != nil {
    	fmt.Println(publicKeyResp.Message)
    }
    ````

- GetPrivateKey

    - Despite of GenPrivateKey functionality, this one returns existing private key content. 
      
      Example Usage: 
      ````go
      privateKey, err := client.GetPrivateKey(context.Background(), &wg.PrivKeyReq{PrivateKeyName: "random_privatekey"})
      if err != nil {
      	fmt.Println(fmt.Sprintf("Get content of private key error %s", err.Error()))
      	panic(err)       
	  }	      
      if privateKey != nil {
      	fmt.Println(privateKey.Message)      
      }
      ````
- GetPublicKey

    - Returns content of existing public key content
    ````go
    publicKey, err := client.GetPublicKey(context.Background(), &wg.PubKeyReq{PubKeyName: "random_publickey"})
    if err != nil {
    	fmt.Println(fmt.Sprintf("Get content of public key error %s", err.Error()))
    	panic(err)
    }
    if publicKey != nil {
    	fmt.Println(publicKey.Message)
    }
    ````

- InitializeI 
    -  It is for initializing wireguard interface in configuration folder which is provided in configuration file. It requires
       wireguard interface specifications, which are 
       ```raw
       Address: <subnet-of-interface>
       ListenPort: <where-users-will-be-connected-to>
       SaveConfig: <whether-save-config-or-not>
       PrivateKey: <private-key-of-server> 
       Eth : <main-ethernet-point-to-outside>
       IName: <interface-name-required-in-grpc-call> 
       ```  
       Example usage: 
       ````go
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
       fmt.Println(interfaceGenResp.Message)
       ````
       
- GetNICInfo
    - Returns information regarding to requested wireguard interface. 
   
   Example Usage:
   
   ````go
  	nicInfoResp, err := client.GetNICInfo(context.Background(), &wg.NICInfoReq{Interface: "wg1"})
  	if err != nil {
  		fmt.Println(fmt.Sprintf("Getting information of interface error %s", err.Error()))
  		panic(err)
  	}
  	if nicInfoResp != nil {
  		fmt.Println(nicInfoResp.Message)
  	}
   ````
 
- ManageNIC 
    - It can up or down given wg interface. 
   
   Example Usage: 
   ````go 
   downI, err := client.ManageNIC(context.Background(), &wg.ManageNICReq{Cmd: "down", Nic: "wg1"})
   	if err != nil {
   		fmt.Println(fmt.Sprintf("down interface is failed %s", err.Error()))
   		panic(err)
   	}
   fmt.Println(downI.Message) 
   ````
       

