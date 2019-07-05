package natslogr_test

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

func processMsg(msg *stan.Msg) {
	fmt.Printf("Received on [%s]: '%s'\n", msg.Subject, msg)
	msg.Ack()
}

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println("Close error:", err)
	}
}

func Example_subscribe() {
	conn, err := stan.Connect(
		"example-cluster",
		"example-client-2",
		stan.NatsURL(stan.DefaultNatsURL),
		stan.ConnectWait(5*time.Second),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalln("Connection lost, reason: ", reason)
		}),
	)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, stan.DefaultNatsURL)
	}
	defer logCloser(conn)

	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", stan.DefaultNatsURL, "example-cluster", "example-client-2")

	sub, err := conn.QueueSubscribe(
		"nats-log-example",
		"test", func(msg *stan.Msg) {
			processMsg(msg)
		}, stan.SetManualAckMode(), stan.DurableName("i-remember"), stan.DeliverAllAvailable(), stan.AckWait(time.Second),
	)
	defer logCloser(sub)

}
