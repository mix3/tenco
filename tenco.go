package tenco

import (
	"fmt"
	"io"

	"github.com/kayac/go-config"
)

type Tenco struct {
	conf Config
}

func Load(confs ...string) (*Tenco, error) {
	var conf Config
	if err := config.Load(&conf, confs...); err != nil {
		return nil, fmt.Errorf("config load failed. %w", err)
	}
	return &Tenco{
		conf: conf,
	}, nil
}

func LoadWithEnv(confs ...string) (*Tenco, error) {
	var conf Config
	if err := config.LoadWithEnv(&conf, confs...); err != nil {
		return nil, fmt.Errorf("config load failed. %w", err)
	}
	return &Tenco{
		conf: conf,
	}, nil
}

func (t *Tenco) Write(w io.Writer, offset int, g Generator) error {
	return g.Generate(w, t.conf, offset)
}
