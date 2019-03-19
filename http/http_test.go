package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDoc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`foo`))
	}))

	_, err := GetDoc(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
