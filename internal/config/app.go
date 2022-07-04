package config

import (
	"os"
	"time"
)

type App struct {
	ListenPort       string
	TokenTimeExpired time.Duration
	TokenSigningKey  []byte
}

func AppNew() App {
	return App{
		ListenPort:       os.Getenv("LISTEN_PORT"),
		TokenTimeExpired: 24 * time.Hour,
		TokenSigningKey:  []byte("sn-app"),
	}
}
