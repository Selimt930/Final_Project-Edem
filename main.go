package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Edeeeem/Final_Progect/handlers"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Hello, World!")
}

// Elasticsearch Client
var es *elasticsearch.Client

// Initialize Elasticsearch client
func initElasticsearch() {
	var err error

	// Configure the Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Replace with your Elasticsearch URL
		},
		// Optional: Provide authentication if security is enabled
		Username: "elastic",
		Password: "ahLXNssbiZ5f04cwab5f", // Replace with your Elasticsearch password
	}

	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Check the connection
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting Elasticsearch info: %s", err)
	}
	defer res.Body.Close()

	// Parse and log the Elasticsearch version
	var info map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		log.Fatalf("Error parsing Elasticsearch info response: %s", err)
	}
	log.Printf("Connected to Elasticsearch version: %s", info["version"].(map[string]interface{})["number"])
}

func main() {
	// Initialize Elasticsearch
	initElasticsearch()

	// Set up routes
	http.HandleFunc("/", helloHandler)                     // Default route
	http.HandleFunc("/upload", handlers.UploadFileHandler) // File upload
	http.HandleFunc("/books", handlers.HandleBooks)        // Book management
	http.HandleFunc("/books/", handlers.HandleBookByID)    // Individual book operations
	http.HandleFunc("/buy", handlers.HandlePurchase)       // Payment processing

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// )

// func main() {
// 	fmt.Println("h")
// 	app := fiber.New()
// 	log.Fatal(app.Listen(":4000"))
// }
