package natslogr_test

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/nats-io/stan.go"
)

func processMsg(msg *stan.Msg) {
	fmt.Printf("Received on [%s]: '%s'\n", msg.Subject, msg)
	msg.Ack()
}

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Close error: %v", err)
		os.Exit(1)
	}
}

func Example_subscribe() {
	conn, err := stan.Connect(
		"example-cluster",
		"example-client-2",
		stan.NatsURL(stan.DefaultNatsURL),
		stan.ConnectWait(5*time.Second),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			fmt.Fprintf(os.Stderr, "Connection lost, reason: ", reason)
			os.Exit(1)
		}),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, stan.DefaultNatsURL)
		os.Exit(1)
	}
	defer logCloser(conn)

	fmt.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", stan.DefaultNatsURL, "example-cluster", "example-client-2")

	sub, err := conn.QueueSubscribe(
		"nats-log-example",
		"test", func(msg *stan.Msg) {
			processMsg(msg)
		}, stan.SetManualAckMode(), stan.DurableName("i-remember"), stan.DeliverAllAvailable(), stan.AckWait(time.Second),
	)
	defer logCloser(sub)

}
