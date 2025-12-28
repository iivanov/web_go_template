package database

type Config struct {
	Host     string `default:"localhost"`
	User     string `default:"gonewproject"`
	Password string `default:"gonewproject"`
	Name     string `default:"gonewproject"`
	Port     string `default:"5433"`
	SSLMode  string `default:"disable"`
}
