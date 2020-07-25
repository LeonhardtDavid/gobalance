package balancer

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Router struct {
	readTimeout    time.Duration
	writeTimeout   time.Duration
	maxHeaderBytes int
	requestTimeout time.Duration
	domain         string
	port           int
	path           string
}

func handleError(w http.ResponseWriter) {
	http.Error(w, "Error accessing service", http.StatusBadGateway)
}

func handleResponse(w http.ResponseWriter, response *http.Response) {
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		log.Println("Error parsing response")
		handleError(w)
	} else {
		w.WriteHeader(response.StatusCode)
		copyHeaders(response.Header, w.Header())
		w.Write(body)
	}
}

func copyHeaders(from http.Header, to http.Header) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}

func (router *Router) handler(strategy RouteStrategy) func(w http.ResponseWriter, r *http.Request) {
	in := make(chan string)
	out := make(chan string)

	go strategy.handleNextDestination(in, out)

	return func(w http.ResponseWriter, r *http.Request) {
		in <- "next"
		url := <-out
		client := http.Client{
			Timeout: router.requestTimeout * time.Millisecond,
		}
		request, _ := http.NewRequest(r.Method, url+r.RequestURI, r.Body)
		request.WithContext(r.Context())

		copyHeaders(r.Header, request.Header)

		if response, err := client.Do(request); err != nil {
			log.Println("Error calling service", err)
			handleError(w)
		} else {
			defer response.Body.Close()
			handleResponse(w, response)
		}
	}
}

func (router *Router) Run(strategy RouteStrategy) {
	muxRouter := mux.NewRouter()

	subRouter := muxRouter
	if router.domain != "" {
		subRouter = muxRouter.Host(router.domain).Subrouter()
	}

	subRouter.PathPrefix(router.path).HandlerFunc(router.handler(strategy))

	log.Println("Listening on port", router.port)

	server := http.Server{
		Addr:           fmt.Sprint(":", router.port),
		Handler:        subRouter,
		ReadTimeout:    router.readTimeout * time.Millisecond,
		WriteTimeout:   router.writeTimeout * time.Millisecond,
		MaxHeaderBytes: router.maxHeaderBytes,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}
