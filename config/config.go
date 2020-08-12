package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	configuration *Config
)

type Config struct {
	WgInterface WgConfig
	GrpcConfig  ConnConfig
}

type WgConfig struct {
	Eth     string
	Dir     string
	UdpPort uint
}

type ConnConfig struct {
	Domain struct {
		Endpoint string
		Port     uint
	}
	Tls  CertConfig
	Auth struct {
		AKey string
		SKey string
	}
}

type CertConfig struct {
	Enabled   bool
	Directory string
	CertFile  string
	CertKey   string
	CAFile    string
}

func InitializeConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: config \n ", err)
		return nil, err
	}
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Println("Unmarshalling fatal error config file: config \n ", err)
		return nil, err
	}
	return configuration, nil
}
