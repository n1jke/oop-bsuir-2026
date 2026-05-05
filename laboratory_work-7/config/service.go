package config

import "time"

type ServiceConfig struct {
	Key       string        `env:"SECRET_KEY,required,notEmpty"`
	Addr      string        `env:"ADDRESS,required,notEmpty"`
	Shutdown  time.Duration `env:"TIMEOUTS_SHUTDOWN" envDefault:"10s"`
	HTTPRead  time.Duration `env:"TIMEOUTS_HTTP_READ" envDefault:"5s"`
	HTTPWrite time.Duration `env:"TIMEOUTS_HTTP_WRITE" envDefault:"10s"`
	HTTPIdle  time.Duration `env:"TIMEOUTS_HTTP_IDLE" envDefault:"60s"`
}
