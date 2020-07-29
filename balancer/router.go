package balancer

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Router interface {
	Run(strategy Strategy)
}

type RouterHandler struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestTimeout time.Duration
	MaxHeaderBytes int
	Domain         string
	Port           int
	Path           string
	Client         Client
}

func handleError(w http.ResponseWriter) {
	http.Error(w, "Error accessing service", http.StatusBadGateway)
}

func handleResponseError(w http.ResponseWriter) {
	http.Error(w, "Service response error", http.StatusBadGateway)
}

func copyHeaders(from http.Header, to http.Header) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}

func handleResponse(w http.ResponseWriter, response *http.Response) {
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Println("Error parsing response")
		handleResponseError(w)
	} else {
		copyHeaders(response.Header, w.Header())
		w.WriteHeader(response.StatusCode)
		w.Write(body)
	}
}

func (router *RouterHandler) handler(strategy Strategy) http.HandlerFunc {
	in := make(chan string)
	out := make(chan string)

	go func() {
		log.Fatalln(strategy.handleNextDestination(in, out))
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		in <- "next"
		url := <-out

		options := ClientOptions{
			timeout: router.RequestTimeout,
			ctx:     r.Context(),
			method:  r.Method,
			url:     url + r.RequestURI,
			body:    r.Body,
			header:  r.Header,
		}

		if response, err := router.Client.execute(options); err != nil {
			log.Println("Error calling service", err)
			handleError(w)
		} else {
			defer response.Body.Close()
			handleResponse(w, response)
		}
	}
}

func (router *RouterHandler) Run(strategy Strategy) {
	muxRouter := mux.NewRouter()

	subRouter := muxRouter
	if router.Domain != "" {
		subRouter = muxRouter.Host(router.Domain).Subrouter()
	}

	subRouter.PathPrefix(router.Path).HandlerFunc(router.handler(strategy))

	log.Println("Listening on Port", router.Port)

	server := http.Server{
		Addr:           fmt.Sprint(":", router.Port),
		Handler:        subRouter,
		ReadTimeout:    router.ReadTimeout,
		WriteTimeout:   router.WriteTimeout,
		MaxHeaderBytes: router.MaxHeaderBytes,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}
