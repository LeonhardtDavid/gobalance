package balancer

import (
	"github.com/LeonhardtDavid/gobalance/configurations"
)

type Server struct {
	Targets []configurations.Target
}

// Creates different routers based in the given configuration.
func (server *Server) Run(f func(configurations.Target) Router) {
	for _, target := range server.Targets {
		router := f(target)

		strategy := RouteStrategy{balancerType: target.Type, destinations: target.Destinations}

		router.Run(strategy)
	}
}
