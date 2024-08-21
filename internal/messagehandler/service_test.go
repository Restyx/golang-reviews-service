package messagehandler_test

import (
	"testing"

	"github.com/Restyx/golang-reviews-service/internal/messagehandler"
	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store/testingstorage"
	"github.com/stretchr/testify/assert"
)

func TestMessageHandlerService_Create(t *testing.T) {
	service := messagehandler.NewService(testingstorage.New())

	mockBehaviourFunc := func(review *model.Review) error {
		return service.Create(review)
	}

	testTable := []struct {
		name           string
		inputReview    *model.Review
		mockBehaviour  func(*model.Review) error
		expectedReview *model.Review
		expectError    bool
	}{
		{
			name: "valid",
			inputReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "Review Title",
				Description: "Description of the review",
			},
			mockBehaviour: mockBehaviourFunc,
			expectedReview: &model.Review{
				ID:          1,
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
			mockBehaviour:  mockBehaviourFunc,
			expectedReview: nil,
			expectError:    true,
		},
		{
			name: "empty fields",
			inputReview: &model.Review{
				Author:      "example_mail@example.com",
				Rating:      3,
				Title:       "",
				Description: "",
			},
			mockBehaviour:  mockBehaviourFunc,
			expectedReview: nil,
			expectError:    true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			err := testcase.mockBehaviour(testcase.inputReview)

			if !testcase.expectError {
				assert.NoError(t, err)
				assert.EqualValues(t, testcase.inputReview, testcase.expectedReview)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestMessageHandlerService_Update(t *testing.T) {
	service := messagehandler.NewService(testingstorage.New())

	mockBehaviourFunc := func(id uint, updateReview *model.Review) error {
		updateReview.ID = id

		return service.Update(updateReview)
	}

	baseReview := &model.Review{
		Author:      "example_mail@example.com",
		Rating:      3,
		Title:       "Review Title",
		Description: "Description of the review",
	}

	testTable := []struct {
		name           string
		inputReview    *model.Review
		inputUpdate    *model.Review
		mockBehaviour  func(uint, *model.Review) error
		expectedReview *model.Review
		expectError    bool
	}{
		{
			name:        "valid",
			inputReview: baseReview,
			inputUpdate: &model.Review{
				Author:      "updated_mail@example.com",
				Rating:      4,
				Title:       "updated Review Title",
				Description: "updated Description of the review",
			},
			mockBehaviour: mockBehaviourFunc,
			expectedReview: &model.Review{
				ID:          1,
				Author:      "updated_mail@example.com",
				Rating:      4,
				Title:       "updated Review Title",
				Description: "updated Description of the review",
			},
		},
		{
			name:        "ivalid id",
			inputReview: baseReview,
			inputUpdate: &model.Review{
				ID:          0,
				Author:      "updated_mail@example.com",
				Rating:      4,
				Title:       "updated Review Title",
				Description: "updated Description of the review",
			},
			mockBehaviour: func(u uint, r *model.Review) error {
				return service.Update(r)
			},
			expectError: true,
		},
		{
			name:           "empty fields",
			inputReview:    baseReview,
			inputUpdate:    &model.Review{},
			mockBehaviour:  mockBehaviourFunc,
			expectedReview: baseReview,
			expectError:    false,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			service.Create(testcase.inputReview)
			err := testcase.mockBehaviour(testcase.inputReview.ID, testcase.inputUpdate)

			actualReview, _ := service.ReadOne(testcase.inputReview.ID)

			if !testcase.expectError {
				assert.NoError(t, err)
				assert.EqualValues(t, actualReview, testcase.inputUpdate)
				assert.EqualValues(t, actualReview, testcase.expectedReview)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestMessageHandlerService_Delete(t *testing.T) {
	service := messagehandler.NewService(testingstorage.New())

	baseReview := &model.Review{
		Author:      "example_mail@example.com",
		Rating:      3,
		Title:       "Review Title",
		Description: "Description of the review",
	}

	testTable := []struct {
		name          string
		mockBehaviour func(uint) error
		expectError   bool
	}{
		{
			name:          "valid",
			mockBehaviour: service.Delete,
		},
		{
			name: "invalid id",
			mockBehaviour: func(u uint) error {
				return service.Delete(0)
			},
			expectError: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			service.Create(baseReview)
			err := testcase.mockBehaviour(baseReview.ID)

			if !testcase.expectError {
				assert.NoError(t, err)

				actualReview, err := service.ReadOne(baseReview.ID)
				assert.Nil(t, actualReview)
				assert.Error(t, err)
			} else {
				assert.Error(t, err)
				actualReview, err := service.ReadOne(baseReview.ID)
				assert.NotNil(t, actualReview)
				assert.NoError(t, err)

			}
		})
	}
}

func TestMessageHandlerService_ReadOne(t *testing.T) {
	service := messagehandler.NewService(testingstorage.New())

	baseReview := &model.Review{
		Author:      "example_mail@example.com",
		Rating:      3,
		Title:       "Review Title",
		Description: "Description of the review",
	}

	testTable := []struct {
		name           string
		inputReview    *model.Review
		mockBehaviour  func(uint) (*model.Review, error)
		expectedReview *model.Review
		expectError    bool
	}{
		{
			name:           "valid",
			inputReview:    baseReview,
			mockBehaviour:  service.ReadOne,
			expectedReview: baseReview,
		},
		{
			name:        "invalid id",
			inputReview: baseReview,
			mockBehaviour: func(u uint) (*model.Review, error) {
				return service.ReadOne(0)
			},
			expectError: true,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			service.Create(baseReview)

			resultReview, err := testcase.mockBehaviour(baseReview.ID)

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

func TestMessageHandlerService_ReadAll(t *testing.T) {
	baseReview := &model.Review{
		Author:      "example_mail@example.com",
		Rating:      3,
		Title:       "Review Title",
		Description: "Description of the review",
	}

	mockBehaviourFunc := func(len int, service messagehandler.ServiceI) ([]model.Review, error) {
		for range len {
			service.Create(baseReview)
		}
		return service.ReadAll()
	}

	testTable := []struct {
		name          string
		inputLen      int
		mockBehaviour func(int, messagehandler.ServiceI) ([]model.Review, error)
		expectError   bool
	}{
		{
			name:          "valid",
			inputLen:      1,
			mockBehaviour: mockBehaviourFunc,
		},
		{
			name:          "valid",
			inputLen:      3,
			mockBehaviour: mockBehaviourFunc,
		},
	}

	for _, testcase := range testTable {
		t.Run(testcase.name, func(t *testing.T) {
			service := messagehandler.NewService(testingstorage.New())

			resultReviews, err := testcase.mockBehaviour(testcase.inputLen, service)

			if !testcase.expectError {
				assert.NoError(t, err)
				assert.Len(t, resultReviews, testcase.inputLen)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resultReviews)
			}
		})
	}
}
