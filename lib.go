package natslogr

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/nats-io/stan.go"
)

type natsLogger struct {
	level       int
	prefix      string
	values      []interface{}
	stanConn    stan.Conn
	natsSubject string
}

// The Options contains necessary information about nats-connection
type Options struct {
	ClusterID, ClientID, NatsURL, Subject string
}

// NewLogger returns a logr.Logger which is implemented by nats-logr
func NewLogger(opts Options) logr.Logger {
	if err := checkNatsOptions(&opts); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return &natsLogger{
		level:       0,
		prefix:      "",
		values:      nil,
		stanConn:    connectToNatsStreamingServer(opts),
		natsSubject: opts.Subject,
	}
}

func (l *natsLogger) Info(msg string, keysAndValues ...interface{}) {
	if l.Enabled() {
		lvlStr := flatten("level", l.level)
		msgStr := flatten("msg", msg)
		trimmed := trimDuplicates(l.values, keysAndValues)
		fixedStr := flatten(trimmed[0]...)
		userStr := flatten(trimmed[1]...)
		logging.printDepth(infoLog, framesToCaller(), l.stanConn, l.natsSubject, l.prefix, " ", lvlStr, " ", msgStr, " ", fixedStr, " ", userStr)
	}
}

func (l *natsLogger) Enabled() bool {
	return bool(V(Level(l.level)))
}

func (l *natsLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	msgStr := flatten("msg", msg)
	var loggableErr interface{}
	if err != nil {
		loggableErr = err.Error()
	}
	errStr := flatten("error", loggableErr)
	trimmed := trimDuplicates(l.values, keysAndValues)
	fixedStr := flatten(trimmed[0]...)
	userStr := flatten(trimmed[1]...)
	logging.printDepth(errorLog, framesToCaller(), l.stanConn, l.natsSubject, l.prefix, " ", msgStr, " ", errStr, " ", fixedStr, " ", userStr)
}

func (l *natsLogger) V(level int) logr.InfoLogger {
	logger := l.clone()
	logger.level = level
	return logger
}

func (l *natsLogger) WithName(name string) logr.Logger {
	logger := l.clone()
	if len(l.prefix) > 0 {
		logger.prefix = l.prefix + "/"
	}
	logger.prefix += name
	return logger
}

func (l *natsLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	logger := l.clone()
	logger.values = append(logger.values, keysAndValues...)
	return logger
}

var _ logr.Logger = &natsLogger{}
