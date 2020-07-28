package balancer

import (
	"github.com/LeonhardtDavid/gobalance/configurations"
	"testing"
)

func TestServerRun_EmptyTargets(t *testing.T) {
	var targets []configurations.Target
	count := 0

	server := Server{Targets: targets}

	server.Run(func(target configurations.Target) Router {
		count++
		return struct {Router}{}
	})

	if count != 0 {
		t.Error(count, "routes created, expected 0")
	}
}

func TestServerRun(t *testing.T) {
	targets := []configurations.Target{
		{},
		{},
	}
	count := 0
	router := TestRouter{}

	server := Server{Targets: targets}

	server.Run(func(target configurations.Target) Router {
		count++
		return &router
	})

	if count != 2 {
		t.Error(count, "routes created, expected 0")
	}

	if len(router.strategy) != 2 {
		t.Error("2 strategies spected")
	}
}

type TestRouter struct {
	strategy []Strategy
}

func (r *TestRouter) Run(strategy Strategy) {
	r.strategy = append(r.strategy, strategy)
}
