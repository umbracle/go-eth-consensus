package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHttp_Post(t *testing.T) {
	handler := func(m *http.ServeMux) {
		m.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Query().Get("t") {
			case "a":
				w.Write([]byte(`{"data":null}`))
			case "b":
				w.Write([]byte(`null`))
			case "c":
				w.Write([]byte(``))
			case "d":
				w.Write([]byte(`x`))
			}
		})
	}

	addr := newMockHttpServer(t, handler)
	clt := New("http://" + addr)

	var input struct{}

	require.NoError(t, clt.Post("/do?t=a", &input, nil))
	require.NoError(t, clt.Post("/do?t=b", &input, nil))
	require.NoError(t, clt.Post("/do?t=c", &input, nil))
	require.Error(t, clt.Post("/do?t=d", &input, nil))
}

func TestHttp_ErrorStatus(t *testing.T) {
	handler := func(m *http.ServeMux) {
		m.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
			errorCode := r.URL.Query().Get("t")

			switch errorCode {
			case "400":
				w.WriteHeader(http.StatusBadRequest)
			case "404":
				w.WriteHeader(http.StatusNotFound)
			case "500":
				w.WriteHeader(http.StatusInternalServerError)
			case "503":
				w.WriteHeader(http.StatusServiceUnavailable)
			}

			w.Write([]byte(`{"message": "incorrect data", "code": ` + errorCode + `}`))
		})
	}

	addr := newMockHttpServer(t, handler)
	clt := New("http://" + addr)

	require.ErrorIs(t, clt.Get("/do?t=400", nil), ErrorBadRequest)
	require.ErrorIs(t, clt.Get("/do?t=404", nil), ErrorNotFound)
	require.ErrorIs(t, clt.Get("/do?t=500", nil), ErrorInternalServerError)
	require.ErrorIs(t, clt.Get("/do?t=503", nil), ErrorServiceUnavailable)
}

func newMockHttpServer(t *testing.T, handler func(m *http.ServeMux)) string {
	m := http.NewServeMux()
	if handler != nil {
		handler(m)
	}

	addr := "127.0.0.1:8000"
	s := &http.Server{
		Addr:    addr,
		Handler: m,
	}

	go func() {
		s.ListenAndServe()
	}()
	time.Sleep(500 * time.Millisecond)

	t.Cleanup(func() {
		s.Close()
	})
	return addr
}
