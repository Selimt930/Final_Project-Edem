package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Directory to store uploaded files
const uploadDir = "./uploads"

// UploadFileHandler handles file upload requests
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating upload directory: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Save the file to the uploads directory
	filePath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving the file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing file: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}
