package main

import (
	"github.com/LeonhardtDavid/gobalance/balancer"
	"github.com/LeonhardtDavid/gobalance/configurations"
	"time"
)

func main() {
	config := configurations.Load()
	server := balancer.Server{Targets: config.Targets}

	server.Run(func(target configurations.Target) balancer.Router {
		client := balancer.HttpClient{}

		return &balancer.RouterHandler{
			ReadTimeout:    config.ReadTimeout * time.Millisecond,
			WriteTimeout:   config.WriteTimeout * time.Millisecond,
			RequestTimeout: target.Timeout * time.Millisecond,
			MaxHeaderBytes: config.MaxHeaderBytes,
			Domain:         target.Domain,
			Port:           target.Port,
			Path:           target.Path,
			Client:         client,
		}
	})

	select {}
}
