package schemas

import "github.com/Restyx/golang-reviews-service/internal/model"

type Message struct {
	Pattern string       `json:"pattern"`
	Data    model.Review `json:"data"`
}
