package webserver

import "time"

type Config struct {
	Port              int           `default:"8080"`
	ReadHeaderTimeout time.Duration `default:"10s"`
}
