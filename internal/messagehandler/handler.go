package messagehandler

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Restyx/golang-reviews-service/internal/rabbitmq"
	"github.com/Restyx/golang-reviews-service/internal/store/postgres"
)

func Start(config *Config) error {
	database, err := connectDB(config.PgUser, config.PgPassword, config.PgHost, config.PgPort, config.PgDB)
	if err != nil {
		return err
	}
	defer database.Close()

	rmq := rabbitmq.New()
	if err := rmq.Connect(config.RmqUser, config.RmqPassword, config.RmqHost, config.RmqPort); err != nil {
		return err
	}
	defer rmq.Close()

	store := postgres.New(database)

	reviewsService := NewService(store)

	reviewsRouter := New(reviewsService, rmq.Channel)

	return listen(reviewsRouter)
}

func connectDB(user, password, host, port, datatbase string) (*sql.DB, error) {
	databaseURL := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, datatbase)
	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

func listen(Router *Server) error {
	queue, err := Router.Channel.QueueDeclare("reviews_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for _, s := range []string{readReviewPattern, readReviewsPattern, createReviewPattern, updateReviewPattern, deleteReviewPattern} {
		log.Printf("Binding queue %s to exchange %s with routing key %s", queue.Name, "reviews", s)

		if err := Router.Channel.QueueBind(queue.Name, s, "reviews", false, nil); err != nil {
			return err
		}
	}

	msgs, err := Router.Channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan interface{})

	go Router.HandleMessages(msgs)
	log.Printf("[*] Waiting for logs. To exit press CTRL+C")

	<-forever

	return nil
}
