package postgres_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestReviewRepository_Create(t *testing.T) {
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
			name: "invalid email",
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

			id, err := store.Review().Create(testcase.inputReview)

			if testcase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, id, testcase.expectedID)
			}
		})
	}
}

func TestReviewRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := postgres.New(db)

	type mockBehavior func(review *model.Review)

	testTable := []struct {
		name           string
		inputReview    *model.Review
		mockBehavior   mockBehavior
		expectRowCount int
		expectError    bool
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				ID:          1,
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
			mockBehavior: func(review *model.Review) {
				mockStmt := mock.ExpectPrepare("UPDATE reviews")
				mockRes := mockStmt.ExpectExec().WithArgs(review.ID, review.Author, review.Rating, review.Title, review.Description)
				mockRes.WillReturnResult(sqlmock.NewResult(1, 1))

			},
			expectRowCount: 1,
		},
		{
			name: "invalid id",
			inputReview: &model.Review{
				ID:          413,
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
			mockBehavior: func(review *model.Review) {
				mockStmt := mock.ExpectPrepare("UPDATE reviews")
				mockRes := mockStmt.ExpectExec().WithArgs(review.ID, review.Author, review.Rating, review.Title, review.Description)
				mockRes.WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectError: false,
		},
		{
			name: "emty id",
			inputReview: &model.Review{
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
			mockBehavior: func(review *model.Review) {
			},
			expectError:    true,
			expectRowCount: 0,
		},
		{
			name: "empty fields",
			inputReview: &model.Review{
				ID:          1,
				Author:      "",
				Rating:      4,
				Title:       "review title",
				Description: "review description",
			},
			mockBehavior: func(review *model.Review) {
				mockStmt := mock.ExpectPrepare("UPDATE reviews")
				mockRes := mockStmt.ExpectExec().WithArgs(review.ID, review.Author, review.Rating, review.Title, review.Description)
				mockRes.WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectRowCount: 1,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior(testcase.inputReview)

			// originalInput := *testcase.inputReview

			rowCount, err := store.Review().Update(testcase.inputReview)

			assert.EqualValues(t, testcase.expectRowCount, rowCount)

			if testcase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReviewRepository_Delete(t *testing.T) {
	// TODO
}

func TestReviewRepository_FindOne(t *testing.T) {
	// TODO
}

func TestReviewRepository_FindAll(t *testing.T) {
	// TODO
}
