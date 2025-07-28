package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Path string
}

const (
	IonPath     = "ion"
	StoragePath = "storage"
)

type IConfig interface {
	GetOrCreateStorage() IConfig
	SetPath(path string)
	GetPath() string
}

func (c *Config) GetOrCreateStorage() IConfig {
	ionDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	storagePath := filepath.Join(ionDir, IonPath)
	storagePath = filepath.Join(storagePath, StoragePath)
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.MkdirAll(storagePath, 0755)
	}
	c.SetPath(storagePath)
	return c
}

func (c *Config) SetPath(path string) {
	c.Path = path
}

func (c *Config) GetPath() string {
	return c.Path
}

func Init() IConfig {
	config := Config{}
	config.GetOrCreateStorage()
	return &config
}
