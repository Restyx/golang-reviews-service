package main_test

import (
	"encoding/json"
	"testing"

	"github.com/Restyx/golang-reviews-service/internal/messagehandler"
	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestMain_Create(t *testing.T) {
	rmq := prepareTest(t)
	defer rmq.Close()

	queue, err := rmq.Channel.QueueDeclare("reviews-create", false, false, false, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name         string
		inputReview  *model.Review
		expectedCode int32
	}{
		{
			name:        "valid",
			inputReview: model.TestReview(t),
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			body, err := json.Marshal(testcase.inputReview)
			assert.NoError(t, err)

			err = publishEvent(t, rmq, queue, body)
			assert.NoError(t, err)
		})

	}
}

func TestMain_Update(t *testing.T) {
	rmq := prepareTest(t)
	defer rmq.Close()

	queue, err := rmq.Channel.QueueDeclare("reviews-update", false, false, false, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name         string
		inputReview  *model.Review
		expectedCode int32
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				ID:          57,
				Author:      "updated_mail@example.com",
				Rating:      5,
				Title:       "Review Title",
				Description: "",
			},
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			body, err := json.Marshal(testcase.inputReview)
			assert.NoError(t, err)

			err = publishEvent(t, rmq, queue, body)
			assert.NoError(t, err)
		})

	}
}

func TestMain_Delete(t *testing.T) {
	rmq := prepareTest(t)
	defer rmq.Close()

	queue, err := rmq.Channel.QueueDeclare("reviews-delete", false, false, false, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name         string
		inputId      int
		expectedCode int32
	}{
		{
			name:    "valid",
			inputId: 58,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			body, err := encodeId(t, testcase.inputId)
			assert.NoError(t, err)

			err = publishEvent(t, rmq, queue, body)
			assert.NoError(t, err)
		})

	}
}

func prepareTest(t *testing.T) *rabbitmq.Rabbitmq {
	t.Helper()

	config := messagehandler.NewConfig()

	config.RmqHost = "localhost"
	config.RmqPort = "5672"
	config.RmqUser = "admin"
	config.RmqPassword = "admin"

	rmq := rabbitmq.New()

	err := rmq.Connect(config.RmqUser, config.RmqPassword, config.RmqHost, config.RmqPort)
	if err != nil {
		t.Fatal(err)
	}

	return rmq
}

func publishEvent(t *testing.T, rmq *rabbitmq.Rabbitmq, queue amqp091.Queue, body []byte) error {
	t.Helper()

	return rmq.Channel.Publish("reviews", queue.Name, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        []byte(body),
	})
}

func encodeId(t *testing.T, id int) ([]byte, error) {
	t.Helper()

	body := &struct {
		Id int `json:"id"`
	}{
		Id: id,
	}

	return json.Marshal(body)
}
