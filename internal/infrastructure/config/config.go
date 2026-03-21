package config

import (
	"time"
)

type Config struct {
	App      AppConfig      `koanf:"app"`
	HTTP     HTTPConfig     `koanf:"http"`
	Postgres PostgresConfig `koanf:"postgres"`
	Log      LogConfig      `koanf:"log"`
}
type AppConfig struct {
	Name string `koanf:"name"`
	Env  string `koanf:"env"`
}
type HTTPConfig struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}
type PostgresConfig struct {
	Host            string        `koanf:"host"`
	Port            int           `koanf:"port"`
	User            string        `koanf:"user"`
	Password        string        `koanf:"password"`
	Database        string        `koanf:"database"`
	SSLMode         string        `koanf:"sslmode"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
}

type LogConfig struct {
	Level string `koanf:"level"`
}
