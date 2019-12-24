package common

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/hsyan2008/hfw/common"
)

var config AppConfigs
var configFile string

var configLoaded bool

func LoadConfig() (err error) {
	if configLoaded {
		return
	}
	configLoaded = true

	var tables string
	flag.StringVar(&configFile, "c", "config.toml", "config file")
	flag.StringVar(&tables, "t", "", "specify tables")
	flag.Parse()

	configFile = filepath.Join(common.GetAppPath(), configFile)

	if !common.IsExist(configFile) {
		return errors.New("config file not exist")
	}

	if _, err = toml.DecodeFile(configFile, &config); err != nil {
		return fmt.Errorf("decode toml file faild: %s", err)
	}

	if tables != "" {
		config.Tables = strings.Split(tables, ",")
	}

	return
}

type AppConfigs struct {
	TargetDir     string   `toml:"target_dir"`
	Driver        string   `toml:"driver"`
	Source        string   `toml:"source"`
	TagType       []string `toml:"tag_type"`
	Tables        []string `toml:"tables"` // -t
	ExcludeTables []string `toml:"exclude_tables"`
	TryComplete   bool     `toml:"try_complete"`
	JsonOmitempty bool     `toml:"json_omitempty"`
}

func Configs() AppConfigs {
	return config
}
