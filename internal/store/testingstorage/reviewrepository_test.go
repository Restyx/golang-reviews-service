package testingstorage_test

import (
	"testing"

	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store/testingstorage"
	"github.com/stretchr/testify/assert"
)

func TestReviewRepository_Create(t *testing.T) {
	store := testingstorage.New()

	testTable := []struct {
		name        string
		inputReview *model.Review
		create      func(*model.Review) (int, error)
		expectError bool
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
			create:      store.Review().Create,
			expectError: false,
		},
		{
			name: "invalid email",
			inputReview: &model.Review{
				Author:      "invalid",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
			create:      store.Review().Create,
			expectError: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			id, err := testcase.create(testcase.inputReview)

			if testcase.expectError {
				assert.Error(t, err)
				assert.Zero(t, id)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, id)
			}
		})
	}
}

func TestReviewRepository_FindOne(t *testing.T) {
	store := testingstorage.New()

	baseReview := model.TestReview(t)

	id, err := store.Review().Create(baseReview)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name           string
		inputId        int
		review         func(int) (*model.Review, error)
		expectedReview *model.Review
	}{
		{
			name:           "valid",
			inputId:        int(id),
			expectedReview: baseReview,
		},
		{
			name:           "not existing review",
			inputId:        0,
			expectedReview: nil,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
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
	store := testingstorage.New()

	testTable := []struct {
		name     string
		find     func() ([]model.Review, error)
		expected int
	}{
		{
			name: "empty table",
			find: func() ([]model.Review, error) {
				return store.Review().FindAll()
			},
			expected: 0,
		},
		{
			name: "3 rows",
			find: func() ([]model.Review, error) {
				testingReview := model.TestReview(t)

				store.Review().Create(testingReview)
				store.Review().Create(testingReview)
				store.Review().Create(testingReview)

				return store.Review().FindAll()
			},
			expected: 3,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			reviews, err := testcase.find()

			assert.NoError(t, err)
			assert.NotNil(t, reviews)
			assert.Len(t, reviews, testcase.expected)
		})
	}
}

func TestReviewRepository_Update(t *testing.T) {
	store := testingstorage.New()

	baseReview := &model.Review{
		Author:      "example_mail@example.com",
		Rating:      3,
		Title:       "review title",
		Description: "review description",
	}

	id, err := store.Review().Create(baseReview)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name           string
		inputReview    *model.Review
		expectedReview *model.Review
		expectError    bool
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				ID:          id,
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
			expectedReview: &model.Review{
				ID:          id,
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
		},
		{
			name: "invalid id",
			inputReview: &model.Review{
				ID:          0,
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
			expectError: true,
		},
		{
			name: "emty author",
			inputReview: &model.Review{
				ID:          id,
				Rating:      4,
				Title:       "review title",
				Description: "review description",
			},
			expectedReview: &model.Review{
				ID:          id,
				Author:      "updated_mail@example.com",
				Rating:      4,
				Title:       "review title",
				Description: "review description",
			},
		},
		{
			name: "emty fields",
			inputReview: &model.Review{
				ID: id,
			},
			expectedReview: &model.Review{
				ID:          id,
				Author:      "updated_mail@example.com",
				Rating:      4,
				Title:       "review title",
				Description: "review description",
			},
		},
		{
			name: "short fields",
			inputReview: &model.Review{
				ID:          id,
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "as",
				Description: "ds",
			},
			expectError: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			err := store.Review().Update(testcase.inputReview)

			if testcase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				actualReview, err := store.Review().FindOne(int(id))
				if err != nil {
					t.Fatal(err)
				}
				assert.EqualValues(t, testcase.expectedReview, actualReview)
			}
		})
	}
}

func TestReviewRepository_Delete(t *testing.T) {
	store := testingstorage.New()

	baseReview := model.TestReview(t)
	id, err := store.Review().Create(baseReview)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name        string
		inputId     int
		expectError bool
	}{
		{
			name:        "valid",
			inputId:     id,
			expectError: false,
		},
		{
			name:        "not existing record",
			inputId:     id,
			expectError: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			err := store.Review().Delete(testcase.inputId)

			if testcase.expectError {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)

				_, err = store.Review().FindOne(testcase.inputId)
				assert.Error(t, err)
			}
		})
	}
}
