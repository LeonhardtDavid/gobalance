package balancer

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRouterHandler_handler(t *testing.T) {
	client := TestClient{}
	strategy := TestStrategy{}
	w := httptest.NewRecorder()
	expectedBody := `{"some_json_key": "some_json_value"}`
	expectedHeaderName := "Some-Header"
	expectedHeaderValue := "Some header value"
	r := httptest.NewRequest(
		"POST",
		"/path",
		strings.NewReader(expectedBody),
	)
	r.Header.Set(expectedHeaderName, expectedHeaderValue)

	router := RouterHandler{
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		RequestTimeout: 10 * time.Second,
		MaxHeaderBytes: 0,
		Domain:         "www.test.com",
		Port:           8080,
		Path:           "/path",
		Client:         &client,
	}

	handlerFunc := router.handler(strategy)

	handlerFunc(w, r)

	if client.options == nil {
		t.Error("Method execute on client was unused")
	}
	if strategy.error != nil {
		t.Error("There was an error handling the next destination")
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(client.options.body)
	body := buf.String()

	if body != expectedBody {
		t.Error("Invalid body")
	}
	if client.options.header.Get(expectedHeaderName) != expectedHeaderValue {
		t.Error("Invalid header")
	}
	if client.options.method != "POST" {
		t.Error("Invalid method")
	}
	if client.options.url != "https://www.test.com:443/path" {
		t.Error("Invalid URL")
	}
	if client.options.timeout != 10 * time.Second {
		t.Error("Invalid timeout")
	}
}

type TestClient struct {
	options *ClientOptions
}

func (c *TestClient) execute(options ClientOptions) (*http.Response, error) {
	c.options = &options

	response := http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
	}

	return &response, nil
}

type TestStrategy struct {
	error error
}

func (s TestStrategy) handleNextDestination(in <-chan string, out chan<- string) error {
	for received := range in {
		if received != "next" {
			s.error = errors.New("expected next")
			return s.error
		}

		out <- "https://www.test.com:443"
	}

	return nil
}
