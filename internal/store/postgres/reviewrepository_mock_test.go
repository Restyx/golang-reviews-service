package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPostgresReviewRepositoryMock_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := postgres.New(db)

	type mockBehavior func(review *model.Review)

	testTable := []struct {
		name         string
		inputReview  *model.Review
		mockBehavior mockBehavior
		expectedID   uint
		expectError  bool
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
			mockBehavior: func(review *model.Review) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO reviews").WithArgs(review.Author, review.Rating, review.Title, review.Description).WillReturnRows(rows)
			},
			expectedID: 1,
		},
		{
			name: "invalid id",
			inputReview: &model.Review{
				Author:      "invalid",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
			mockBehavior: func(review *model.Review) {},
			expectError:  true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior(testcase.inputReview)

			err := store.Review().Create(testcase.inputReview)

			if testcase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, testcase.inputReview.ID, testcase.expectedID)
			}
		})
	}
}
