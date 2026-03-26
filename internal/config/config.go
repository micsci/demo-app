package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	AppPort      string `env:"APP_PORT" envDefault:"8080"`
	GRPCPort     string `env:"GRPC_PORT" envDefault:"9090"`
	DBHost       string `env:"POSTGRES_HOST" envDefault:"localhost"`
	DBPort       string `env:"POSTGRES_PORT" envDefault:"5432"`
	DBUser       string `env:"POSTGRES_USER" envDefault:"marketplace"`
	DBPassword   string `env:"POSTGRES_PASSWORD" envDefault:"marketplace"`
	DBName       string `env:"POSTGRES_DB" envDefault:"marketplace"`
	DBSSLMode    string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	MigrationsUp bool   `env:"MIGRATIONS_UP" envDefault:"true"`
}

func (c Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
