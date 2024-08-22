package testingstorage

import (
	"fmt"

	"github.com/Restyx/golang-reviews-service/internal/model"
)

type ReviewRepository struct {
	store   *Store
	reviews map[int]*model.Review
}

func (r *ReviewRepository) Create(review *model.Review) (int64, error) {
	if err := review.Validate(); err != nil {
		return 0, err
	}

	review.ID = uint(len(r.reviews) + 1)

	r.reviews[int(review.ID)] = review

	return int64(review.ID), nil
}

func (r *ReviewRepository) FindAll() ([]model.Review, error) {
	result := make([]model.Review, len(r.reviews))

	for index, value := range r.reviews {
		result[index-1] = *value
	}

	return result, nil
}

func (r *ReviewRepository) FindOne(id int) (*model.Review, error) {
	review, ok := r.reviews[id]
	if !ok {
		return nil, fmt.Errorf("record not found")
	}

	return review, nil
}

func (r *ReviewRepository) Update(updatedReview *model.Review) (int64, error) {
	if updatedReview.ID == 0 {
		return 0, fmt.Errorf("id is required")
	}

	if err := updatedReview.Validate(); err != nil {
		return 0, err
	}

	review, ok := r.reviews[int(updatedReview.ID)]
	if !ok {
		return 0, fmt.Errorf("record not found")
	}

	if updatedReview.Author != "" {
		review.Author = updatedReview.Author
	} else {
		updatedReview.Author = review.Author
	}
	if updatedReview.Rating != 0 {
		review.Rating = updatedReview.Rating
	} else {
		updatedReview.Rating = review.Rating
	}
	if updatedReview.Title != "" {
		review.Title = updatedReview.Title
	} else {
		updatedReview.Title = review.Title
	}
	if updatedReview.Description != "" {
		review.Description = updatedReview.Description
	} else {
		updatedReview.Description = review.Description
	}

	return 1, nil
}

func (r *ReviewRepository) Delete(id int) (int64, error) {
	_, ok := r.reviews[id]
	if !ok {
		return 0, nil
	}

	delete(r.reviews, id)

	return 1, nil
}
