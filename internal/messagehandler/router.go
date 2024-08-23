package messagehandler

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	readReviewPattern   string = "reviews-get-one"
	readReviewsPattern  string = "reviews-get-all"
	createReviewPattern string = "reviews-create"
	updateReviewPattern string = "reviews-update"
	deleteReviewPattern string = "reviews-delete"
)

type Server struct {
	logger  logrus.Logger
	service ServiceI
	Channel *amqp091.Channel
}

func New(service ServiceI, channel *amqp091.Channel) *Server {
	return &Server{
		logger:  *logrus.New(),
		service: service,
		Channel: channel,
	}
}

func (s *Server) HandleMessages(messages <-chan amqp091.Delivery) {
	for msg := range messages {
		s.logger.Infof("received message key: '%s'", msg.RoutingKey)

		var (
			nack   bool
			reason error
			body   []byte
		)

		switch msg.RoutingKey {
		case readReviewPattern:
			id, err := DecodeId(msg.Body)
			if err != nil {
				nack = true
				reason = err
				break
			}

			review, err := s.service.ReadOne(id)
			if err != nil {
				nack = true
				reason = err
				break
			}

			body, err = json.Marshal(review)
			if err != nil {
				nack = true
				reason = err
			}

		case readReviewsPattern:
			reviews, err := s.service.ReadAll()
			if err != nil {
				nack = true
				reason = err
				break

			}

			body, err = json.Marshal(reviews)
			if err != nil {
				nack = true
				reason = err
			}

		case createReviewPattern:
			review, err := DecodeReview(msg.Body)
			if err != nil {
				nack = true
				reason = err
				break
			}

			err = s.service.Create(review)
			if err != nil {
				nack = true
				reason = err
			}

		case updateReviewPattern:
			review, err := DecodeReview(msg.Body)
			if err != nil {
				nack = true
				reason = err
				break
			}

			err = s.service.Update(review)
			if err != nil {
				nack = true
				reason = err
			}

		case deleteReviewPattern:
			id, err := DecodeId(msg.Body)
			if err != nil {
				nack = true
				reason = err
				break
			}

			err = s.service.Delete(id)
			if err != nil {
				nack = true
				reason = err
			}

		default:
			nack = true
			reason = errors.New("invalid message routing key")
		}

		if msg.ReplyTo != "" {
			if nack {
				body = []byte(fmt.Sprint(reason))
			}

			s.logger.Infof("sending reply with status code %v", getStatusCode(reason))
			err := s.reply(msg, getStatusCode(reason), body)

			if err != nil {
				nack = true
				reason = err
			}
		}

		if nack {
			s.logger.Errorf("message rejected: %s", reason)
			msg.Nack(false, false)
		} else {
			s.logger.Infof("message acknowledged")
			msg.Ack(false)
		}
	}
}

func (s *Server) reply(msg amqp091.Delivery, statusCode int32, body []byte) error {
	return s.Channel.Publish(
		"",
		msg.ReplyTo,
		false,
		false,
		amqp091.Publishing{
			Headers: amqp091.Table{
				"code": statusCode,
			},
			CorrelationId: msg.CorrelationId,
			ContentType:   "application/json",
			Body:          []byte(body),
		},
	)
}

func getStatusCode(inputError error) int32 {
	var statusCode int32
	switch {
	case inputError == nil:
		statusCode = 200
	case errors.As(inputError, &store.ErrRecordNotFound):
		statusCode = 404
	case errors.As(inputError, &store.ErrFieldMissing):
		statusCode = 400
	default:
		statusCode = 500
	}

	return statusCode
}

func DecodeReview(body []byte) (*model.Review, error) {
	review := &model.Review{}

	if err := json.Unmarshal(body, review); err != nil {
		return nil, err
	}

	return review, nil
}

func DecodeReviewSlice(body []byte) ([]model.Review, error) {
	var review []model.Review

	if err := json.Unmarshal(body, &review); err != nil {
		return nil, err
	}

	return review, nil
}

func DecodeId(body []byte) (int, error) {
	var data struct {
		Id int `json:"id"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data.Id, nil
}
