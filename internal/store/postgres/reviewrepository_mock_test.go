package postgres_test

import (
	"fmt"
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
		expectedID   int
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
		name         string
		inputReview  *model.Review
		mockBehavior mockBehavior
		expectError  bool
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
				mock.ExpectPrepare("UPDATE reviews").ExpectExec().WithArgs(review.ID, review.Author, review.Rating, review.Title, review.Description).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectError: true,
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
			expectError: true,
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
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior(testcase.inputReview)

			err := store.Review().Update(testcase.inputReview)

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := postgres.New(db)

	type mockBehavior func(id int)

	testTable := []struct {
		name         string
		inputId      int
		mockBehavior mockBehavior
		expectError  bool
	}{
		{
			name:    "valid",
			inputId: 1,
			mockBehavior: func(id int) {
				mock.ExpectPrepare("DELETE FROM reviews").ExpectExec().WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:    "empty id",
			inputId: 0,
			mockBehavior: func(id int) {
			},
			expectError: true,
		},
		{
			name:    "invalid id",
			inputId: 112314,
			mockBehavior: func(id int) {
				mock.ExpectPrepare("DELETE FROM reviews").ExpectExec().WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 0))
			},
			expectError: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior(testcase.inputId)

			err := store.Review().Delete(testcase.inputId)

			if testcase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestReviewRepository_FindOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := postgres.New(db)

	type mockBehavior func(id int)

	testTable := []struct {
		name           string
		inputId        int
		mockBehavior   mockBehavior
		expectedReview *model.Review
	}{
		{
			name:    "valid",
			inputId: 1,
			mockBehavior: func(id int) {
				rows := mock.NewRows([]string{"id", "author", "rating", "title", "description"}).AddRow(1, "example_mail.@example.com", 3, "review title", "review description")
				mock.ExpectQuery("SELECT id, author, rating, title, description FROM reviews").WithArgs(id).WillReturnRows(rows)
			},
			expectedReview: &model.Review{
				ID:          1,
				Author:      "example_mail.@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
		},
		{
			name:    "not existing review",
			inputId: 123,
			mockBehavior: func(id int) {
				mock.ExpectQuery("SELECT FROM reviews").WithArgs(id).WillReturnError(fmt.Errorf("some errro"))
			},

			expectedReview: nil,
		},
		{
			name:    "empty id",
			inputId: 0,
			mockBehavior: func(id int) {
			},

			expectedReview: nil,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior(testcase.inputId)
			returnedReview, err := store.Review().FindOne(testcase.inputId)

			if testcase.expectedReview != nil {
				assert.NoError(t, err)
				assert.EqualValues(t, testcase.expectedReview, returnedReview)
			} else {
				assert.Error(t, err)
				assert.Nil(t, returnedReview)
			}
		})
	}
}

func TestReviewRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := postgres.New(db)

	type mockBehavior func()

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		expectedLen  int
	}{
		{
			name: "1 review",
			mockBehavior: func() {
				rows := mock.NewRows([]string{"id", "author", "rating", "title", "description"}).AddRow(1, "example_mail.@example.com", 3, "review title", "review description")
				mock.ExpectQuery("SELECT id, author, rating, title, description FROM reviews").WithoutArgs().WillReturnRows(rows)
			},

			expectedLen: 1,
		},
		{
			name: "3 review",
			mockBehavior: func() {
				rows := mock.NewRows([]string{"id", "author", "rating", "title", "description"}).AddRow(1, "example_mail.@example.com", 3, "review title", "review description").AddRow(1, "example_mail.@example.com", 3, "review title", "review description").AddRow(1, "example_mail.@example.com", 3, "review title", "review description")
				mock.ExpectQuery("SELECT id, author, rating, title, description FROM reviews").WithoutArgs().WillReturnRows(rows)

			},

			expectedLen: 3,
		},
		{
			name: "0 review",
			mockBehavior: func() {
				rows := mock.NewRows([]string{"id", "author", "rating", "title", "description"})
				mock.ExpectQuery("SELECT id, author, rating, title, description FROM reviews").WithoutArgs().WillReturnRows(rows)
			},
			expectedLen: 0,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			testcase.mockBehavior()
			returnedReviews, err := store.Review().FindAll()
			assert.NoError(t, err)
			assert.Len(t, returnedReviews, testcase.expectedLen)
		})
	}
}
