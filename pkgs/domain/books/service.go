package books

import "context"

type Service interface {
	CreateBooks(ctx context.Context, newBooks []*Book) error
	UpdateBook(ctx context.Context, updatedBook *Book) error
	GetBooks(ctx context.Context, params *GetBooksParams) ([]*Book, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateBooks implements Service.
func (s *service) CreateBooks(ctx context.Context, newBooks []*Book) error {
	// All entity agnostic business logic to do with creating a book goes here
	// This logic is not coupled directly with an entity but contains domain logic
	// that may be pertinant to work flow of a particular usecase such as a retrieval
	// of records from the db for calculations etc
	return s.repo.InsertBooks(ctx, newBooks)
}

// GetBooks implements Service.
func (s *service) GetBooks(ctx context.Context, params *GetBooksParams) ([]*Book, int, error) {
	// All entity agnostic business logic to do with getting books
	return s.repo.GetBooks(ctx, params)
}

// UpdateBook implements Service.
func (s *service) UpdateBook(ctx context.Context, updatedBook *Book) error {
	// All entity agnostic business logic to do with updating a book goes here
	return s.UpdateBook(ctx, updatedBook)
}
