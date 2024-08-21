package store

import "github.com/Restyx/golang-reviews-service/internal/model"

type ReviewRepositoryI interface {
	Create(*model.Review) error
	FindAll() ([]model.Review, error)
	FindOne(uint) (*model.Review, error)
	Update(*model.Review) error
	Delete(uint) error
}
