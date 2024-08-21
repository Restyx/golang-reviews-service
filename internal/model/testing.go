package model

import "testing"

func TestReview(t *testing.T) *Review {
	t.Helper()

	return &Review{
		Author:      "example_mail@example.com",
		Rating:      3,
		Title:       "Title",
		Description: "Description of the review",
	}
}
