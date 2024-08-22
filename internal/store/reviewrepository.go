package store

import "github.com/Restyx/golang-reviews-service/internal/model"

type ReviewRepositoryI interface {
	Create(*model.Review) (int64, error)
	FindOne(int) (*model.Review, error)
	FindAll() ([]model.Review, error)
	Update(*model.Review) (int64, error)
	Delete(int) (int64, error)
}
