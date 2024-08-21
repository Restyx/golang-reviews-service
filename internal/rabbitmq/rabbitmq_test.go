package rabbitmq_test

import (
	"os"
	"testing"

	"github.com/Restyx/golang-reviews-service/internal/rabbitmq"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	amqpUser string
	amqpPass string
	amqpHost string
	amqpPort string
)

func TestMain(m *testing.M) {
	godotenv.Load()

	if amqpUser = os.Getenv("RABBITMQ_USER"); amqpUser == "" {
		amqpUser = "admin"
	}
	if amqpPass = os.Getenv("RABBITMQ_PASS"); amqpPass == "" {
		amqpPass = "admin"
	}
	if amqpHost = os.Getenv("RABBITMQ_HOST"); amqpHost == "" {
		amqpHost = "localhost"
	}
	if amqpPort = os.Getenv("RABBITMQ_PORT"); amqpPort == "" {
		amqpPort = "5672"
	}

	os.Exit(m.Run())
}

func TestConnect(t *testing.T) {
	testcases := []struct {
		name        string
		connect     func(*rabbitmq.Rabbitmq) error
		expectError bool
	}{
		{
			name: "valid",
			connect: func(rmq *rabbitmq.Rabbitmq) error {
				return rmq.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
			},
			expectError: false,
		},
		{
			name: "invalid user",
			connect: func(rmq *rabbitmq.Rabbitmq) error {
				return rmq.Connect("invalid", amqpPass, amqpHost, amqpPort)
			},
			expectError: true,
		},
		{
			name: "invalid password",
			connect: func(rmq *rabbitmq.Rabbitmq) error {
				return rmq.Connect(amqpUser, "invalid", amqpHost, amqpPort)
			},
			expectError: true,
		},
		{
			name: "invalid host",
			connect: func(rmq *rabbitmq.Rabbitmq) error {
				return rmq.Connect(amqpUser, amqpPass, "invalid", amqpPort)
			},
			expectError: true,
		},
		{
			name: "invalid port",
			connect: func(rmq *rabbitmq.Rabbitmq) error {
				return rmq.Connect(amqpUser, amqpPass, amqpHost, "invalid")
			},
			expectError: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			rmq := rabbitmq.New()
			err := testcase.connect(rmq)

			if !testcase.expectError {
				assert.NoError(t, err)
				assert.NotNil(t, rmq.Channel)
				assert.NotNil(t, rmq.CloseConnection)
				rmq.CloseConnection()
			} else {
				assert.Error(t, err)
				assert.Nil(t, rmq.Channel)
				assert.Nil(t, rmq.CloseConnection)
			}
		})
	}
}
