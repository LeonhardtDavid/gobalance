package balancer

import (
	"fmt"
	"github.com/LeonhardtDavid/gobalance/configurations"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

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

func chooseBalancingStrategy(balancerType configurations.BalancerType, length int) func(int) int {
	switch balancerType {
	case configurations.RoundRobin:
		return func(index int) int {
			return (index + 1) % length
		}
	case configurations.Random:
		return func(index int) int {
			return rand.Intn(length)
		}
	default:
		return func(index int) int {
			return 0
		}
	}
}

func handleNextDestination(in <-chan string, out chan<- string, balancerType configurations.BalancerType, destinations []configurations.Destination) {
	index := 0
	nextDestination := chooseBalancingStrategy(balancerType, len(destinations))

	for range in {
		destination := destinations[index]
		protocol := "http"
		if destination.Secure {
			protocol = "https"
		}

		out <- fmt.Sprint(
			protocol,
			"://",
			destination.Host,
			":",
			destination.Port,
		)

		index = nextDestination(index)
	}
}

func copyHeaders(from http.Header, to http.Header) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}

func handler(timeout time.Duration, balancerType configurations.BalancerType, destinations []configurations.Destination) func(w http.ResponseWriter, r *http.Request) {
	in := make(chan string)
	out := make(chan string)

	go handleNextDestination(in, out, balancerType, destinations)

	return func(w http.ResponseWriter, r *http.Request) {
		in <- "next"
		url := <-out
		client := http.Client{
			Timeout: timeout * time.Millisecond,
		}
		request, _ := http.NewRequest(r.Method, url+r.RequestURI, r.Body)
		// request.WithContext(r.Context())

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

func initRouter(readTimeout time.Duration, writeTimeout time.Duration, maxHeaderBytes int, target configurations.Target) {
	router := mux.NewRouter()

	subRouter := router
	if target.Domain != "" {
		subRouter = router.Host(target.Domain).Subrouter()
	}

	subRouter.PathPrefix(target.Path).HandlerFunc(handler(target.Timeout, target.Type, target.Destinations))

	log.Println("Listening on port", target.Port)

	server := http.Server{
		Addr:           fmt.Sprint(":", target.Port),
		Handler:        subRouter,
		ReadTimeout:    readTimeout * time.Millisecond,
		WriteTimeout:   writeTimeout * time.Millisecond,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

// Creates different routers based in the given configuration.
func Run(config configurations.Config) {
	for _, target := range config.Targets {
		initRouter(config.ReadTimeout, config.WriteTimeout, config.MaxHeaderBytes, target)
	}
}
