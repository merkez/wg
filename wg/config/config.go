package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	configuration *Config
)

type Config struct {
	wgInterface WgConfig
	grpcConfig  ConnConfig
}

type WgConfig struct {
	eth     string
	dir     string
	udpPort uint
}

type ConnConfig struct {
	grpcEndpoint string
	port         uint
	tls          CertConfig
	auth         string
}

type CertConfig struct {
	Enabled   bool
	Directory string
	CertFile  string
	CertKey   string
	CAFile    string
}

func initializeConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: config \n ", err)
		return err
	}
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Println("Unmarshalling fatal error config file: config \n ", err)
		return err
	}
	return nil
}
