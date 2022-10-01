package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	_ "embed"

	"github.com/r3labs/sse"
)

//go:embed testcases/events.json
var testcasesEvents []byte

var testStreamId = "stream"

func TestEvents(t *testing.T) {
	// events cannot be tested with prism.Mock since it does
	// not support sse streams.

	server := sse.New()
	server.CreateStream(testStreamId)

	doneCh := make(chan struct{})

	handler := func(m *http.ServeMux) {
		m.HandleFunc("/eth/v1/events", func(w http.ResponseWriter, r *http.Request) {

			// r3labs/sse requires an url field 'stream' with the name of the stream to subscribe.
			// Unsure if that is part of the official sse spec or something specific of r3labs/sse
			q := r.URL.Query()
			q.Set("stream", testStreamId)
			r.URL.RawQuery = q.Encode()

			close(doneCh)
			server.HTTPHandler(w, r)
		})
	}

	addr := newMockHttpServer(t, handler)
	clt := New("http://"+addr, WithLogger(log.New(os.Stdout, "", 0)), WithUntrackedKeys())

	errCh := make(chan error)
	objCh := make(chan interface{})

	go func() {
		err := clt.Events(context.Background(), []string{"head"}, func(obj interface{}) {
			objCh <- obj
		})
		errCh <- err
	}()

	// wait for the request to be made
	select {
	case <-doneCh:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout to ping sse server")
	}

	var cases []struct {
		Event string
		Data  json.RawMessage
	}
	if err := json.Unmarshal(testcasesEvents, &cases); err != nil {
		t.Fatal(err)
	}

	for indx, c := range cases {
		// we need to remove any '\n\t' since it is used by sse to split
		// the message data during transport
		data := string(c.Data)
		data = strings.Replace(data, "\n", "", -1)
		data = strings.Replace(data, "\t", "", -1)

		server.Publish(testStreamId, &sse.Event{
			Event: []byte(c.Event),
			Data:  []byte(data),
		})

		select {
		case <-errCh:
			t.Fatal("sse client closed unexpectedly")
		case <-objCh:
		case <-time.After(1 * time.Second):
			t.Fatalf("message %d not received", indx)
		}
	}
}
