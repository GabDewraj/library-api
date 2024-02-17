package books

import (
	"errors"
	"fmt"

	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/go-playground/validator"
)

type Availability string

const (
	Available    Availability = "available"
	NotAvailable Availability = "not_available"
)

type Book struct {
	ID           int              `json:"id" db:"id"`
	ISBN         string           `json:"isbn" db:"isbn" validate:"required"`
	Title        string           `json:"title" db:"title" validate:"required"`
	Author       string           `json:"author" db:"author" validate:"required"`
	Publisher    string           `json:"publisher" db:"publisher" validate:"required"`
	Published    utils.CustomDate `json:"published" db:"published" validate:"required"`
	Genre        string           `json:"genre" db:"genre" validate:"required"`
	Language     string           `json:"language" db:"language" validate:"required"`
	Pages        int              `json:"pages" db:"pages" validate:"required"`
	Availability Availability     `json:"available" db:"available" validate:"required"`
	UpdatedAt    utils.CustomTime `json:"updated_at" db:"updated_at"`
	CreatedAt    utils.CustomTime `json:"created_at" db:"created_at"`
	DeletedAt    utils.CustomTime `json:"deleted_at" db:"deleted_at"`
}

type GetBooksParams struct {
	ID           int
	Page         int
	PerPage      int
	UpdatedAt    utils.CustomTime
	ISBN         string
	Title        string
	Author       string
	Publisher    string
	Published    utils.CustomDate
	Genre        string
	Language     string
	Pages        int
	Availability Availability
}

// Object methods for aggregate root
// Validation for creating a Book
func (b *Book) ValidateCreateBook() error {
	validate := validator.New()

	err := validate.Struct(b)

	if err != nil {
		return validationErrMessage(err.(validator.ValidationErrors))
	}
	return nil
}

// Internal helper funcs for methods
func validationErrMessage(errs validator.ValidationErrors) error {
	for _, err := range errs {
		var errMessage string
		switch err.Tag() {
		case "required":
			errMessage = fmt.Sprintf("%s field is required", err.Field())
		case "min":
			errMessage = fmt.Sprintf("%s field is too short", err.Field())
		case "max":
			errMessage = fmt.Sprintf("%s field is too long", err.Field())
		default:
			errMessage = fmt.Sprintf("%s field has the following error %s", err.Field(), err.Tag())
		}
		return errors.New(errMessage)
	}
	return nil
}
