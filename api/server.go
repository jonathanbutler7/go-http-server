package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UpdateItem struct {
	Name string `json:"name"`
}

type Server struct {
	*mux.Router

	shoppingItems []Item
}

func NewServer() *Server {
	s := &Server{
		Router:        mux.NewRouter(),
		shoppingItems: []Item{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/shopping-items", s.listShoppingItems()).Methods("GET")
	s.HandleFunc("/shopping-items", s.createShoppingItems()).Methods(("POST"))
	s.HandleFunc("/shopping-items/{id}", s.removeShoppingItem()).Methods("DELETE")
	s.HandleFunc("/shopping-items/{id}", s.updateShoppingItem()).Methods("PUT")
}

func (s *Server) updateShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL
		idStr := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		// Decode the updated item data from the request body
		var updatedData UpdateItem
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		log.Printf("Decoded request body: %+v\n", updatedData)
		// Update the item in the shoppingItems slice
		found := false
		var updatedItem Item
		for i, item := range s.shoppingItems {
			if item.ID == id {
				// Update the item's name and preserve the original ID
				s.shoppingItems[i].Name = updatedData.Name
				updatedItem = s.shoppingItems[i] // Keep the updated item with original ID
				found = true
				break
			}
		}

		// If the item was not found, return a 404 error
		if !found {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		// Return the updated item as the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(updatedItem); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) removeShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		for i, item := range s.shoppingItems {
			if item.ID == id {
				s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
				break
			}
		}
	}
}

func (s *Server) createShoppingItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i Item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		s.shoppingItems = append(s.shoppingItems, i)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listShoppingItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.shoppingItems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
