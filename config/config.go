package config

import "os"

type AppConf struct {
	DB  DB     `yaml:"db"`
	MD  string `yaml:"md"`
	GEO Geo    `yaml:"geo"`
}
type DB struct {
	Net      string `yaml:"net"`
	Driver   string `yaml:"driver"`
	Name     string `yaml:"name"`
	User     string `json:"-" yaml:"user"`
	Password string `json:"-" yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type Geo struct {
	APIKey string `yaml:"api_key"`
	GEOKey string `yaml:"geo_key"`
}

func NewAppConf() AppConf {
	return AppConf{
		DB{
			Net:      os.Getenv("DB_NET"),
			Driver:   os.Getenv("DB_DRIVER"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		os.Getenv("SECRET_KEY"),
		Geo{
			APIKey: os.Getenv("API_KEY"),
			GEOKey: os.Getenv("GEO_KEY"),
		},
	}
}
