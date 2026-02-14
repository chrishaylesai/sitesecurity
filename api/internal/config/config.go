package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

type AuthConfig struct {
	Provider     string
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type CORSConfig struct {
	Origins string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "sitesecurity"),
			Password: getEnv("DB_PASSWORD", "sitesecurity_dev"),
			Name:     getEnv("DB_NAME", "sitesecurity"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			Provider:     getEnv("AUTH_PROVIDER", "keycloak"),
			IssuerURL:    getEnv("AUTH_ISSUER_URL", "http://localhost:8180/realms/sitesecurity"),
			ClientID:     getEnv("AUTH_CLIENT_ID", "sitesecurity-api"),
			ClientSecret: getEnv("AUTH_CLIENT_SECRET", "sitesecurity-api-secret"),
			RedirectURL:  getEnv("AUTH_REDIRECT_URL", "http://localhost:3000/auth/callback"),
		},
		CORS: CORSConfig{
			Origins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
