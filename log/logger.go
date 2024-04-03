package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
	Writer() io.Writer
}

func New(cfg *Config) (Logger, error) {

	log, err := newLogrus(cfg.Format, cfg.Level) // TODO, enable configurable logger type
	if err != nil {
		return nil, err
	}
	return log, nil
}
