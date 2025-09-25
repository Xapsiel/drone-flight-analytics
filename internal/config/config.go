package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabaseConfig `yaml:"database"`
	HostConfig     `yaml:"host"`
	OidcConfig     `yaml:"oidc"`
}

type HostConfig struct {
	IsProduction bool   `yaml:"isProduction"`
	Port         string `yaml:"port"`   //`env:"PORT" env-default:"8080"`
	Domain       string `yaml:"domain"` //`env:"DOMAIN" env-default:"http://127.0.0.1:8080"`
}

type DatabaseConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Name           string `yaml:"name"`
	MaxConnections int32  `yaml:"maxConnections"`
	Sslmode        string `yaml:"sslmode"`
}

type OidcConfig struct {
	ClientID       string   `yaml:"client_id"`
	RedirectURI    string   `yaml:"redirect_uri"`
	Scopes         []string `yaml:"scopes"`
	KeycloakURL    string   `yaml:"keycloak_url"`
	KeycloakRealm  string   `yaml:"keycloak_realm"`
	KeycloakSecret string   `yaml:"keycloak_secret"`
}

func New(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
