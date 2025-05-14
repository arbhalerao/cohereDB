package utils

import (
	"github.com/BurntSushi/toml"
)

func LoadTomlConfig(config interface{}, filePath string) error {
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return err
	}
	return nil
}
