package postgres_test

import (
	"testing"

	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPostgresReviewRepository_Create(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

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

func TestPostgresReviewRepository_FindOne(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

	baseReview := model.TestReview(t)

	id, err := store.Review().Create(baseReview)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name           string
		inputId        int
		expectedReview *model.Review
	}{
		{
			name:           "valid",
			inputId:        id,
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

func TestPostgresReviewRepository_FindAll(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

	baseReview := model.TestReview(t)

	testTable := []struct {
		name        string
		fillAndFind func() ([]model.Review, error)
		expectedLen int
	}{
		{
			name: "empty table",
			fillAndFind: func() ([]model.Review, error) {
				return store.Review().FindAll()
			},
			expectedLen: 0,
		},
		{
			name: "3 rows",
			fillAndFind: func() ([]model.Review, error) {
				store.Review().Create(baseReview)
				store.Review().Create(baseReview)
				store.Review().Create(baseReview)

				return store.Review().FindAll()
			},
			expectedLen: 3,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			reviews, err := testcase.fillAndFind()

			assert.NoError(t, err)
			assert.NotNil(t, reviews)
			assert.Len(t, reviews, testcase.expectedLen)
		})
	}
}

func TestPostgresReviewRepository_Update(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

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

	testcases := []struct {
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
			name:        "invalid id",
			inputReview: &model.Review{},
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

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			err := store.Review().Update(testcase.inputReview)

			if testcase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				actualReview, err := store.Review().FindOne(id)
				if err != nil {
					t.Fatal(err)
				}
				assert.EqualValues(t, testcase.expectedReview, actualReview)
			}

		})
	}
}

func TestPostgresReviewRepository_Delete(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

	baseReview := model.TestReview(t)
	id, err := store.Review().Create(baseReview)
	if err != nil {
		t.Fatal(err)
	}

	testcases := []struct {
		name        string
		inputId     int
		expectError bool
	}{
		{
			name:    "valid",
			inputId: id,
		},
		{
			name:        "not existing record",
			inputId:     id,
			expectError: true,
		},
	}

	for _, testcase := range testcases {
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
