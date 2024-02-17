package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractQueryParams(t *testing.T) {
	assertWithTest := assert.New(t)
	// Create a request with query parameters
	req, err := http.NewRequest("GET", "/books??page=1&per_page=10&updated_at=2022-01-01T12:00:00Z&isbn=123456789&title=ExampleTitle&author=JohnDoe&publisher=ExamplePublisher&published=2022-01-01&genre=Fiction&language=English&pages=300&availability=Available", nil)
	assertWithTest.Nil(err)

	assertWithTest.NotNil(req)
}
