package main

import (
	"flag"

	stan "github.com/nats-io/go-nats-streaming"
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
	flag.CommandLine.Parse([]string{})
	natslogFlags := flag.NewFlagSet("natslog", flag.ExitOnError)
	natslog.InitFlags(natslogFlags)

	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := natslogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			f2.Value.Set(value)
		}
	})

	logger := natslogr.New().WithValues(natslogr.ClusterID, "pharmer-cluster", natslogr.ClientID, "temporary", natslogr.NatsURL, stan.DefaultNatsURL, natslogr.ConnectWait, 5, natslogr.Subject, "Create-Cluster")
	logger = logger.WithName("Something")
	logger.V(2).Info("test", "key", "values")
	logger.Error(newError("it's an error"), "error msg", "key2", "value2")
	natslog.Flush()
}
