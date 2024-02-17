package repo

import (
	"context"
	"time"

	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type booksRepo struct {
	dbClient *sqlx.DB
}

// ArchiveBook implements books.Repository.
func (*booksRepo) ArchiveBook(ctx context.Context, id int) error {
	panic("unimplemented")
}

// DeleteBook implements books.Repository.
func (*booksRepo) DeleteBook(ctx context.Context, id int) error {
	panic("unimplemented")
}

// InsertBooks implements books.Repository.
func (repo *booksRepo) InsertBooks(ctx context.Context, newBooks []*books.Book) error {
	// Start transaction
	tx, err := repo.dbClient.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer concludeTx(tx, &err)
	if err := repo.insertBooks(ctx, tx, newBooks); err != nil {
		return err
	}
	return nil
}

// UpdateBook implements books.Repository.
func (*booksRepo) UpdateBook(ctx context.Context, arg *books.Book) error {
	panic("unimplemented")
}

func NewlibraryDB(db *sqlx.DB) books.Repository {
	return &booksRepo{
		dbClient: db,
	}
}

// GetBooks implements books.Repository.
func (b *booksRepo) GetBooks(ctx context.Context, params books.GetBooksParams) ([]*books.Book, int, error) {
	// Enter nil for sqlx.ExtContext as this query does not form part of a transaction chain
	return b.getBooks(ctx, nil, &params)
}

func (p *booksRepo) insertBooks(ctx context.Context, ext sqlx.ExtContext, books []*books.Book) error {
	// Make an efficient insert using a sql statement builder
	ib := squirrel.Insert("books").Columns(
		"isbn", "title", "author", "publisher", "published",
		"genre", "language", "pages", "available", "updated_at", "created_at",
	)

	// Add values for each user
	for _, book := range books {
		book.CreatedAt = utils.CustomTime{
			Time: time.Now(),
		}
		book.UpdatedAt = utils.CustomTime{
			Time: time.Now(),
		}
		ib = ib.Values(
			book.ISBN, book.Title, book.Author, book.Publisher, book.Published.Time, book.Genre,
			book.Language, book.Pages, book.Available, book.UpdatedAt.Time, book.CreatedAt.Time,
		)
	}
	// Build the final SQL query and arguments
	sql, args, err := ib.ToSql()
	if err != nil {
		return err
	}
	// Execute the query with ExecContext
	_, err = ext.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

// Get requests can be included in a transaction or independent of one
// Allow flexibility to include in a transaction or not
// Return map to allow for linear allocation of child objects to parents
// Return slice of project ids to allow for project scop child selection
func (repo *booksRepo) getBooks(ctx context.Context, ext sqlx.ExtContext,
	params *books.GetBooksParams) ([]*books.Book, int, error) {
	var userBooks []*books.Book
	sb := squirrel.Select("id", "isbn", "title", "author", "publisher", "published",
		"genre", "language", "pages", "available", "updated_at", "created_at").From("books")
	sb = sb.Where("deleted_at IS NULL")
	// Availability is a binary value that always holds a statement significant to the business context
	if params.Available != "" {
		sb = sb.Where(squirrel.Eq{"available": params.Available})
	}
	if params.ISBN != "" {
		sb = sb.Where(squirrel.Eq{"isbn": params.ISBN})
	}
	if params.Title != "" {
		sb = sb.Where(squirrel.Like{"title": "%" + params.Title + "%"})
	}
	if params.Author != "" {
		sb = sb.Where(squirrel.Like{"author": "%" + params.Author + "%"})
	}
	if params.Publisher != "" {
		sb = sb.Where(squirrel.Like{"publisher": "%" + params.Publisher + "%"})
	}
	if params.Genre != "" {
		sb = sb.Where(squirrel.Like{"genre": "%" + params.Genre + "%"})
	}
	if params.Language != "" {
		sb = sb.Where(squirrel.Eq{"language": params.Language})
	}
	if (params.Published != utils.CustomDate{}) {
		sb = sb.Where(squirrel.Eq{"published": params.Published.Time})
	}
	if (params.UpdatedAt != utils.CustomTime{}) {
		sb = sb.Where("updated_at >= ?", params.UpdatedAt.Time)
	}
	// We always want to order the retrieved data by the updated_at
	sb = sb.OrderBy("updated_at")
	// If we choose a specific page of results
	if params.Page > 0 {
		offset := (params.Page - 1) * params.PerPage
		sb = sb.Offset(uint64(offset))
	}
	// We always limit number of results retrieved
	// Default can be set in domain or handler
	// If we choose a specific page of results
	if params.PerPage > 0 {
		sb = sb.Limit(uint64(params.PerPage))
	}
	query, args, err := sb.ToSql()
	if err != nil {
		return nil, -1, err
	}
	logrus.Infoln(query, args)

	var bookRows *sqlx.Rows
	switch ext {
	case nil:
		queryRows, err := repo.dbClient.QueryxContext(ctx, query, args...)
		if err != nil {
			return nil, -1, err
		}
		bookRows = queryRows
	default:
		txRows, err := ext.QueryxContext(ctx, query, args...)
		if err != nil {
			return nil, -1, err
		}
		bookRows = txRows
	}
	defer bookRows.Close()
	rowsErr := bookRows.Err()
	if rowsErr != nil {
		logrus.Error(rowsErr)
		return nil, -1, rowsErr
	}
	for bookRows.Next() {
		var book books.Book
		if err := bookRows.StructScan(&book); err != nil {
			return nil, -1, err
		}
		userBooks = append(userBooks, &book)
	}
	return userBooks, len(userBooks), nil
}
