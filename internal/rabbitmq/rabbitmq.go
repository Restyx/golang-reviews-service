package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
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

	if err := channel.ExchangeDeclare("reviews", "topic", true, false, false, false, nil); err != nil {
		return err
	}

	if err := channel.Qos(1, 0, false); err != nil {
		return err
	}

	rmq.Connection = connection
	rmq.Channel = channel

	return nil
}

func (rmq *Rabbitmq) Close() {
	rmq.Connection.Close()
}
