package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Edeeeem/Final_Progect/store"
)

// HandleBooks handles listing all books and creating a new book
func HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// List all books
		books := store.GetAllBooks()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	case http.MethodPost:
		// Create a new book
		var book store.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid book data", http.StatusBadRequest)
			return
		}
		store.AddBook(book)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleBookByID handles retrieving, updating, and deleting a book by ID
func HandleBookByID(w http.ResponseWriter, r *http.Request) {
	// Extract the book ID from the URL
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	if id == "" {
		http.Error(w, "Book ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get a book by ID
		book, found := store.GetBookByID(id)
		if !found {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)
	case http.MethodPut:
		// Update a book by ID
		var book store.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid book data", http.StatusBadRequest)
			return
		}
		updatedBook, ok := store.UpdateBook(id, book)
		if !ok {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(updatedBook)
	case http.MethodDelete:
		// Delete a book by ID
		if ok := store.DeleteBook(id); !ok {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
