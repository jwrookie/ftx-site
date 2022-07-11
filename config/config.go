package config

import (
	"bytes"
	_ "embed"

	"fmt"

	"github.com/spf13/viper"
)

var config ftxConfig

//go:embed config.yaml
var configBytes []byte

//go:embed rsa_private_key.pem
var privateKey []byte

//go:embed rsa_public_key.pem
var publicKey []byte

type mysqlConfig struct {
	Url             string `mapstructure:"url"`
	Prefix          string `mapstructure:"prefix"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	SlowThreshold   int    `mapstructure:"slow_threshold"`
}

type redisConfig struct {
	Url      string `mapstructure:"url"`
	UserName string `mapstructure:"user_name"`
	PassWord string `mapstructure:"pass_word"`
	DB       int    `mapstructure:"db"`
}

type jwtConfig struct {
	SignKey string `mapstructure:"sign_key"`
	Issuer  string `mapstructure:"issuer"`
}

type appConfig struct {
	Name        string `mapstructure:"name"`
	Port        int    `mapstructure:"port"`
	Model       string `mapstructure:"model"`
	RoutePrefix string `mapstructure:"route_prefix"`
}

type logConfig struct {
	Dir        string `mapstructure:"dir"`
	FileName   string `mapstructure:"file_name"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type csrfConfig struct {
	Interval uint32 `mapstructure:"interval"`
}

type ftxConfig struct {
	App   appConfig
	Log   logConfig
	Mysql mysqlConfig
	Redis redisConfig
	Jwt   jwtConfig
	Csrf  csrfConfig
}

func init() {
	var ftx ftxConfig
	conf := viper.New()
	conf.SetConfigType("yaml")

	if err := conf.ReadConfig(bytes.NewBuffer(configBytes)); err != nil {
		panic(err)
	}

	{
		logConf := conf.Sub("app")
		if err := logConf.Unmarshal(&ftx.App); err != nil {
			panic(err)
		}
	}

	{
		logConf := conf.Sub("log")
		if err := logConf.Unmarshal(&ftx.Log); err != nil {
			panic(err)
		}
	}

	{
		logConf := conf.Sub("mysql")
		if err := logConf.Unmarshal(&ftx.Mysql); err != nil {
			panic(err)
		}
	}

	{
		logConf := conf.Sub("redis")
		if err := logConf.Unmarshal(&ftx.Redis); err != nil {
			panic(err)
		}
	}

	{
		logConf := conf.Sub("jwt")
		if err := logConf.Unmarshal(&ftx.Jwt); err != nil {
			panic(err)
		}
	}

	{
		logConf := conf.Sub("csrf")
		if err := logConf.Unmarshal(&ftx.Csrf); err != nil {
			panic(err)
		}
	}

	config = ftx

	fmt.Printf("config: %v\n", config)
}

func GetConfig() ftxConfig {
	return config
}

func GetPrivateKey() []byte {
	return privateKey
}

func GetPublicKey() []byte {
	return publicKey
}
