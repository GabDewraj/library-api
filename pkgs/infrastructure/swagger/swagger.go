package swagger

import "github.com/GabDewraj/library-api/pkgs/domain/books"

// File defining all Request bodies

type CreateBookRequestBody struct {
	ID           int
	ISBN         string
	Title        string
	Author       string
	Publisher    string
	Published    int64
	Genre        string
	Language     string
	Pages        int
	Availability string
}

type UpdateBookRequestBody struct {
	ID           int
	ISBN         string
	Title        string
	Author       string
	Publisher    string
	Published    int64
	Genre        string
	Language     string
	Pages        int
	Availability string
}

type GetBooksReponse struct {
	Books []*books.Book `json:"books"`
	Count int           `json:"count"`
}
