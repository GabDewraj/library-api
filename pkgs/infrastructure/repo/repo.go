package repo

// This is just a file to create generic helpers.
// These helpers are generic to a repo implementation, this avoids
// code repetition, but contains no domain logic or transaction logic.
// It crosses no domain boundaries, it provides generic database specific logic.
import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// This allows a test connection for unit testing the db layer in a clean way.
func testConn() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql",
		"root:password@tcp(localhost:3306)/library_dev")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err := cleanTestTB(db); err != nil {
		return nil, err
	}
	return db, nil
}

func cleanTestTB(db *sqlx.DB) error {
	if _, err := db.Exec("DELETE FROM books;"); err != nil {
		return fmt.Errorf("Could not delete books: %v", err)
	}
	return nil
}

func concludeTx(tx *sqlx.Tx, err error) error {
	// Defer panic handling for the function scope on commit or rollback
	defer handlePanic()

	// Commit the transaction if there was no error
	if err == nil {
		if commitErr := tx.Commit(); commitErr != nil {
			return fmt.Errorf("commit error: %v", commitErr)
		}
	} else {
		// Rollback the transaction in case of an error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("rollback error: %v", rollbackErr)
		}
	}

	return nil
}

func handleMysqlErr(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1062:
			return fmt.Errorf("entity already exists")
		default:
			return fmt.Errorf("Unexpected MySQL error: %v", err)
		}
	}
	return err
}

func handlePanic() {
	if r := recover(); r != nil {
		// Log or handle the panic as needed
		fmt.Printf("Panic occurred: %v\n", r)
	}
}
