package messagehandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Restyx/golang-reviews-service/api/schemas"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	createReviewPattern string = "create-review"
	updateReviewPattern string = "update-review"
	deleteReviewPattern string = "delete-review"
	readReviewPattern   string = "get-review"
	readReviewsPattern  string = "get-reviews"
)

type messageRouter struct {
	logger  logrus.Logger
	service ServiceI
	Channel *amqp091.Channel
}

func New(service ServiceI, channel *amqp091.Channel) *messageRouter {
	return &messageRouter{
		logger:  *logrus.New(),
		service: service,
		Channel: channel,
	}
}

func (r *messageRouter) HandleMessages(messages <-chan amqp091.Delivery) {
	for msg := range messages {
		message := &schemas.Message{}
		json.Unmarshal(msg.Body, message)
		r.logger.Infof("received message with pattern: '%s' and body: %+v", message.Pattern, message.Data)

		switch message.Pattern {
		case readReviewPattern:
			review, err := r.service.ReadOne(int(message.Data.ID))
			r.rejectOnError(msg, err)

			parsedReview, err := json.Marshal(review)
			r.rejectOnError(msg, err)

			err = r.rpcResponse(msg.ReplyTo, msg.CorrelationId, parsedReview)
			r.rejectOnError(msg, err)

			msg.Ack(false)

		case readReviewsPattern:
			review, err := r.service.ReadAll()
			r.rejectOnError(msg, err)

			parsedReviews, err := json.Marshal(review)
			r.rejectOnError(msg, err)

			err = r.rpcResponse(msg.ReplyTo, msg.CorrelationId, parsedReviews)
			r.rejectOnError(msg, err)
			msg.Ack(false)

		case createReviewPattern:
			err := r.service.Create(&message.Data)
			r.rejectOnError(msg, err)
			msg.Ack(false)

		case updateReviewPattern:
			err := r.service.Update(&message.Data)
			r.rejectOnError(msg, err)

			msg.Ack(false)

		case deleteReviewPattern:
			err := r.service.Delete(int(message.Data.ID))
			r.rejectOnError(msg, err)

			msg.Ack(false)

		default:
			r.logger.Error("invalid message pattern: ", message.Pattern)
			msg.Nack(false, false)
		}
	}
}

func (r *messageRouter) rejectOnError(msg amqp091.Delivery, err error) {
	if err != nil {
		r.logger.Error(err)
		msg.Nack(false, false)
	}
}

func (r *messageRouter) rpcResponse(replyTo string, correlationId string, parsedBody []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.Channel.PublishWithContext(ctx,
		"",
		replyTo,
		false,
		false,
		amqp091.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			Body:          parsedBody,
		},
	)
}
