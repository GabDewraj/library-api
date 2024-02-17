package books

import "github.com/GabDewraj/library-api/pkgs/infrastructure/utils"

type Book struct {
	ID        int              `json:"id" db:"id"`
	ISBN      string           `json:"isbn" db:"isbn"`
	Title     string           `json:"title" db:"title"`
	Author    string           `json:"author" db:"author"`
	Publisher string           `json:"publisher" db:"publisher"`
	Published utils.CustomDate `json:"published" db:"published"`
	Genre     string           `json:"genre" db:"genre"`
	Language  string           `json:"language" db:"language"`
	Pages     int              `json:"pages" db:"pages"`
	Available bool             `json:"available" db:"available"`
	UpdatedAt utils.CustomTime `json:"updated_at" db:"updated_at"`
	CreatedAt utils.CustomTime `json:"created_at" db:"created_at"`
	DeletedAt utils.CustomTime `json:"deleted_at" db:"deleted_at"`
}

type GetBooksParams struct {
	Page      int
	PerPage   int
	UpdatedAt utils.CustomTime
	ISBN      string
	Title     string
	Author    string
	Publisher string
	Published utils.CustomDate
	Genre     string
	Language  string
	Pages     int
	Available bool
}
