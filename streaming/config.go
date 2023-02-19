package streaming

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config is an application config
// Should be used only in main packages for config parsing and dependency initialization.
type Config struct {
	SignalAction SignalActionConfig
}

type SignalActionConfig struct {
	UpdateVariant      string
	UpdateRendition    string
	UpdateSegment      string
	UpdatePart         string
	UpdateDemuxSegment string
	UpdateDemuxPart    string
	Ping               string
	Abort              string
	Terminate          string
	AckVariant         string
	AckRendition       string
	AckSegment         string
	AckPart            string
	AckDemuxSegment    string
	AckDemuxPart       string
	Pong               string
	Aborted            string
	Terminated         string
}

var (
	conf *Config
)

func initConf() *Config {
	var configPath string = ""
	configPath = "config.yml"
	viper.SetConfigFile(configPath)
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Errorf("read config file, %w", err))
	}
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Println(fmt.Errorf("decode config, %w", err))
	}
	return conf
}

func LoadedConfig() *Config {
	return conf
}

func InitConfig() {
	conf = initConf()
}
