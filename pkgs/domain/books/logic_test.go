package books

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	assertWithTest := assert.New(t)
	var newTest Book
	v := reflect.ValueOf(newTest)
	typeOfT := v.Type()

	var columns []string

	for i := 0; i < v.NumField(); i++ {
		field := typeOfT.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columns = append(columns, dbTag)
		}
	}

	assertWithTest.Nil(columns)
}
