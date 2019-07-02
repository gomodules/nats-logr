package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"gopkg.in/macaron.v1"

	stan "github.com/nats-io/go-nats-streaming"
)

func ProcessMsg(msg *stan.Msg) {
	fmt.Printf("Received on [%s]: '%s'\n", msg.Subject, msg)
	msg.Ack()
}

func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println("Close error:", err)
	}
}

func main() {
	conn, err := stan.Connect(
		"pharmer-cluster",
		"subscriber",
		stan.NatsURL(stan.DefaultNatsURL),
		stan.ConnectWait(10*time.Second),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalln("Connection lost, reason:", reason)
		}),
	)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, stan.DefaultNatsURL)
	}
	defer logCloser(conn)

	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", stan.DefaultNatsURL, "test-cluster", "subscriber")

	sub, err := conn.QueueSubscribe(
		"create-cluster",
		"test", func(msg *stan.Msg) {
			ProcessMsg(msg)
		}, stan.SetManualAckMode(), stan.DurableName("i-remember"), stan.DeliverAllAvailable(), stan.AckWait(time.Second),
	)
	defer logCloser(sub)

	m := macaron.Classic()

	m.Get("/", func() string {
		return "Hello stranger...!"
	})
	m.Run()
}
