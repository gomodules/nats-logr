package natslogr

import (
	"time"

	"github.com/go-logr/logr"
	stan "github.com/nats-io/stan.go"
	natslog "gomodules.xyz/nats-logr/nats-log"
)

type natsLogger struct {
	level    int
	prefix   string
	values   []interface{}
	stanConn stan.Conn
}

type natsInfo struct {
	subject, clusterID, clientID, natsUrl string
	connectWait                           time.Duration
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
	if l.Enabled() {
		lvlStr := flatten("level", l.level)
		msgStr := flatten("msg", msg)
		trimmed := trimDuplicates(l.values, keysAndValues)
		fixedStr := flatten(trimmed[0]...)
		userStr := flatten(trimmed[1]...)
		nats, _ := getNatsInfo(l.values)
		natslog.InfoDepth(framesToCaller(), l.stanConn, nats.subject, l.prefix, " ", lvlStr, " ", msgStr, " ", fixedStr, " ", userStr)
	}
}

func (l natsLogger) Enabled() bool {
	return bool(natslog.V(natslog.Level(l.level)))
}

func (l natsLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	msgStr := flatten("msg", msg)
	var loggableErr interface{}
	if err != nil {
		loggableErr = err.Error()
	}
	errStr := flatten("error", loggableErr)
	trimmed := trimDuplicates(l.values, keysAndValues)
	fixedStr := flatten(trimmed[0]...)
	userStr := flatten(trimmed[1]...)
	nats, _ := getNatsInfo(l.values)
	natslog.ErrorDepth(framesToCaller(), l.stanConn, nats.subject, l.prefix, " ", msgStr, " ", errStr, " ", fixedStr, " ", userStr)
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
	_ logr.Logger = natsLogger{}
)
