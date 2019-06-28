package natslog

import stan "github.com/nats-io/go-nats-streaming"

func publishToNatsServer(conn stan.Conn, subject string, data []byte) error {
	return conn.Publish(subject, data)
}
