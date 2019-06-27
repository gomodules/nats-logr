package natslogr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sort"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

// Constants for Nats Streaming Server connection
const (
	ClusterID   = "cluster_id"
	ClientID    = "client_id"
	NatsURL     = "nats_url"
	ConnectWait = "connect_wait"
)

// Magic string for intermediate frames that we should ignore.
const autogeneratedFrameName = "<autogenerated>"

// Discover how many frames we need to climb to find the caller. This approach
// was suggested by Ian Lance Taylor of the Go team, so it *should* be safe
// enough (famous last words).
func framesToCaller() int {
	// 1 is the immediate caller.  3 should be too many.
	for i := 1; i < 3; i++ {
		_, file, _, _ := runtime.Caller(i + 1) // +1 for this function's frame
		if file != autogeneratedFrameName {
			return i
		}
	}
	return 1 // something went wrong, this is safe
}

// trimDuplicates will deduplicate elements provided in multiple KV tuple
// slices, whilst maintaining the distinction between where the items are
// contained.
func trimDuplicates(kvLists ...[]interface{}) [][]interface{} {
	// maintain a map of all seen keys
	seenKeys := map[interface{}]struct{}{}
	// build the same number of output slices as inputs
	outs := make([][]interface{}, len(kvLists))
	// iterate over the input slices backwards, as 'later' kv specifications
	// of the same key will take precedence over earlier ones
	for i := len(kvLists) - 1; i >= 0; i-- {
		// initialise this output slice
		outs[i] = []interface{}{}
		// obtain a reference to the kvList we are processing
		kvList := kvLists[i]

		// start iterating at len(kvList) - 2 (i.e. the 2nd last item) for
		// slices that have an even number of elements.
		// We add (len(kvList) % 2) here to handle the case where there is an
		// odd number of elements in a kvList.
		// If there is an odd number, then the last element in the slice will
		// have the value 'null'.
		for i2 := len(kvList) - 2 + (len(kvList) % 2); i2 >= 0; i2 -= 2 {
			k := kvList[i2]
			// if we have already seen this key, do not include it again
			if _, ok := seenKeys[k]; ok {
				continue
			}
			// make a note that we've observed a new key
			seenKeys[k] = struct{}{}
			// attempt to obtain the value of the key
			var v interface{}
			// i2+1 should only ever be out of bounds if we handling the first
			// iteration over a slice with an odd number of elements
			if i2+1 < len(kvList) {
				v = kvList[i2+1]
			}
			// add this KV tuple to the *start* of the output list to maintain
			// the original order as we are iterating over the slice backwards
			outs[i] = append([]interface{}{k, v}, outs[i]...)
		}
	}
	return outs
}

func flatten(keyAndValues ...interface{}) string {
	keys := make([]string, 0, len(keyAndValues))
	vals := make(map[string]interface{}, len(keyAndValues))
	for i := 0; i < len(keyAndValues); i += 2 {
		k, ok := keyAndValues[i].(string)
		if !ok {
			panic(fmt.Sprintf("key is not a string: %s", pretty(keyAndValues[i])))
		}
		var v interface{}
		if i+1 < len(keyAndValues) {
			v = keyAndValues[i+1]
		}
		keys = append(keys, k)
		vals[k] = v
	}
	sort.Strings(keys)
	buf := bytes.Buffer{}
	for i, k := range keys {
		v := vals[k]
		if i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(pretty(k))
		buf.WriteString("=")
		buf.WriteString(pretty(v))
	}
	return buf.String()
}

func pretty(value interface{}) string {
	jb, _ := json.Marshal(value)
	return string(jb)
}

func copySlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	copy(out, in)
	return out
}

func (l natsLogger) clone() natsLogger {
	return natsLogger{
		level:    l.level,
		prefix:   l.prefix,
		values:   copySlice(l.values),
		stanConn: l.stanConn,
	}
}

func getNatsInfo(keyValues ...[]interface{}) (natsInfo, bool) {
	nats := natsInfo{}

	outs := make([][]interface{}, len(keyValues))

	for i := len(keyValues) - 1; i >= 0; i-- {
		outs[i] = []interface{}{}
		keyValue := keyValues[i]

		for j := len(keyValue) - 2 + (len(keyValue) % 2); j >= 0; j -= 2 {
			key := keyValue[j]
			value := keyValue[j+1]
			fmt.Println(key, "=", value)
			if key == ClusterID {
				nats.clusterID = value.(string)
			} else if key == ClientID {
				nats.clientID = value.(string)
			} else if key == NatsURL {
				nats.natsUrl = value.(string)
			} else if key == ConnectWait {
				nats.connectWait = time.Duration(value.(int)) * time.Second
			}
		}
	}

	return nats, true
}

func connectToNatsStreamingServer(info natsInfo) stan.Conn {
	conn, err := stan.Connect(
		info.clusterID, info.clusterID, stan.NatsURL(info.natsUrl), stan.ConnectWait(info.connectWait),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalln("Connection lost, reason:", reason)
		}),
	)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, info.natsUrl)
	}
	return conn
}
