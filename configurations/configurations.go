package configurations

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Destination struct {
	Host   string
	Port   int
	Secure bool
}

type BalancerType string

const (
	RoundRobin BalancerType = "ROUND_ROBIN"
	Random     BalancerType = "RANDOM"
)

type Target struct {
	Domain       string
	Port         int
	Path         string
	Timeout      time.Duration `mapstructure:"timeout-millis"`
	Type         BalancerType
	Destinations []Destination
}

type Config struct {
	ReadTimeout    time.Duration `mapstructure:"read-timeout-millis"`
	WriteTimeout   time.Duration `mapstructure:"write-timeout-millis"`
	MaxHeaderBytes int           `mapstructure:"max-header-bytes"`
	Targets        []Target
}

func initViper() {
	viper.SetConfigName("config")          // name of config file (without extension)
	viper.AddConfigPath("/etc/gobalance/") // path to look for the config file in
	viper.AddConfigPath(".")               // call multiple times to add many search paths, this optionally look for config in the working directory
	err := viper.ReadInConfig()            // Find and read the config file
	if err != nil {                        // Handle errors reading the config file
		log.Fatalln("Fatal error config file: ", err)
	}
}

func validateConfigFileExistence() {
	configFileName := "config.yml"
	rootConfig := "./" + configFileName
	_, errHere := os.Stat(rootConfig)
	_, errEtc := os.Stat("/etc/gobalance/" + configFileName)
	if os.IsNotExist(errHere) && os.IsNotExist(errEtc) {
		log.Println("No config file found, creating one from template")

		input, err := ioutil.ReadFile("./config.template.yml")
		if err != nil {
			log.Fatalln("Error reading template file")
		}

		err = ioutil.WriteFile(rootConfig, input, 0644)
		if err != nil {
			log.Fatal("Couldn't create config file", err)
		}
	}
}

func (config Config) validatesBalancingTypes() error {
	for _, target := range config.Targets {
		switch target.Type {
		case RoundRobin:
		case Random:
		default:
			return fmt.Errorf("invalid balancing type %s", target.Type)
		}
	}

	return nil
}

// Loads the configuration from config.yml file beeing on project root or /etc/gobalance
func Load() Config {
	validateConfigFileExistence()
	initViper()

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln("Error parsing config file")
	}

	if err := config.validatesBalancingTypes(); err != nil {
		log.Fatalln(err)
	}

	return config
}
