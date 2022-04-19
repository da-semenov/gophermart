package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"
	"os"
)

type AppConfig struct {
	ServerAddress        string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabaseDSN          string `env:"DATABASE_URI" envDefault:"postgresql://practicum:practicum@127.0.0.1:5432/gophermart"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:3000"`
	ReInit               bool   `env:"REINIT" envDefault:"false"`
}

func (config *AppConfig) Init() error {
	fmt.Println(os.Args)
	if err := env.Parse(config); err != nil {
		fmt.Println("unable to load server settings", err)
		return err
	}

	pflag.StringVarP(&config.ServerAddress, "a", "a", config.ServerAddress, "Http-server address")
	pflag.StringVarP(&config.DatabaseDSN, "d", "d", config.DatabaseDSN, "Database connection string")
	pflag.StringVarP(&config.AccrualSystemAddress, "r", "r", config.AccrualSystemAddress, "Accrual system address")
	pflag.BoolVarP(&config.ReInit, "c", "c", config.ReInit, "Re-init database")
	pflag.Parse()

	return nil
}

func NewConfig() *AppConfig {
	return &AppConfig{}
}
