package books

import "net/http"

type Handler interface {
	CreateBook(res http.ResponseWriter, req *http.Request)
	UpdateBook(res http.ResponseWriter, req *http.Request)
	GetBooks(res http.ResponseWriter, req *http.Request)
	GetBookByID(res http.ResponseWriter, req *http.Request)
	DeleteBook(res http.ResponseWriter, req *http.Request)
}
