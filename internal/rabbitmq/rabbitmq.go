package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	Channel         *amqp.Channel
	CloseConnection func() error
}

func New() *Rabbitmq {
	return &Rabbitmq{}
}

func (rmq *Rabbitmq) Connect(user string, password string, host string, port string) error {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	connection, err := amqp.Dial(address)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	rmq.Channel = channel
	rmq.CloseConnection = connection.Close

	return nil
}
