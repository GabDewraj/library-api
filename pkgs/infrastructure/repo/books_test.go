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
					ISBN:         "978-1234567890",
					Title:        "The Great Gatsby",
					Author:       "F. Scott Fitzgerald",
					Publisher:    "Scribner",
					Published:    utils.CustomDate{Time: time.Date(1990, 4, 10, 0, 0, 0, 0, time.UTC)},
					Genre:        "Fiction",
					Language:     "English",
					Pages:        180,
					Availability: books.Available,
					UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
					CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
					DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
				},
				{
					ISBN:         "978-0451524935",
					Title:        "1984",
					Author:       "George Orwell",
					Publisher:    "Signet Classic",
					Published:    utils.CustomDate{Time: time.Date(1980, 6, 8, 0, 0, 0, 0, time.UTC)},
					Genre:        "Dystopian",
					Language:     "English",
					Pages:        328,
					Availability: books.Available,
					UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
					CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
					DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
				},
				{
					ID:           3,
					ISBN:         "978-0061120084",
					Title:        "To Kill a Mockingbird",
					Author:       "Harper Lee",
					Publisher:    "Harper Perennial Modern Classics",
					Published:    utils.CustomDate{Time: time.Date(1980, 7, 11, 0, 0, 0, 0, time.UTC)},
					Genre:        "Classics",
					Language:     "English",
					Pages:        336,
					Availability: books.NotAvailable,
					UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
					CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
					DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
				},
				{
					ID:           4,
					ISBN:         "978-0142407332",
					Title:        "The Outsiders",
					Author:       "S.E. Hinton",
					Publisher:    "Penguin Books",
					Published:    utils.CustomDate{Time: time.Date(1980, 4, 24, 0, 0, 0, 0, time.UTC)},
					Genre:        "Young Adult",
					Language:     "English",
					Pages:        192,
					Availability: books.NotAvailable,
					UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
					CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
					DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
				},
				{
					ID:           5,
					ISBN:         "978-1400032493",
					Title:        "The Kite Runner",
					Author:       "Khaled Hosseini",
					Publisher:    "Riverhead Books",
					Published:    utils.CustomDate{Time: time.Date(2003, 5, 29, 0, 0, 0, 0, time.UTC)},
					Genre:        "Fiction",
					Language:     "English",
					Pages:        371,
					Availability: books.NotAvailable,
					UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
					CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
					DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
				},
			},
			Description: "Successful insert of books",
		},
	}
	for _, test := range testCases {
		err := booksRepo.InsertBooks(context.Background(), test.Input)
		assertWithTest.Equal(test.ExpectedErr, err, test.Description)

	}

}

func TestUpdateBook(t *testing.T) {
	assertWithTest := assert.New(t)
	booksRepo, err := testingBooksDB()
	assertWithTest.Nil(err, "Test org db conn successful")
	ctx := context.Background()
	seed := []*books.Book{
		{
			ISBN:         "978-1234567890",
			Title:        "The Great Gatsby",
			Author:       "F. Scott Fitzgerald",
			Publisher:    "Scribner",
			Published:    utils.CustomDate{Time: time.Date(1990, 4, 10, 0, 0, 0, 0, time.UTC)},
			Genre:        "Fiction",
			Language:     "English",
			Pages:        180,
			Availability: books.Available,
			UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
			CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
			DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
		},
		{
			ISBN:         "978-0451524935",
			Title:        "1984",
			Author:       "George Orwell",
			Publisher:    "Signet Classic",
			Published:    utils.CustomDate{Time: time.Date(1980, 6, 8, 0, 0, 0, 0, time.UTC)},
			Genre:        "Dystopian",
			Language:     "English",
			Pages:        328,
			Availability: books.Available,
			UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
			CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
			DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
		},
	}
	err = booksRepo.InsertBooks(ctx, seed)
	assertWithTest.Nil(err)
	testCases := []struct {
		ExpectedErr error
		Input       books.Book
		Description string
	}{
		{
			ExpectedErr: nil,
			Input: books.Book{
				ID:           seed[0].ID,
				ISBN:         "787877",
				Title:        "The Great Gatsby Updated",
				Author:       "F. Scott Fitzgerald Updated",
				Publisher:    "Scribner Updated",
				Published:    utils.CustomDate{Time: time.Date(1999, 4, 10, 0, 0, 0, 0, time.UTC)},
				Genre:        "Fiction Updated",
				Language:     "English Updated",
				Pages:        220,
				Availability: books.NotAvailable,
				DeletedAt:    utils.CustomTime{Time: time.Now()},
			},
			Description: "Update all fields",
		},
	}
	for _, test := range testCases {
		retrievedBooks, _, retrieveErr := booksRepo.GetBooks(ctx, &books.GetBooksParams{
			Title: "The Great Gatsby",
		})
		assertWithTest.NotNil(retrievedBooks, test.Description)
		assertWithTest.Equal(seed[0].ID, retrievedBooks[0].ID)
		assertWithTest.Nil(retrieveErr, test.Description)
		err := booksRepo.UpdateBook(ctx, &test.Input)
		assertWithTest.Nil(err)
	}

}

func TestGetBooks(t *testing.T) {
	assertWithTest := assert.New(t)
	booksRepo, err := testingBooksDB()
	assertWithTest.Nil(err, "Test org db conn successful")
	ctx := context.Background()
	seed := []*books.Book{
		{
			ISBN:         "978-1234567890",
			Title:        "The Great Gatsby",
			Author:       "F. Scott Fitzgerald",
			Publisher:    "Scribner",
			Published:    utils.CustomDate{Time: time.Date(1990, 4, 10, 0, 0, 0, 0, time.UTC)},
			Genre:        "Fiction",
			Language:     "English",
			Pages:        180,
			Availability: books.Available,
			UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
			CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
			DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
		},
		{
			ISBN:         "978-0451524935",
			Title:        "1984",
			Author:       "George Orwell",
			Publisher:    "Signet Classic",
			Published:    utils.CustomDate{Time: time.Date(1980, 6, 8, 0, 0, 0, 0, time.UTC)},
			Genre:        "Dystopian",
			Language:     "English",
			Pages:        328,
			Availability: books.Available,
			UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
			CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
			DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
		},
		{
			ID:           3,
			ISBN:         "978-0061120084",
			Title:        "To Kill a Mockingbird",
			Author:       "Harper Lee",
			Publisher:    "Harper Perennial Modern Classics",
			Published:    utils.CustomDate{Time: time.Date(1980, 7, 11, 0, 0, 0, 0, time.UTC)},
			Genre:        "Classics",
			Language:     "English",
			Pages:        336,
			Availability: books.NotAvailable,
			UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
			CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
			DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
		},
		{
			ID:           4,
			ISBN:         "978-0142407332",
			Title:        "The Outsiders",
			Author:       "S.E. Hinton",
			Publisher:    "Penguin Books",
			Published:    utils.CustomDate{Time: time.Date(1980, 4, 24, 0, 0, 0, 0, time.UTC)},
			Genre:        "Young Adult",
			Language:     "English",
			Pages:        192,
			Availability: books.NotAvailable,
			UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
			CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
			DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
		},
		{
			ID:           5,
			ISBN:         "978-1400032493",
			Title:        "The Kite Runner",
			Author:       "Khaled Hosseini",
			Publisher:    "Riverhead Books",
			Published:    utils.CustomDate{Time: time.Date(2003, 5, 29, 0, 0, 0, 0, time.UTC)},
			Genre:        "Fiction",
			Language:     "English",
			Pages:        371,
			Availability: books.NotAvailable,
			UpdatedAt:    utils.CustomTime{Time: time.Now().Add(-24 * time.Hour)},
			CreatedAt:    utils.CustomTime{Time: time.Now().Add(-48 * time.Hour)},
			DeletedAt:    utils.CustomTime{Time: time.Now().Add(-72 * time.Hour)},
		},
	}
	err = booksRepo.InsertBooks(ctx, seed)
	assertWithTest.Nil(err)
	testCases := []struct {
		ExpectedOutput struct {
			Count int
			Error error
		}
		Input       books.GetBooksParams
		Description string
	}{
		{
			ExpectedOutput: struct {
				Count int
				Error error
			}{
				Count: 5,
				Error: nil,
			},
			Input:       books.GetBooksParams{},
			Description: "Get all books successfully",
		},
		{
			ExpectedOutput: struct {
				Count int
				Error error
			}{
				Count: 2,
				Error: nil,
			},
			Input: books.GetBooksParams{
				Availability: books.Available,
				Page:         1,
				PerPage:      2,
			},
			Description: "Get all available books",
		},
		{
			ExpectedOutput: struct {
				Count int
				Error error
			}{
				Count: 3,
				Error: nil,
			},
			Input: books.GetBooksParams{
				Availability: books.NotAvailable,
			},
			Description: "Get all unavailable books",
		},
		{
			ExpectedOutput: struct {
				Count int
				Error error
			}{
				Count: 2,
				Error: nil,
			},
			Input: books.GetBooksParams{
				Availability: books.NotAvailable,
				Title:        "The",
			},
			Description: "Get all unavailable books that have `THE` in the title",
		},
		{
			ExpectedOutput: struct {
				Count int
				Error error
			}{
				Count: 1,
				Error: nil,
			},
			Input: books.GetBooksParams{
				ISBN:         seed[0].ISBN,
				Title:        seed[0].Title,
				Author:       seed[0].Author,
				Publisher:    seed[0].Publisher,
				Published:    seed[0].Published,
				Genre:        seed[0].Genre,
				Language:     seed[0].Language,
				Availability: seed[0].Availability,
				UpdatedAt: utils.CustomTime{
					Time: seed[0].UpdatedAt.Time.Add(-5 * time.Hour),
				},
			},
			Description: "make a specific search to find the first inserted book",
		},
		{
			ExpectedOutput: struct {
				Count int
				Error error
			}{
				Count: 1,
				Error: nil,
			},
			Input: books.GetBooksParams{
				ID: seed[0].ID,
			},
			Description: "Get book by ID",
		},
	}
	for _, test := range testCases {
		retrievedBooks, count, retrieveErr := booksRepo.GetBooks(ctx, &test.Input)
		assertWithTest.NotNil(retrievedBooks, test.Description)
		assertWithTest.Equal(test.ExpectedOutput.Count, count, test.Description)
		assertWithTest.Nil(retrieveErr, test.Description)
		if test.Description == "Get book by ID" {
			assertWithTest.Equal(seed[0].Title, retrievedBooks[0].Title)
		}
	}

}
