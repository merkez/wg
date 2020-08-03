#!/usr/bin/env bash

# For any issues on installation script below,
# refer official installation instructions :  https://www.wireguard.com/install/


# wireguard installation requires root privileges in order to
# manage network interface cards
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

# be careful to install WireGuard on CentOS
if [ -f /etc/centos-release ]; then
  majorVersion=$(cat /etc/centos-release | tr -dc '0-9.'|cut -d \. -f1)
  if [ $majorVersion -eq 7 ]; then
    sudo yum install -y epel-release https://www.elrepo.org/elrepo-release-7.el7.elrepo.noarch.rpm
    sudo yum install -y yum-plugin-elrepo
    sudo yum install -y kmod-wireguard wireguard-tools
  elif [ $majorVersion -eq 8 ]; then
    sudo yum install -y elrepo-release epel-release
    sudo yum install -y kmod-wireguard wireguard-tools
    # Users running non-standard kernels may wish to use the DKMS package instead of the above prebuilt kmod package, using these alternative instructions:
    fi
fi

if [ -f /etc/lsb-release ]; then
   sudo apt install -y  wireguard
fi

# there is a GUI client for Mac and Windows
# to install it, please refer https://www.wireguard.com/install/
if [ "$(uname)" == "Darwin" ]; then
#    todo: check whether brew is installed or not
    sudo brew install wireguard-tools
#    alternatively it can be installed through `port`
#    sudo port install wireguard-tools
fi









