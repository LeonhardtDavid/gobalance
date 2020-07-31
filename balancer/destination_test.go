package balancer

import (
	"github.com/LeonhardtDavid/gobalance/configurations"
	"math/rand"
	"testing"
	"time"
)

func TestNextDestinationStrategy_RoundRobin(t *testing.T) {
	strategy := RouteStrategy{
		balancerType: configurations.RoundRobin,
		destinations: []configurations.Destination{
			{Host: "localhost", Port: 9000, Secure: false},
		},
	}

	switch v := strategy.nextDestinationStrategy().(type) {
	case roundRobinPicker:
	default:
		t.Errorf("Invalid picker, roundRobinPicker expected, got %T", v)
	}
}

func TestNextDestinationStrategy_Random(t *testing.T) {
	strategy := RouteStrategy{
		balancerType: configurations.Random,
		destinations: []configurations.Destination{
			{Host: "localhost", Port: 9000, Secure: false},
		},
	}

	switch v := strategy.nextDestinationStrategy().(type) {
	case randomPicker:
	default:
		t.Errorf("Invalid picker, randomPicker expected, got %T", v)
	}
}

func TestNextDestinationStrategy_Default(t *testing.T) {
	strategy := RouteStrategy{
		destinations: []configurations.Destination{
			{Host: "localhost", Port: 9000, Secure: false},
		},
	}

	switch v := strategy.nextDestinationStrategy().(type) {
	case defaultPicker:
	default:
		t.Errorf("Invalid picker, defaultPicker expected, got %T", v)
	}
}

func TestRoundRobinPicker_picker(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	destination1 := configurations.Destination{Host: "localhost", Port: rand.Intn(10000), Secure: false}
	destination2 := configurations.Destination{Host: "google.com", Port: rand.Intn(10000), Secure: true}

	picker := roundRobinPicker{
		destinations: []configurations.Destination{
			destination1,
			destination2,
		},
	}

	pickerFunc := picker.picker()

	picked1 := pickerFunc()
	picked2 := pickerFunc()
	picked3 := pickerFunc()

	if picked1 != destination1 {
		t.Error("Unexpected destination")
	}
	if picked2 != destination2 {
		t.Error("Unexpected destination")
	}
	if picked3 != destination1 {
		t.Error("Unexpected destination")
	}
}

func TestRandomPicker_picker(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	destination1 := configurations.Destination{Host: "localhost", Port: rand.Intn(10000), Secure: false}
	destination2 := configurations.Destination{Host: "google.com", Port: rand.Intn(10000), Secure: true}

	picker := randomPicker{
		destinations: []configurations.Destination{
			destination1,
			destination2,
		},
	}

	pickerFunc := picker.picker()

	picked1 := pickerFunc()
	picked2 := pickerFunc()
	picked3 := pickerFunc()

	if picked1 != destination1 && picked1 != destination2 {
		t.Error("Unexpected destination")
	}
	if picked2 != destination1 && picked2 != destination2 {
		t.Error("Unexpected destination")
	}
	if picked3 != destination1 && picked3 != destination2 {
		t.Error("Unexpected destination")
	}
}

func TestDefaultPicker_picker(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	destination := configurations.Destination{Host: "localhost", Port: rand.Intn(10000), Secure: false}

	picker := defaultPicker{
		destination: destination,
	}

	pickerFunc := picker.picker()

	picked := pickerFunc()

	if picked != destination {
		t.Error("Unexpected destination")
	}
}
