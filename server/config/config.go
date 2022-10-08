package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const cfgFilename = "cfg.yml"

// Config отвечает за конфигурацию
type Config struct {
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	TimeoutConnPsql time.Duration `yaml:"timeout_conn_psql"`
	StringConnPsql  string        `yaml:"connection_psql"`
}

var cfg Config

func init() {
	update()
}

func update() {
	b, err := os.ReadFile(cfgFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// GetConfig возвращает копию конфига
func GetConfig() Config {
	return cfg
}

// GetHostPort возвращает полный адрес
func (c Config) GetHostPort() string {
	return c.Host + ":" + c.Port
}

// GetConnectionTimeout возвращает продолжительность ожидания подключения к бд
func (c Config) GetConnectionTimeout() time.Duration {
	return c.TimeoutConnPsql
}

// GetConnectionString возвращает строку подключения к бд
func (c Config) GetConnectionString() string {
	return c.StringConnPsql
}
