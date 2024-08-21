package postgres

import (
	"database/sql"

	"github.com/Restyx/golang-reviews-service/internal/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db               *sql.DB
	reviewRepository *ReviewRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Review() store.ReviewRepositoryI {
	if s.reviewRepository == nil {
		s.reviewRepository = &ReviewRepository{
			store: s,
		}
	}

	return s.reviewRepository
}
