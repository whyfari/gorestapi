package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// any route handler function has to take a req and res
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe in production since a duplicate can be generated
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"] // id decoded from the json param will be ignored, uncomment to use it instead
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {

	// Init Router
	r := mux.NewRouter()

	// Mock data ~todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "1111", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "1112", Title: "Book Two", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "3", Isbn: "1113", Title: "Book Three", Author: &Author{Firstname: "Will", Lastname: "Smith"}})

	// Router Handers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	/* tested using Postman chrome extention
	* getBooks: 'GET', url: http://localhost:8000/api/books, Header: 'Content-type' : application/json // got back all books
	* getBook: 'GET', url: http://localhost:8000/api/books/2, Header: 'Content-type' : application/json // got back book with id 2
	* createBook: 'POST', url: http://localhost:8000/api/books, Header: 'Content-type' : application/json, body(raw)
		 {
			"isbn" : "5555",
			"title" : "Book cuatro",
			"author": {"firstname":"Jio", "lastname": "Dio"}
		}
	 // got back the new book created with a random id
	* getBook: 'UPDATE', url: http://localhost:8000/api/books/2, Header: 'Content-type' : application/json // got back all books including the updated book 2
	*
	* {
    *  "isbn": "1112-2",
    *  "title": "Book 2 ",
    *  "author": {
    *     "firstname": "Jonny",
    *     "lastname": "Doe"
    *   }
    * }
	*/

	log.Fatal(http.ListenAndServe(":8000", r))
}
