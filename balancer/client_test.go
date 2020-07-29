package balancer

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHttpClientExecute(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()

		if r.Method == "POST" && body == `{"key":"value""}` && r.Header.Get("Some-Header") == "header content" {
			w.Header().Set("some-key", "some_value")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"response_key":"response_value""}`))
		} else {
			http.Error(w, "error", http.StatusBadRequest)
		}
	}))
	defer ts.Close()
	ctx := context.Background()

	options := ClientOptions{
		timeout : 1 * time.Second,
		ctx    : ctx,
		method : "POST",
		url    : ts.URL,
		body   : strings.NewReader(`{"key":"value""}`),
		header : http.Header{
			"Some-Header": []string{"header content"},
		},
	}

	client := HttpClient{}

	if response, err := client.execute(options); err != nil {
		t.Error("Response error", err)
	} else if response.StatusCode != http.StatusOK {
		t.Error("Expected status 200 received", response.StatusCode)
	} else if response.Header.Get("some-key") != "some_value" {
		t.Error("Unexpected response header")
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		body := buf.String()

		if body != `{"response_key":"response_value""}` {
			t.Error("Unexpected body")
		}
	}
}
