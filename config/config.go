package config

type Config struct {
	BaseURL string `json:"base_url"`
	Port    string `json:"port"`
}

var AppConfig Config
