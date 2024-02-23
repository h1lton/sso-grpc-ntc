package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// MustLoad загружает конфиг по пути который указан
// в переменной окружения "CONFIG_PATH"
// или в флаге командной строки "--config".
//
// "Must" в имени функции обозначает
// что функция не возвращает ошибку а сразу паникует.
func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("Не указан путь к конфигу")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Файл конфига не существует: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Не удалось прочитать конфиг: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath извлекает путь конфигурации
// из флага командной строки или переменной среды.
//
// Приоритет: флаг > окружение > по умолчанию.
//
// Значение по умолчанию — пустая строка.
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
