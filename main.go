package main

import (
	"github.com/LeonhardtDavid/gobalance/balancer"
	"github.com/LeonhardtDavid/gobalance/configurations"
)

func main() {
	config := configurations.Load()
	server := balancer.Server{Config: config}

	server.Run()

	select {}
}
