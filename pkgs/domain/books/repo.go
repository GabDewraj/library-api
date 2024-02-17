package books

import "context"

type Repository interface {
	InsertBooks(ctx context.Context, newBooks []*Book) error
	GetBooks(ctx context.Context, params GetBooksParams) ([]*Book, error)
	UpdateBook(ctx context.Context, arg *Book) error
	ArchiveBook(ctx context.Context, id int) error
}
