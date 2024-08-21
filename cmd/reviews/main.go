package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Restyx/golang-reviews-service/internal/messagehandler"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/reviews.toml", "path to config file")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	flag.Parse()
	config := messagehandler.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	failOnError(err, "failed to initialize config file")

	err = messagehandler.Start(config)
	failOnError(err, "failed to start the service")
}
