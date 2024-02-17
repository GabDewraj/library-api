package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/cache"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/go-chi/chi"
	"go.uber.org/fx"
)

// Uber fx: Package by uber used for dependency management in app and server lifecycle
type BooksHandlerParams struct {
	fx.In
	BookService books.Service
	// Other infrastructure layer services: Generic packages that don't contain business logic
	CacheService cache.Service
}

type booksHandler struct {
	bookService  books.Service
	cacheService cache.Service
}

type BooksHandler interface {
	CreateBook(res http.ResponseWriter, req *http.Request)
	UpdateBook(res http.ResponseWriter, req *http.Request)
	GetBooks(res http.ResponseWriter, req *http.Request)
	GetBookByID(res http.ResponseWriter, req *http.Request)
	ArchiveBook(res http.ResponseWriter, req *http.Request)
}

func NewBooksHandler(p BooksHandlerParams) BooksHandler {
	return &booksHandler{
		bookService:  p.BookService,
		cacheService: p.CacheService,
	}
}

// CreateBook implements BooksHandler.
func (h *booksHandler) CreateBook(res http.ResponseWriter, req *http.Request) {
	var newBook books.Book
	if err := json.NewDecoder(req.Body).Decode(&newBook); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// Validate the Request
	if err := newBook.ValidateCreateBook(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Serve domain data to context domain service function
	if err := h.bookService.CreateBooks(req.Context(), []*books.Book{&newBook}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	payload, err := json.Marshal(newBook)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := res.Write(payload); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetBookByID implements BooksHandler.
func (h *booksHandler) GetBookByID(res http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParamFromCtx(req.Context(), "book_id")
	// Scope the input to a urlParam
	bookID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	retrievedBook, _, err := h.bookService.GetBooks(req.Context(), &books.GetBooksParams{ID: bookID})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	payload, err := json.Marshal(retrievedBook)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := res.Write(payload); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetBooks implements BooksHandler.
func (h *booksHandler) GetBooks(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	var params books.GetBooksParams
	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		params.Page = page
	}
	if perPageStr := query.Get("per_page"); perPageStr != "" {
		perPage, err := strconv.Atoi(perPageStr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		params.PerPage = perPage
	}
	if updatedAtStr := query.Get("updated_at"); updatedAtStr != "" {
		convertedUpdatedAt, err := strconv.Atoi(updatedAtStr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		params.UpdatedAt = utils.CustomTime{
			Time: time.Unix(int64(convertedUpdatedAt), 0),
		}
	}
	if bookPagesStr := query.Get("book_pages"); bookPagesStr != "" {
		pages, err := strconv.Atoi(bookPagesStr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		params.BookPages = pages
	}
	if publishedStr := query.Get("published"); publishedStr != "" {
		published, err := strconv.Atoi(publishedStr)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		params.Published = utils.CustomDate{
			Time: time.Unix(int64(published), 0),
		}
	}
	// String values
	params.ISBN = query.Get("isbn")
	params.Title = query.Get("title")
	params.Author = query.Get("author")
	params.Publisher = query.Get("publisher")
	params.Genre = query.Get("genre")
	params.Language = query.Get("language")
	params.Availability = books.Availability(query.Get("availability"))
	retrievedBooks, count, err := h.bookService.GetBooks(req.Context(), &params)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		Books []*books.Book `json:"books"`
		Count int           `json:"count"`
	}{
		Books: retrievedBooks,
		Count: count,
	}
	payload, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := res.Write(payload); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

}

// UpdateBook implements BooksHandler.
func (h *booksHandler) UpdateBook(res http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParamFromCtx(req.Context(), "book_id")
	// Scope the input to a urlParam
	bookID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	requestBody := struct {
		ISBN         string             `json:"isbn"`
		Title        string             `json:"title"`
		Author       string             `json:"author"`
		Publisher    string             `json:"publisher"`
		Published    utils.CustomDate   `json:"published"`
		Genre        string             `json:"genre"`
		Language     string             `json:"language"`
		Pages        int                `json:"pages"`
		Availability books.Availability `json:"availability"`
	}{}
	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	updatedBook := books.Book{
		ID:           bookID,
		ISBN:         requestBody.ISBN,
		Title:        requestBody.Title,
		Author:       requestBody.Author,
		Publisher:    requestBody.Publisher,
		Published:    requestBody.Published,
		Genre:        requestBody.Genre,
		Language:     requestBody.Language,
		Pages:        requestBody.Pages,
		Availability: requestBody.Availability,
	}
	if err = h.bookService.UpdateBook(req.Context(), &updatedBook); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	payload, err := json.Marshal(updatedBook)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := res.Write(payload); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

}

// ArchiveBook implements BooksHandler.
func (h *booksHandler) ArchiveBook(res http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParamFromCtx(req.Context(), "book_id")
	// Scope the input to a urlParam
	bookID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedBook := books.Book{
		ID:        bookID,
		DeletedAt: utils.CustomTime{Time: time.Now()},
	}
	// Function is extensible to soft delete by updating the book deleted_at field
	if err = h.bookService.UpdateBook(req.Context(), &updatedBook); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(200)
}
