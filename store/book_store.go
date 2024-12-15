package store

import "sync"

// Book represents a book in the store
type Book struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

// In-memory database and a mutex for thread safety
var (
	books = make(map[string]Book)
	mu    sync.RWMutex // Read-write mutex for thread-safe access
)

// GetAllBooks retrieves all books
func GetAllBooks() []Book {
	mu.RLock()
	defer mu.RUnlock()

	bookList := make([]Book, 0, len(books))
	for _, book := range books {
		bookList = append(bookList, book)
	}
	return bookList
}

// GetBookByID retrieves a book by its ID
func GetBookByID(id string) (Book, bool) {
	mu.RLock()
	defer mu.RUnlock()

	book, exists := books[id]
	return book, exists
}

// AddBook adds a new book to the store
func AddBook(book Book) {
	mu.Lock()
	defer mu.Unlock()

	books[book.ID] = book
}

// UpdateBook updates an existing book by its ID
func UpdateBook(id string, updatedBook Book) (Book, bool) {
	mu.Lock()
	defer mu.Unlock()

	_, exists := books[id]
	if !exists {
		return Book{}, false
	}
	updatedBook.ID = id // Ensure the ID remains unchanged
	books[id] = updatedBook
	return updatedBook, true
}

// DeleteBook deletes a book by its ID
func DeleteBook(id string) bool {
	mu.Lock()
	defer mu.Unlock()

	_, exists := books[id]
	if !exists {
		return false
	}
	delete(books, id)
	return true
}
