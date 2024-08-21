package model_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func generateRandomString(t *testing.T, length int) string {
	t.Helper()

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result)
}

func TestReview_Validate(t *testing.T) {
	testcases := []struct {
		name    string
		review  func() *model.Review
		isValid bool
	}{
		{
			name: "valid",
			review: func() *model.Review {
				return model.TestReview(t)
			},
			isValid: true,
		},
		{
			name: "empty author",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Author = ""
				return review
			},
			isValid: false,
		},
		{
			name: "invalid author",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Author = "invalid"
				return review
			},
			isValid: false,
		},
		{
			name: "0 rating",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Rating = 0
				return review
			},
			isValid: false,
		},
		{
			name: "0 rating with id",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.ID = 1
				review.Rating = 0
				return review
			},
			isValid: true,
		},
		{
			name: "invalid rating: too high",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Rating = 11
				return review
			},
			isValid: false,
		},
		{
			name: "empty title",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Title = ""
				return review
			},
			isValid: false,
		},
		{
			name: "empty title with id",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.ID = 1
				review.Title = ""
				return review
			},
			isValid: true,
		},
		{
			name: "whitespace title",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Title = "   "
				return review
			},
			isValid: false,
		},
		{
			name: "short title",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Title = generateRandomString(t, 2)
				return review
			},
			isValid: false,
		},
		{
			name: "long title",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Title = generateRandomString(t, 51)
				return review
			},
			isValid: false,
		},
		{
			name: "short description",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Description = generateRandomString(t, 2)
				return review
			},
			isValid: false,
		},
		{
			name: "long description",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Description = generateRandomString(t, 501)
				return review
			},
			isValid: false,
		},
		{
			name: "emty description",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Description = ""
				return review
			},
			isValid: false,
		},
		{
			name: "emty description with id",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.ID = 1
				review.Description = ""
				return review
			},
			isValid: true,
		},
		{
			name: "whitespace description",
			review: func() *model.Review {
				review := model.TestReview(t)
				review.Description = "   "
				return review
			},
			isValid: false,
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			if testcase.isValid {
				assert.NoError(t, testcase.review().Validate())
			} else {
				assert.Error(t, testcase.review().Validate())
			}
		})
	}

}
