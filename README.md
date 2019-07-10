# nats-logr
A [logr](https://github.com/go-logr/logr) implementation using https://nats.io

Usage
---

### Code Example


### Description

To use [nats-logr](https://github.com/gomodules/nats-logr), you have to do the following:

- Run the following in a window:
	```$ nats-streaming-server --cluster_id=example-cluster```
- Run the following `.go` file in second window:
```go
package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
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

func main() {
	conn, err := stan.Connect(
		"example-cluster",
		"example-client-2",
		stan.NatsURL(stan.DefaultNatsURL),
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
	
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

```
- Finally, run the following code in another window: 
```go
package main

import (
	"errors"
	"flag"
	"fmt"

	natslogr "gomodules.xyz/nats-logr"

	"github.com/nats-io/stan.go"
)

func main() {
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

	logger.Error(errors.New("an error has been occured"), "error msg", "logr", "nats-logr")
}

```

Now, you will see the logs in the second window.

## Acknowledgement
The logger parts of this library has been adapted from [k8s.io/klog](https://github.com/kubernetes/klog) which itself is a fork of [golang/glog](https://github.com/golang/glog).
