package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

type Environment string

const (
	Development Environment = "development"
	Staging                 = "staging"
	Production              = "production"
)

const DefaultPathToConfigPath = "."

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Environment Environment
	Host        string
	Port        int

	IsDev   bool
	IsStage bool
	IsProd  bool

	Logger  LoggerConf
	Storage StorageConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

type StorageConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

// TODO: check panic messages
func NewConfig(pathToConfig string) Config {
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

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	defineEnv(&conf)
	return conf
}

func defineEnv(cfg *Config) {
	cfg.IsDev = cfg.Environment == Development
	cfg.IsStage = cfg.Environment == Staging
	cfg.IsProd = cfg.Environment == Production
}

// TODO
