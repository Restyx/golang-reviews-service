package postgres

import (
	"fmt"

	"github.com/Restyx/golang-reviews-service/internal/model"
)

type ReviewRepository struct {
	store *Store
}

func (r *ReviewRepository) Create(review *model.Review) (int, error) {
	if err := review.Validate(); err != nil {
		return 0, err
	}

	err := r.store.db.QueryRow("INSERT INTO reviews (author, rating, title, description) VALUES ($1, $2, $3, $4) RETURNING id", review.Author, review.Rating, review.Title, review.Description).Scan(&review.ID)
	if err != nil {
		return 0, err
	}

	return review.ID, nil
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

func (r *ReviewRepository) FindOne(id int) (*model.Review, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}

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
	RETURNING author, rating, title, description`

	stmt, err := r.store.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(updateReview.ID, updateReview.Author, updateReview.Rating, updateReview.Title, updateReview.Description)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}

func (r *ReviewRepository) Delete(id int) error {
	if id == 0 {
		return fmt.Errorf("id is required")
	}

	stmt, err := r.store.db.Prepare("DELETE FROM reviews WHERE id=$1")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}
