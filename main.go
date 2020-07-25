package main

import (
	"github.com/LeonhardtDavid/gobalance/balancer"
	"github.com/LeonhardtDavid/gobalance/configurations"
)

func main() {
	config := configurations.Load()

	balancer.Run(config)

	select {}
}
