package routes

import (
	"Golang/handler"
	"net/http"
	"github.com/gorilla/mux"
)
	

func Register( r *mux.Router ){
	r.HandleFunc("/books/search", handler.HandleSearch ).Methods(http.MethodGet)
}