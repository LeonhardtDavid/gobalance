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

type strategyFunc = func() configurations.Destination

type destinationPicker interface {
	picker() strategyFunc
}

type roundRobinPicker struct {
	destinations []configurations.Destination
}

func (p roundRobinPicker) picker() strategyFunc {
	length := len(p.destinations)
	currentIndex := 0

	return func() configurations.Destination {
		index := currentIndex
		currentIndex = (currentIndex + 1) % length
		return p.destinations[index]
	}
}

type randomPicker struct {
	destinations []configurations.Destination
}

func (p randomPicker) picker() strategyFunc {
	length := len(p.destinations)
	currentIndex := 0
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	return func() configurations.Destination {
		index := currentIndex
		currentIndex = random.Intn(length)
		return p.destinations[index]
	}
}

type defaultPicker struct {
	destination configurations.Destination
}

func (p defaultPicker) picker() strategyFunc {
	return func() configurations.Destination {
		return p.destination
	}
}

func (rs RouteStrategy) nextDestinationStrategy() destinationPicker {
	switch rs.balancerType {
	case configurations.RoundRobin:
		return roundRobinPicker{
			destinations: rs.destinations,
		}
	case configurations.Random:
		return randomPicker{
			destinations: rs.destinations,
		}
	default:
		return defaultPicker{
			destination: rs.destinations[0],
		}
	}
}

func (rs RouteStrategy) handleNextDestination(in <-chan string, out chan<- string) error {
	nextDestination := rs.nextDestinationStrategy().picker()

	for range in {
		destination := nextDestination()
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
	}

	return fmt.Errorf("channel in has been closed")
}
