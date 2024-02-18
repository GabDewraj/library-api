package repo

import (
	"context"
	"time"

	"github.com/GabDewraj/library-api/pkgs/domain/books"
	"github.com/GabDewraj/library-api/pkgs/infrastructure/utils"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type booksRepo struct {
	dbClient *sqlx.DB
}

// InsertBooks implements books.Repository.

func NewBooksDB(db *sqlx.DB) books.Repository {
	return &booksRepo{
		dbClient: db,
	}
}
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
func (repo *booksRepo) UpdateBook(ctx context.Context, arg *books.Book) error {
	// Start transaction
	tx, err := repo.dbClient.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer concludeTx(tx, &err)
	if err := repo.updatebook(ctx, tx, arg); err != nil {
		return err
	}
	return nil
}

func (repo *booksRepo) DeleteBookByID(ctx context.Context, id int) error {
	// Start transaction
	tx, err := repo.dbClient.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer concludeTx(tx, &err)
	if err := repo.deleteBookByID(ctx, tx, id); err != nil {
		return err
	}
	return nil
}

// GetBooks implements books.Repository.
func (b *booksRepo) GetBooks(ctx context.Context, params *books.GetBooksParams) ([]*books.Book, int, error) {
	// Enter nil for sqlx.ExtContext as this query does not form part of a transaction chain
	return b.getBooks(ctx, nil, params)
}

func (p *booksRepo) updatebook(ctx context.Context, ext sqlx.ExtContext, updatedBook *books.Book) error {
	updateBuilder := squirrel.Update("books")

	if updatedBook.ISBN != "" {
		updateBuilder = updateBuilder.Set("isbn", updatedBook.ISBN)
	}
	if updatedBook.Title != "" {
		updateBuilder = updateBuilder.Set("title", updatedBook.Title)
	}
	if updatedBook.Author != "" {
		updateBuilder = updateBuilder.Set("author", updatedBook.Author)
	}
	if updatedBook.Publisher != "" {
		updateBuilder = updateBuilder.Set("publisher", updatedBook.Publisher)
	}
	if (updatedBook.Published != utils.CustomDate{}) {
		updateBuilder = updateBuilder.Set("published", updatedBook.Published.Time)
	}
	if updatedBook.Genre != "" {
		updateBuilder = updateBuilder.Set("genre", updatedBook.Genre)
	}
	if updatedBook.Language != "" {
		updateBuilder = updateBuilder.Set("language", updatedBook.Language)
	}
	if updatedBook.Pages != 0 {
		updateBuilder = updateBuilder.Set("pages", updatedBook.Pages)
	}
	if updatedBook.Availability != "" {
		updateBuilder = updateBuilder.Set("availability", updatedBook.Availability)
	}
	if (updatedBook.DeletedAt != utils.CustomTime{}) {
		updateBuilder = updateBuilder.Set("deleted_at", updatedBook.DeletedAt.Time)
	}
	// Always update the updated at field
	updateBuilder = updateBuilder.Set("updated_at", utils.CustomTime{Time: time.Now()}.Time)
	updateBuilder = updateBuilder.Where(squirrel.Eq{"id": updatedBook.ID})
	// Build the final SQL query and arguments
	sql, args, err := updateBuilder.ToSql()
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

func (p *booksRepo) insertBooks(ctx context.Context, ext sqlx.ExtContext, books []*books.Book) error {
	// Make an efficient insert using a sql statement builder
	ib := squirrel.Insert("books").Columns(
		"isbn", "title", "author", "publisher", "published",
		"genre", "language", "pages", "availability", "updated_at", "created_at",
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
			book.Language, book.Pages, book.Availability, book.UpdatedAt.Time, book.CreatedAt.Time,
		)
	}
	// Build the final SQL query and arguments
	sql, args, err := ib.ToSql()
	if err != nil {
		return err
	}
	// Execute the query with ExecContext
	result, err := ext.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	// Retrieve last insert ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// Write DB primary key ID back to the pointer
	for _, book := range books {
		book.ID = int(lastInsertID)
		lastInsertID++
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
		"genre", "language", "pages", "availability", "updated_at", "created_at").From("books")
	sb = sb.Where("deleted_at IS NULL")
	// Select by id
	if params.ID != 0 {
		sb = sb.Where(squirrel.Eq{"id": params.ID})
	}
	// Availability is a binary value that always holds a statement significant to the business context
	if params.Availability != "" {
		sb = sb.Where(squirrel.Eq{"availability": params.Availability})
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

// hard delete
func (repo *booksRepo) deleteBookByID(ctx context.Context, ext sqlx.ExtContext, id int) error {
	if _, err := repo.dbClient.Exec("DELETE FROM books where id=?;", id); err != nil {
		return err
	}
	return nil
}
