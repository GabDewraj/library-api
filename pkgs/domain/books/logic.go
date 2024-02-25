package books

import (
	"errors"
	"fmt"
	"strings"

	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/go-playground/validator"
)

type Availability string

// Create global errors that are specific to this domain
var (
	ErrBookAlreadyExists = errors.New("book already exists")
)

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
	Availability Availability     `json:"availability" db:"availability" validate:"required,eq=available|eq=not_available"`
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
	BookPages    int
	Availability Availability
}

// Object methods for aggregate root
// Validation for creating a Book
func (b *Book) ValidateCreateBook() error {
	validate := validator.New()

	err := validate.Struct(b)

	if err != nil {
		err, _ := validationErrMessage(err.(validator.ValidationErrors))
		return err
	}
	return nil
}
func (b *Book) ValidateUpdateBook() error {
	validate := validator.New()
	err := validate.Struct(b)
	if err != nil {
		err, tag := validationErrMessage(err.(validator.ValidationErrors))
		if tag == "eq=available|eq=not_available" {
			return err
		}
		return nil
	}
	return nil
}

// Internal helper funcs for methods
func validationErrMessage(errs validator.ValidationErrors) (error, string) {
	for _, err := range errs {
		var errMessage string
		switch err.Tag() {
		case "required":
			errMessage = fmt.Sprintf("%s field is required", strings.ToLower(err.Field()))
		case "min":
			errMessage = fmt.Sprintf("%s field is too short", strings.ToLower(err.Field()))
		case "max":
			errMessage = fmt.Sprintf("%s field is too long", strings.ToLower(err.Field()))
		case "eq=available|eq=not_available":
			errMessage = fmt.Sprintf("value for %s is not recognised, please use available or not_available", strings.ToLower(err.Field()))
		default:
			errMessage = fmt.Sprintf("value for %s is not recognized", strings.ToLower(err.Field()))
		}
		return errors.New(errMessage), err.Tag()
	}
	return nil, ""
}
