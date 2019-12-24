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
	TargetDir     string            `toml:"target_dir"`     //生成的model保存的目录
	Driver        string            `toml:"driver"`         //数据库类型
	Source        string            `toml:"source"`         //数据库连接信息
	TagType       []string          `toml:"tag_type"`       //生成的tag类型
	Tables        []string          `toml:"tables"`         // -t，指定生成的tables
	ExcludeTables []string          `toml:"exclude_tables"` //排除tables
	TryComplete   bool              `toml:"try_complete"`   //是否跳过错误的table
	JsonOmitempty bool              `toml:"json_omitempty"` //json是否带上omitempty
	ReplaceType   map[string]string `toml:"replace_type"`   //替换的类型，如Time替换成自己实现的
}

func Configs() AppConfigs {
	return config
}
