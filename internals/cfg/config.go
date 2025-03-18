package cfg

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	configPath = "./config/"
	configFile = "config.yaml"
	configType = "yaml"
)
const (
	appLogging    = "debug"
	appListenPort = 8080
)

type App struct {
	Logging    string `yaml:"logging"`
	ListenPort uint16 `yaml:"listen_port"`
}

type Configuration struct {
	App App `yaml:"app"`
}

// LoadAndStoreConfig - Load configuration from file.
func LoadAndStoreConfig() (*Configuration, error) {
	v := viper.New()

	v.AddConfigPath(configPath)
	v.SetConfigName(configFile)
	v.SetConfigType(configType)

	v.SetDefault("App.Logging", appLogging)
	v.SetDefault("App.ListenPort", appListenPort)

	cfg := &Configuration{}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := v.Unmarshal(cfg, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = configType
	}); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, nil
}
