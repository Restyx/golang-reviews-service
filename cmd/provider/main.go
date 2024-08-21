package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Restyx/golang-reviews-service/api/schemas"
	"github.com/Restyx/golang-reviews-service/internal/messagehandler"
	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/rabbitmq"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
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

	rmq := rabbitmq.New()

	err = rmq.Connect(config.RmqUser, config.RmqPassword, config.RmqHost, config.RmqPort)
	failOnError(err, "failed to rmq connection")

	defer rmq.CloseConnection()

	queue, err := rmq.Channel.QueueDeclare("reviews_queue", true, false, false, false, nil)
	failOnError(err, "failed to declare a queue")

	replyQueue, err := rmq.Channel.QueueDeclare("", false, false, true, false, nil)
	failOnError(err, "failed to declare a reply queue")

	msgs, err := rmq.Channel.Consume(replyQueue.Name, "", true, false, false, false, nil)
	failOnError(err, "failed to consume from reply queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(schemas.Message{
		Pattern: "get-reviews",
		// Data:    model.Review{ID: 226},
	})

	failOnError(err, "failed to marshal")

	corrId := uuid.New().String()

	err = rmq.Channel.PublishWithContext(ctx, "", queue.Name, false, false, amqp.Publishing{
		DeliveryMode:  amqp.Persistent,
		CorrelationId: corrId,
		ReplyTo:       replyQueue.Name,
		ContentType:   "application/json",
		Body:          []byte(body),
	})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s, corID: %s", body, corrId)

	for msg := range msgs {
		if corrId == msg.CorrelationId {
			review := []model.Review{}
			err := json.Unmarshal(msg.Body, &review)
			failOnError(err, "Failed to convert body")
			log.Printf("%+v", review)
			break
		}
	}
}
