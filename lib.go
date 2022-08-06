package natslogr

import (
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/nats-io/nats.go"
)

var LogNatsError = true

// prefix will be always "", since the NewFormatter() is used without any prefix

func NewAsync(nc *nats.Conn, subj string) logr.Logger {
	return NewAsyncWithOptions(nc, subj, funcr.Options{})
}

func NewAsyncWithOptions(nc *nats.Conn, subj string, opts funcr.Options) logr.Logger {
	return funcr.New(func(_, args string) {
		if err := nc.Publish(subj, []byte(args)); err != nil && LogNatsError {
			_, _ = fmt.Fprintln(os.Stderr, "failed to publish to nats", err)
		}
	}, opts)
}

func NewSync(nc *nats.Conn, subj string, timeout time.Duration) logr.Logger {
	return NewSyncWithOptions(nc, subj, timeout, funcr.Options{})
}

func NewSyncWithOptions(nc *nats.Conn, subj string, timeout time.Duration, opts funcr.Options) logr.Logger {
	return funcr.New(func(_, args string) {
		if _, err := nc.Request(subj, []byte(args), timeout); err != nil && LogNatsError {
			_, _ = fmt.Fprintln(os.Stderr, "failed to publish to nats", err)
		}
	}, opts)
}
