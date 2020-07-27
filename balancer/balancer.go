package balancer

import (
	"github.com/LeonhardtDavid/gobalance/configurations"
)

type Server struct {
	Targets []configurations.Target
}

type RouterCreator = func(configurations.Target) Router

// Creates different routers based in the given configuration.
func (server *Server) Run(create RouterCreator) {
	for _, target := range server.Targets {
		router := create(target)

		strategy := RouteStrategy{balancerType: target.Type, destinations: target.Destinations}

		router.Run(strategy)
	}
}
