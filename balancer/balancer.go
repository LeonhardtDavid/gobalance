package balancer

import (
	"github.com/LeonhardtDavid/gobalance/configurations"
)

type Server struct {
	Config configurations.Config
}

// Creates different routers based in the given configuration.
func (server *Server) Run() {
	config := server.Config
	for _, target := range config.Targets {
		router := Router{
			readTimeout:    config.ReadTimeout,
			writeTimeout:   config.WriteTimeout,
			maxHeaderBytes: config.MaxHeaderBytes,
			requestTimeout: target.Timeout,
			domain:         target.Domain,
			port:           target.Port,
			path:           target.Path,
		}

		strategy := RouteStrategy{balancerType: target.Type, destinations: target.Destinations}

		router.Run(strategy)
	}
}
