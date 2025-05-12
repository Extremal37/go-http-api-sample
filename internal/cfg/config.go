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
	appLogging       = "debug"
	appListenPort    = ":8080"
	postgresHost     = "localhost"
	postgresPort     = 5432
	postgresDatabase = "postgres"
)

type App struct {
	Logging string `yaml:"logging"`
	Address string `yaml:"address"`
	Storage string `yaml:"storage"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Configuration struct {
	App      App      `yaml:"app"`
	Postgres Postgres `yaml:"postgres"`
}

// LoadAndStoreConfig - Load configuration from file.
func LoadAndStoreConfig() (*Configuration, error) {
	v := viper.New()

	v.AddConfigPath(configPath)
	v.SetConfigName(configFile)
	v.SetConfigType(configType)

	v.SetDefault("App.Logging", appLogging)
	v.SetDefault("App.ListenPort", appListenPort)
	v.SetDefault("Postgres.host", postgresHost)
	v.SetDefault("Postgres.Port", postgresPort)
	v.SetDefault("Postgres.Database", postgresDatabase)

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
