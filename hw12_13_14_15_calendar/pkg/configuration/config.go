package configuration

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Environment string

const (
	Development Environment = "development"
	Staging                 = "staging"
	Production              = "production"
)

const DefaultPathToConfigPath = "."

type Configuration struct {
	Environment Environment
	Host        string
	Port        int

	IsDev   bool
	IsStage bool
	IsProd  bool

	Logger  LoggerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	InMemory bool
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Dialect  string
}

func NewConfiguration(pathToConfig string) Configuration {
	if len(pathToConfig) == 0 {
		pathToConfig = DefaultPathToConfigPath
	}

	dir, file := filepath.Split(pathToConfig)
	cfgExt := filepath.Ext(file)[1:]
	fileName := file[:len(file)-(len(cfgExt)+1)]
	if len(fileName) == 0 && len(cfgExt) == 0 {
		panic("invalid configuration path")
	}

	absoluteDirPath, err := filepath.Abs(dir)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.AddConfigPath(absoluteDirPath)
	viper.SetConfigName(fileName)
	viper.SetConfigType(cfgExt)

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var conf Configuration
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	defineEnv(&conf)
	return conf
}

func defineEnv(cfg *Configuration) {
	cfg.IsDev = cfg.Environment == Development
	cfg.IsStage = cfg.Environment == Staging
	cfg.IsProd = cfg.Environment == Production
}
