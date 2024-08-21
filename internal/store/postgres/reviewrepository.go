package postgres

import (
	"fmt"

	"github.com/Restyx/golang-reviews-service/internal/model"
)

type ReviewRepository struct {
	store *Store
}

func (r *ReviewRepository) Create(review *model.Review) error {
	if err := review.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO reviews (author, rating, title, description) VALUES ($1, $2, $3, $4) RETURNING id",
		review.Author,
		review.Rating,
		review.Title,
		review.Description,
	).Scan(&review.ID)
}

func (r *ReviewRepository) FindAll() ([]model.Review, error) {
	reviews := make([]model.Review, 0)

	rows, err := r.store.db.Query("SELECT id, author, rating, title, description FROM reviews")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		review := model.Review{}

		if err := rows.Scan(&review.ID, &review.Author, &review.Rating, &review.Title, &review.Description); err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) FindOne(id uint) (*model.Review, error) {
	review := &model.Review{}
	if err := r.store.db.QueryRow("SELECT id, author, rating, title, description FROM reviews WHERE id=$1", id).Scan(&review.ID, &review.Author, &review.Rating, &review.Title, &review.Description); err != nil {
		return nil, err
	}

	return review, nil
}

func (r *ReviewRepository) Update(updateReview *model.Review) error {
	if updateReview.ID == 0 {
		return fmt.Errorf("id is required")
	}

	if err := updateReview.Validate(); err != nil {
		return err
	}

	sqlQuery := `UPDATE reviews
	SET	
	author = COALESCE(NULLIF($2, ''), author), 
	rating = COALESCE(NULLIF($3, 0), rating), 
	title = COALESCE(NULLIF($4, ''), title), 
	description = COALESCE(NULLIF($5, ''), description)
	WHERE id = $1
	RETURNING id, author, rating, title, description`

	return r.store.db.QueryRow(sqlQuery, updateReview.ID, updateReview.Author, updateReview.Rating, updateReview.Title, updateReview.Description).Scan(&updateReview.ID, &updateReview.Author, &updateReview.Rating, &updateReview.Title, &updateReview.Description)
}

func (r *ReviewRepository) Delete(id uint) error {
	if err := r.store.db.QueryRow("DELETE FROM reviews WHERE id=$1 RETURNING id", id).Scan(&id); err != nil {
		return err
	}

	return nil
}
