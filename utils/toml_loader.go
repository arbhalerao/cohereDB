package utils

import (
	"github.com/BurntSushi/toml"
)

// LoadTomlConfig loads a TOML configuration file into the provided config structure.
func LoadTomlConfig(config interface{}, filePath string) error {
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return err
	}
	return nil
}
