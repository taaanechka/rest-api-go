package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"
	"github.com/taaanechka/rest-api-go/pkg/logging"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Users userstorage.Config `yaml:"users"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		lg := logging.GetLogger()
		lg.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			lg.Info(help)
			lg.Fatal(err)
		}
	})

	return instance
}
