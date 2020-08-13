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



