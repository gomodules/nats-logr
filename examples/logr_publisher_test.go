package natslogr_test

import (
	"flag"
	"fmt"

	natslogr "gomodules.xyz/nats-logr"

	"github.com/nats-io/stan.go"
)

type errorr struct {
	msg string
}

func newErrorr(msg string) errorr {
	return errorr{msg: msg}
}

func (e errorr) Error() string {
	return e.msg
}

func Example() {
	flag.Set("v", "4")
	flag.Parse()

	natslogr.InitFlags(nil)
	defer natslogr.Flush()

	logger := natslogr.NewLogger().
		WithName("Example").
		WithValues(natslogr.ClusterID, "example-cluster", natslogr.ClientID, "example-client", natslogr.NatsURL, stan.DefaultNatsURL, natslogr.ConnectWait, 5, natslogr.Subject, "nats-log-example")
	logger.Info("Log Example", "key", "value")

	logger.V(0).Info("Another Log Example", "logr", "nats-logr")

	logger.Error(newErrorr("An error has been occured"), "error msg", "logr", "nats-logr")

	fmt.Println("Example")
	//	Output: Example
}
