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

	testTable := []struct {
		name           string
		inputReview    *model.Review
		create         func(*model.Review) error
		expectedReview *model.Review
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
			create: store.Review().Create,
			expectedReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
		},
		{
			name: "invalid email",
			inputReview: &model.Review{
				Author:      "invalid",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
			create:         store.Review().Create,
			expectedReview: nil,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {

			mock.ExpectExec("INSERT INTO reviews (author, rating, title, description)").WithArgs(testcase.inputReview.Author, testcase.inputReview.Rating, testcase.inputReview.Title, testcase.inputReview.Description)
			err := testcase.create(testcase.inputReview)

			if testcase.expectedReview != nil {
				testcase.expectedReview.ID = testcase.inputReview.ID
				assert.EqualValues(t, testcase.expectedReview, testcase.inputReview)
			} else {
				assert.Zero(t, testcase.inputReview.ID)
				assert.Error(t, err)
			}
		})
	}
}

func TestPostgresReviewRepositoryMock_FindOne(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

	baseReview := model.TestReview(t)
	if err := store.Review().Create(baseReview); err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name           string
		inputId        uint
		review         func(uint) (*model.Review, error)
		expectedReview *model.Review
	}{
		{
			name:           "valid",
			inputId:        baseReview.ID,
			review:         store.Review().FindOne,
			expectedReview: baseReview,
		},
		{
			name:           "not existing review",
			inputId:        0,
			review:         store.Review().FindOne,
			expectedReview: nil,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			returnedReview, err := testcase.review(testcase.inputId)

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

func TestPostgresReviewRepositoryMock_FindAll(t *testing.T) {
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

func TestPostgresReviewRepositoryMock_Update(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

	testUpdateFunc := func(r *model.Review) error {
		testingReview := &model.Review{
			Author:      "example_mail@example.com",
			Rating:      3,
			Title:       "review title",
			Description: "review description",
		}
		if err := store.Review().Create(testingReview); err != nil {
			t.Fatal(err)
		}

		r.ID = testingReview.ID

		return store.Review().Update(r)
	}

	testcases := []struct {
		name           string
		inputReview    *model.Review
		expectedReview *model.Review
		update         func(*model.Review) error
		expectError    bool
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
			update: testUpdateFunc,
			expectedReview: &model.Review{
				Author:      "updated_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
		},
		{
			name:        "invalid id",
			inputReview: &model.Review{},
			update: func(r *model.Review) error {
				return store.Review().Update(r)
			},
			expectError: true,
		},
		{
			name: "emty author",
			inputReview: &model.Review{
				Rating:      4,
				Title:       "review title",
				Description: "review description",
			},
			update: testUpdateFunc,
			expectedReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      4,
				Title:       "review title",
				Description: "review description",
			},
		},
		{
			name:        "emty fields",
			inputReview: &model.Review{},
			update:      testUpdateFunc,
			expectedReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
		},
		{
			name: "short fields",
			inputReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "as",
				Description: "ds",
			},
			update: testUpdateFunc,
			expectedReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "review title",
				Description: "review description",
			},
			expectError: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			err := testcase.update(testcase.inputReview)

			if !testcase.expectError {
				testcase.expectedReview.ID = testcase.inputReview.ID
				actualReview, _ := store.Review().FindOne(testcase.inputReview.ID)

				assert.NoError(t, err)
				assert.NotNil(t, actualReview)
				assert.EqualValues(t, testcase.expectedReview, actualReview)
			} else {
				assert.Error(t, err)
			}

		})
	}
}

func TestPostgresReviewRepositoryMock_Delete(t *testing.T) {
	database, teardown := postgres.TestPostgresDB(t, pgUser, pgPass, pgHost, pgPort, pgDB, pgSSL)
	defer teardown("reviews")

	store := postgres.New(database)

	testDeleteFunc := func(r *model.Review) error {
		return store.Review().Delete(r.ID)
	}

	testcases := []struct {
		name      string
		delete    func(*model.Review) error
		expectErr bool
	}{
		{
			name:      "valid",
			delete:    testDeleteFunc,
			expectErr: false,
		},
		{
			name: "invalid id",
			delete: func(r *model.Review) error {
				return store.Review().Delete(0)
			},
			expectErr: true,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			baseReview := model.TestReview(t)

			if err := store.Review().Create(baseReview); err != nil {
				t.Fatal(err)
			}
			err := testcase.delete(baseReview)

			if !testcase.expectErr {
				selectedReview, _ := store.Review().FindOne(baseReview.ID)
				assert.NoError(t, err)
				assert.Nil(t, selectedReview)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
