package main

import (
	"github.com/davecgh/go-spew/spew"
	stan "github.com/nats-io/go-nats-streaming"
	natslogr "gomodules.xyz/nats-logr"
)

func main() {
	logger := natslogr.New().WithValues(natslogr.ClusterID, "test-cluster", natslogr.ClientID, "publisher", natslogr.NatsURL, stan.DefaultNatsURL, natslogr.ConnectWait, 5)
	spew.Dump(logger)
}
