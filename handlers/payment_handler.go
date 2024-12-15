package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Edeeeem/Final_Progect/store"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

// HandlePurchase creates a Stripe checkout session for a book
func HandlePurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request to get the book ID
	var request struct {
		BookID string `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find the book in the store
	book, found := store.GetBookByID(request.BookID)
	if !found {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Set Stripe API key
	stripe.Key = "sk_test_51PqLh600u1DnbQizLZTOJ0ddVJEwmXVXQvDR3okzQCq5WGt9kR4ywczehnhku13ei3AFR5iHr4ymj4yNDkbJeEDw00TezeW4wz"

	// Create a Stripe checkout session
	sessionParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(book.Title),
						Description: stripe.String(fmt.Sprintf("Author: %s\n%s", book.Author, book.Description)),
					},
					UnitAmount: stripe.Int64(int64(book.Price * 100)), // Convert to cents
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String("http://localhost:8080/success"),
		CancelURL:  stripe.String("http://localhost:8080/cancel"),
	}

	// Create the session
	stripeSession, err := session.New(sessionParams)
	if err != nil {
		log.Printf("Stripe session error: %v", err)
		http.Error(w, fmt.Sprintf("Error creating Stripe session: %s", err), http.StatusInternalServerError)
		return
	}

	// Respond with the Stripe session URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": stripeSession.URL})
}
