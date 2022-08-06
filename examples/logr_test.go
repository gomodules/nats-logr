package natslogr_test

import (
	"flag"
	"fmt"

	natslogr "gomodules.xyz/nats-logr"
)

/*
nats-server -p 4222 -m 8222 -js --user myname --pass password

nats sub msg.test --user myname --password password --connection-name=demo
*/
func Example() {
	// Set a user and plain text password
	nc, err := nats.Connect("127.0.0.1", nats.UserInfo("myname", "password"))
	if err != nil {
		log.Fatal(err)
	}
	// defer nc.Close()
	defer nc.Drain()

	l := natslogr.NewAsync(nc, "msg.test")
	l.Info("default info log", "stringVal", "value", "intVal", 12345)
	l.V(0).Info("V(0) info log", "stringVal", "value", "intVal", 12345)
	l.Error(fmt.Errorf("an error"), "error log", "stringVal", "value", "intVal", 12345)

	select {}
}
