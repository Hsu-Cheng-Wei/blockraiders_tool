package config

import (
	"blockraiders_tool/templates"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var dir = ".blockraiders"
var file = "config"

func GetAndCreateConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, dir)
	if _, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0755); err != nil {
			return "", err
		}
	}

	return path, nil
}

func getConfigPath() (string, error) {
	dir, err := GetAndCreateConfigPath()
	path := filepath.Join(dir, file)

	if err != nil {
		return "", err
	}
	return path, nil
}

//goland:noinspection ALL
func CreateConfig() (*templates.ConfigDto, error) {
	configPath, err := getConfigPath()
	if _, err = os.Stat(configPath); err == nil {
		if err = os.Remove(configPath); err != nil {
			return nil, err
		}
	}

	var cfg templates.ConfigDto

	if err = cfg.InitChart(); err != nil {
		return nil, err
	}

	if err = cfg.InitCodegen(); err != nil {
		return nil, err
	}

	f, err := os.Create(configPath)

	defer f.Close()

	if err != nil {
		return nil, err
	}

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	if err = os.WriteFile(configPath, b, 0755); err != nil {
		return nil, err
	}

	return &cfg, err
}

func ReadConfig() (*templates.ConfigDto, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var res templates.ConfigDto
	if err = yaml.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func EnsureConfig() error {
	if _, err := GetAndCreateConfigPath(); err != nil {
		return err
	}
	if _, err := CreateConfig(); err != nil {
		return err
	}

	return nil
}
