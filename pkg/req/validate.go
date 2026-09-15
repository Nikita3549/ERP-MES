package req

import "github.com/go-playground/validator/v10"

func IsValid[T any](dto *T) error {
	validate := validator.New()
	err := validate.Struct(dto)
	return err
}

