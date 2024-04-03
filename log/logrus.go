package log

import (
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

var _ Logger = &logrusLogger{}

type logrusLogger struct {
	*logrus.Logger
}

func (l logrusLogger) Writer() io.Writer {
	return l.Out
}

func newLogrus(format, level string) (*logrusLogger, error) {
	var formatter logrus.Formatter

	switch format {
	case "text":
		formatter = &logrus.TextFormatter{}
	case "json", "":
		formatter = &logrus.JSONFormatter{}
	default:
		return nil, fmt.Errorf("invalid log encoding format %q", format)
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	l := logrus.New()
	l.SetFormatter(formatter)
	l.SetLevel(logLevel)

	return &logrusLogger{l}, nil
}
