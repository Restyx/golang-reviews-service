package testingstorage

import (
	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store"
)

type Store struct {
	reviewRepository *ReviewRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Review() store.ReviewRepositoryI {
	if s.reviewRepository == nil {
		s.reviewRepository = &ReviewRepository{
			store:   s,
			reviews: make(map[uint]*model.Review),
		}
	}

	return s.reviewRepository
}
