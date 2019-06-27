package natslogr

import (
	"time"

	"github.com/go-logr/logr"
	stan "github.com/nats-io/go-nats-streaming"
)

type natsLogger struct {
	level    int
	prefix   string
	values   []interface{}
	stanConn stan.Conn
}

type natsInfo struct {
	clusterID, clientID, natsUrl string
	connectWait                  time.Duration
}

// New returns a logr.Logger which is implemented by nats-logr
func New() logr.Logger {
	return natsLogger{
		level:    0,
		prefix:   "",
		values:   nil,
		stanConn: nil,
	}
}

func (l natsLogger) Info(msg string, keysAndValues ...interface{}) {

}

func (l natsLogger) Enabled() bool {
	return true
}

func (l natsLogger) Error(err error, msg string, keysAndValues ...interface{}) {

}

func (l natsLogger) V(level int) logr.InfoLogger {
	logger := l.clone()
	logger.level = level
	return logger
}

func (l natsLogger) WithName(name string) logr.Logger {
	logger := l.clone()
	if len(l.prefix) > 0 {
		logger.prefix = l.prefix + "/"
	}
	logger.prefix += name
	return logger
}

func (l natsLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	logger := l.clone()
	logger.values = append(logger.values, keysAndValues...)

	nats, ok := getNatsInfo(logger.values)

	if logger.stanConn == nil && ok {
		logger.stanConn = connectToNatsStreamingServer(nats)
	}

	return logger
}

var (
	_ logr.Logger     = natsLogger{}
	_ logr.InfoLogger = natsLogger{}
)
