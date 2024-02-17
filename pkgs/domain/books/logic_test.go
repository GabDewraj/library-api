package books

import (
	"reflect"
	"testing"
	"time"

	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	assertWithTest := assert.New(t)
	newTest := Book{
		ISBN:         "978-1234567890",
		Title:        "The Great Gatsby",
		Author:       "F. Scott Fitzgerald",
		Publisher:    "Scribner",
		Published:    utils.CustomDate{Time: time.Date(1990, 4, 10, 0, 0, 0, 0, time.UTC)},
		Genre:        "Fiction",
		Language:     "English",
		Pages:        180,
		Availability: Available,
		UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
		CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
		DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
	}
	v := reflect.ValueOf(newTest)
	typeOfT := v.Type()

	var columns []string
	var values []interface{}

	for i := 0; i < v.NumField(); i++ {
		field := typeOfT.Field(i)
		dbTag := field.Tag.Get("db")

		if dbTag != "" {
			columns = append(columns, dbTag)
			values = append(values)
		}
	}

	assertWithTest.Nil(columns)
	assertWithTest.Nil(values)
}
