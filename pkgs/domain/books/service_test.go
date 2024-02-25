package books

import (
	"errors"
	"testing"
	"time"

	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookValidation(t *testing.T) {
	assertWithTest := assert.New(t)
	testCases := []struct {
		Input         Book
		ExpectedError error
		Message       string
	}{
		{
			Input: Book{
				ISBN:         "978-0451524935",
				Title:        "1984",
				Author:       "George Orwell",
				Publisher:    "Signet Classic",
				Published:    utils.CustomDate{Time: time.Date(1980, 6, 8, 0, 0, 0, 0, time.UTC)},
				Genre:        "Dystopian",
				Language:     "English",
				Pages:        328,
				Availability: Available,
				UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
				CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
				DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
			},
			ExpectedError: nil,
			Message:       "Correct format for Book",
		},
		{
			Input: Book{
				ISBN:         "978-0451524935",
				Title:        "1984",
				Author:       "George Orwell",
				Publisher:    "Signet Classic",
				Published:    utils.CustomDate{Time: time.Date(1980, 6, 8, 0, 0, 0, 0, time.UTC)},
				Genre:        "Dystopian",
				Pages:        328,
				Availability: Available,
				UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
				CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
				DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
			},
			ExpectedError: errors.New("language field is required"),
			Message:       "Correct format for Book",
		},
	}
	for _, test := range testCases {
		err := test.Input.ValidateCreateBook()
		assertWithTest.Equal(test.ExpectedError, err)
	}
}
