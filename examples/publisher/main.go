package main

import (
	"flag"
	"fmt"

	stan "github.com/nats-io/stan.go"
	natslogr "gomodules.xyz/nats-logr"
	natslog "gomodules.xyz/nats-logr/nats-log"
)

type errror struct {
	msg string
}

func newError(msg string) errror {
	return errror{msg}
}

func (e errror) Error() string {
	return e.msg
}

func main() {
	natslog.InitFlags(nil)
	defer natslog.Flush()

	flag.Parse()

	flag.VisitAll(func(f1 *flag.Flag) {
		fmt.Println(f1.Name, f1.Value)
	})

	logger := natslogr.New().WithValues(natslogr.ClusterID, "pharmer-cluster", natslogr.ClientID, "publisher", natslogr.NatsURL, stan.DefaultNatsURL, natslogr.ConnectWait, 5, natslogr.Subject, "nats-log-example")
	logger = logger.WithName("Something")
	logger.V(2).Info("test", "key", "values", "v2", "info")
	logger.Error(newError("it's an error"), "error msg", "key2", "value2")
}
