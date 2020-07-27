package balancer

import (
	"fmt"
	"github.com/LeonhardtDavid/gobalance/configurations"
	"math/rand"
	"time"
)

type Strategy interface {
	handleNextDestination(in <-chan string, out chan<- string) error
}

type RouteStrategy struct {
	balancerType configurations.BalancerType
	destinations []configurations.Destination
}

func (rs RouteStrategy) chooseBalancingStrategy() func(int) int {
	length := len(rs.destinations)

	switch rs.balancerType {
	case configurations.RoundRobin:
		return func(index int) int {
			return (index + 1) % length
		}
	case configurations.Random:
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		return func(index int) int {
			return random.Intn(length)
		}
	default:
		return func(index int) int {
			return 0
		}
	}
}

func (rs RouteStrategy) handleNextDestination(in <-chan string, out chan<- string) error {
	index := 0
	nextDestination := rs.chooseBalancingStrategy()

	for range in {
		destination := rs.destinations[index]
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

	return fmt.Errorf("channel in has been closed")
}
