package main

import (
	"context"
	"log"
	"net/http"
)

// TODO. Реализовать REST API CRUD для коллекции "Книги"

type Book struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Year    int    `json:"year"`
	Country string `json:"country"`
	Pages   int    `json:"pages"`
}

type apiResponse struct {
	InsertedID   any  `json:"insertedID"`
	Acknowledged bool `json:"acknowledged"`
}

var ctx = context.TODO()

func main() {

	// Init DB
	InitDB(ctx)
	defer myDB.client.Disconnect(ctx)
	// defer myDB.coll.Drop(ctx)

	// Start server
	log.Println("Starting server ...")
	http.HandleFunc("POST /create", CreateBook)
	http.HandleFunc("GET /book/{id}", GetBook)
	http.HandleFunc("POST /update/{id}", UpdateBook)
	http.HandleFunc("DELETE /delete/{id}", DeleteBook)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
