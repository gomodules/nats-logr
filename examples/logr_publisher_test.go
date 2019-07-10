package natslogr_test

import (
	"flag"
	"fmt"

	natslogr "gomodules.xyz/nats-logr"
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

	opts := natslogr.Options{
		ClusterID: "example-cluster",
		ClientID:  "example-client",
		Subject:   "nats-log-example",
	}
	logger := natslogr.NewLogger(opts).
		WithName("Example").
		WithValues("withKey", "withValue")
	logger.Info("Log Example", "key", "value")

	logger.V(0).Info("Another Log Example", "logr", "nats-logr")

	logger.Error(newErrorr("An error has been occured"), "error msg", "logr", "nats-logr")

	fmt.Println("Example")
	//	Output: Example
}
