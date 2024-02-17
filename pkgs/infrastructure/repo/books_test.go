package repo

import (
	"testing"

	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/stretchr/testify/assert"
)

func testingBooksDB() (booksRepo, error) {
	client, err := testConn()
	if err != nil {
		return booksRepo{}, err
	}
	return booksRepo{dbClient: client}, nil
}

func TestConn(t *testing.T) {
	assertWithTest := assert.New(t)
	// hello
	testingBooksDB()
	_, err := testConn()
	assertWithTest.Nil(err)
}

func TestCreateNewBook(t *testing.T) {
	assertWithTest := assert.New(t)
	booksRepo, err := testingBooksDB()
	assertWithTest.Nil(err, "Test org db conn successful")

	testCases := []struct {
		Error       error
		Input       books.Book
		Description string
	}{
		{
			Error:       nil,
			Input:       books.Book{},
			Description: "Successful insert of books",
		},
	}
	for _, test := range testCases {

	}

}
