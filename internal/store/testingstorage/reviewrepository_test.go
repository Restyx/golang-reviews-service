package testingstorage_test

import (
	"testing"

	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store/testingstorage"
	"github.com/stretchr/testify/assert"
)

func TestReviewRepositoryCreate(t *testing.T) {
	store := testingstorage.New()

	testcases := []struct {
		name        string
		review      func(*model.Review) error
		expectError bool
	}{
		{
			name: "valid",
			review: func(review *model.Review) error {
				return store.Review().Create(review)
			},
		},
		{
			name: "invalid email",
			review: func(review *model.Review) error {
				review.Author = "invalid"

				return store.Review().Create(review)
			},
			expectError: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testingReview := model.TestReview(t)

			err := testcase.review(testingReview)

			if !testcase.expectError {
				assert.NoError(t, err)
				assert.NotZero(t, testingReview.ID)
			} else {
				assert.Error(t, err)
				assert.Zero(t, testingReview.ID)
			}
		})
	}
}

func TestReviewRepositoryFindOne(t *testing.T) {
	store := testingstorage.New()

	baseReview := &model.Review{
		Author:      "example_mail@example.com",
		Rating:      3,
		Title:       "Review Title",
		Description: "Description of the review",
	}

	testTable := []struct {
		name           string
		inputReview    *model.Review
		find           func(uint) (*model.Review, error)
		expectedReview *model.Review
		expectError    bool
	}{
		{
			name:           "valid",
			inputReview:    baseReview,
			find:           store.Review().FindOne,
			expectedReview: baseReview,
		},
		{
			name:        "invalid id",
			inputReview: baseReview,
			find: func(u uint) (*model.Review, error) {
				return store.Review().FindOne(0)
			},
			expectError: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			store.Review().Create(baseReview)

			resultReview, err := testcase.find(baseReview.ID)

			if !testcase.expectError {
				assert.NoError(t, err)
				assert.EqualValues(t, baseReview, resultReview)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resultReview)
			}
		})
	}
}

func TestReviewRepositoryFindAll(t *testing.T) {
	store := testingstorage.New()

	testcases := []struct {
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

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			reviews, err := testcase.find()

			assert.NoError(t, err)
			assert.NotNil(t, reviews)
			assert.Len(t, reviews, testcase.expected)
		})
	}
}

func TestReviewRepositoryUpdate(t *testing.T) {
	store := testingstorage.New()

	testcases := []struct {
		name     string
		update   func(*model.Review) error
		expected bool
	}{
		{
			name: "valid",
			update: func(createdReview *model.Review) error {
				createdReview.Author = "updated_mail@example.com"

				return store.Review().Update(createdReview)
			},
			expected: true,
		},
		{
			name: "invalid id",
			update: func(createdReview *model.Review) error {
				createdReview.ID = 0

				return store.Review().Update(createdReview)
			},
			expected: false,
		},
		{
			name: "emty fields",
			update: func(createdReview *model.Review) error {
				createdReview.Rating = 2
				createdReview.Title = ""
				createdReview.Description = ""

				return store.Review().Update(createdReview)
			},
			expected: true,
		},
		{
			name: "short fields",
			update: func(createdReview *model.Review) error {
				createdReview.Rating = 2
				createdReview.Title = "as"
				createdReview.Description = "ds"

				return store.Review().Update(createdReview)
			},
			expected: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testingReview := model.TestReview(t)
			if err := store.Review().Create(testingReview); err != nil {
				t.Fatal(err)
			}
			err := testcase.update(testingReview)

			if testcase.expected {
				selectedReview, _ := store.Review().FindOne(testingReview.ID)

				assert.NoError(t, err)
				assert.NotNil(t, testingReview)
				assert.EqualValues(t, testingReview, selectedReview)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestReviewRepositoryDelete(t *testing.T) {
	store := testingstorage.New()

	testcases := []struct {
		name     string
		delete   func(*model.Review) error
		expected bool
	}{
		{
			name: "valid",
			delete: func(r *model.Review) error {
				return store.Review().Delete(r.ID)
			},
			expected: true,
		},
		{
			name: "invalid id",
			delete: func(r *model.Review) error {
				return store.Review().Delete(0)
			},
			expected: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			testingReview := model.TestReview(t)

			if err := store.Review().Create(testingReview); err != nil {
				t.Fatal(err)
			}
			err := testcase.delete(testingReview)

			if testcase.expected {
				selectedReview, _ := store.Review().FindOne(testingReview.ID)
				assert.NoError(t, err)
				assert.Nil(t, selectedReview)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
