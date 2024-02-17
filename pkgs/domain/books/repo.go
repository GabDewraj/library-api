package books

import "context"

type Repository interface {
	InsertBooks(ctx context.Context, newBooks []*Book) error
	GetBooks(ctx context.Context, params *GetBooksParams) ([]*Book, int, error)
	UpdateBook(ctx context.Context, arg *Book) error
}
