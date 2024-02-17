package repo

import (
	"context"
	"testing"
	"time"

	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
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
		ExpectedErr error
		Input       []*books.Book
		Description string
	}{
		{
			ExpectedErr: nil,
			Input: []*books.Book{
				{
					ID:        1,
					ISBN:      "978-1234567890",
					Title:     "The Great Gatsby",
					Author:    "F. Scott Fitzgerald",
					Publisher: "Scribner",
					Published: utils.CustomDate{Time: time.Date(1925, 4, 10, 0, 0, 0, 0, time.UTC)},
					Genre:     "Fiction",
					Language:  "English",
					Pages:     180,
					Available: true,
					UpdatedAt: utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
					CreatedAt: utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
					DeletedAt: utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
				},
				// Add more books as needed...
				{
					ID:        2,
					ISBN:      "978-0451524935",
					Title:     "1984",
					Author:    "George Orwell",
					Publisher: "Signet Classic",
					Published: utils.CustomDate{Time: time.Date(1949, 6, 8, 0, 0, 0, 0, time.UTC)},
					Genre:     "Dystopian",
					Language:  "English",
					Pages:     328,
					Available: true,
					UpdatedAt: utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
					CreatedAt: utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
					DeletedAt: utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
				},
			},
			Description: "Successful insert of books",
		},
	}
	for _, test := range testCases {
		err := booksRepo.InsertBooks(context.Background(), test.Input)
		assertWithTest.Nil(err)
	}

}
