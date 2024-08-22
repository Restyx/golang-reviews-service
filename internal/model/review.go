package model

import (
	"github.com/go-playground/validator/v10"
	"github.com/leebenson/conform"
)

type Model interface {
	Validate() error
}

type Review struct {
	ID          int    `json:"id" validate:"omitempty"`
	Author      string `json:"author" validate:"required_without=ID,omitempty,email" conform:"trim"`
	Rating      int8   `json:"rating" validate:"required_without=ID,omitempty,gte=1,lte=10"`
	Title       string `json:"title" validate:"required_without=ID,omitempty,gte=3,lte=50" conform:"trim"`
	Description string `json:"description" validate:"required_without=ID,omitempty,gte=3,lte=500" conform:"trim"`
}

func (r *Review) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := conform.Strings(r); err != nil {
		return err
	}

	return validate.Struct(r)
}
