package testingstorage

import (
	"fmt"

	"github.com/Restyx/golang-reviews-service/internal/model"
)

type ReviewRepository struct {
	store   *Store
	reviews map[int]*model.Review
}

func (r *ReviewRepository) Create(review *model.Review) (int, error) {
	if err := review.Validate(); err != nil {
		return 0, err
	}

	review.ID = int(len(r.reviews) + 1)

	r.reviews[review.ID] = review

	return review.ID, nil
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

func (r *ReviewRepository) Update(updatedReview *model.Review) error {
	if updatedReview.ID == 0 {
		return fmt.Errorf("id is required")
	}

	if err := updatedReview.Validate(); err != nil {
		return err
	}

	review, ok := r.reviews[int(updatedReview.ID)]
	if !ok {
		return fmt.Errorf("record not found")
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

	return nil
}

func (r *ReviewRepository) Delete(id int) error {
	if id == 0 {
		return fmt.Errorf("id is required")
	}

	_, ok := r.reviews[id]
	if !ok {
		return fmt.Errorf("record not found")
	}

	delete(r.reviews, id)

	return nil
}
