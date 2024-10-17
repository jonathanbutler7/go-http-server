package api

import (
	"encoding/json"
	"fmt"
	"io"
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

	s.HandleFunc("/shopping-items", s.listShoppingItems).Methods("GET") // since you aren't calling a function that returns a function, you can just directly reference the function like this s.listShoppingItems instead of s.listShoppingItems().
	s.HandleFunc("/shopping-items", s.createShoppingItems()).Methods(("POST"))
	s.HandleFunc("/shopping-items/{id}", s.removeShoppingItem()).Methods("DELETE")
	s.HandleFunc("/shopping-items/{id}", s.updateShoppingItem()).Methods("PUT")

	return s
}

// this function is just defining a function just to move four lines out of the NewServer function - the abstraction
// is not necessary and the code is more readable without it. I'd just move all the code from this function into the
// NewServer function.
//func (s *Server) routes() {
//
//}

func (s *Server) updateShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL
		idStr := mux.Vars(r)["id"]
		if idStr == "" {
			http.Error(w, "ID of the shopping item is required", http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid UUID.\nError: %v\nReceivedID: %s\n", err, idStr), http.StatusBadRequest)
			return
		}

		// Decode the updated item data from the request body
		var updatedData UpdateItem
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil { // change to json.Unmarshal
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			// if you had used json.Unmarshal, you'd be able to print out the 'bad' body here in the error message
			// to give the caller a better idea of what went wrong and what to change instead of just "invalid request payload"
			return
		}
		log.Printf("Decoded request body: %+v\n", updatedData)
		// Update the item in the shoppingItems slice

		var updatedItem *Item // if we use a pointer here, we can check if the item was found by checking if updatedItem is nil, avoiding the extra 'found' bool
		for i, item := range s.shoppingItems {
			if item.ID == id {
				// Update the item's name and preserve the original ID
				s.shoppingItems[i].Name = updatedData.Name
				updatedItem = &s.shoppingItems[i] // Keep the updated item with original ID
				break
			}
		}

		// If the item was not found, return a 404 error
		if updatedItem != nil {
			http.Error(w, fmt.Sprintf("Item with ID: '%s' not found", idStr), http.StatusNotFound)
			return
		}

		// Return the updated item as the response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(updatedItem); err != nil { // change to json.Marshal
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) removeShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		// check for empty idStr first
		if idStr == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

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

		w.WriteHeader(http.StatusNoContent) // 204 no content is what I expect to see on DELETE requests ("I did something but there's nothing to return")
	}
}

func (s *Server) createShoppingItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i Item
		// we tend to exclusively use json.Unmarshal over json.NewDecoder.Decode, but I not positive as to why :)
		// probably because it allows more finite error handling and it is more explicit of what is happening.
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close() // always close the request body. Defer funcs will run after the current scope ends
		// in this case, r.Body.Close() will run after this createShoppingItems function returns.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, &i)
		if err != nil {
			http.Error(w, fmt.Sprintf("error unmarshalling request body to Item\n\tError:%v\n\tReceived body:%s", err, string(body)), http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		s.shoppingItems = append(s.shoppingItems, i)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201 Created is standard for POST requests that create a new resource although 200 does work fine.
		// you know better than most I guess about http status codes, so you can ignore me!

		if err := json.NewEncoder(w).Encode(i); err != nil { // I'd use json.Marshal instead of json.NewEncoder.Encode here to keep it consistent with the rest of the code
			http.Error(w, err.Error(), http.StatusInternalServerError) // use fmt.Sprintf to give a more detailed error message
			return
		}
	}
}

// Having a function that returns a function is a great way to get your brain all scrambled up. Works, but unnecessary headache.
// That technique can be extremely powerful in some cases (most especially for testing), but in this case, it's just adding complexity.
// I recommend turning listShoppingItems itself into a handler func by making it accept a http.ResponseWriter and *http.Request as arguments.
// Here's what I mean:
func (s *Server) listShoppingItems(w http.ResponseWriter, _ *http.Request) {
	// if you don't use an argument, it is good form for readability to "ignore" it by replacing the
	// parameter variable name with an underscore like I did with the request parameter
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(s.shoppingItems) // for reasons I don't really know, we tend to use json.Marshal instead of json.NewEncoder.Encode
	// afaik there is no real difference between the two, but I could be wrong. I does allow you to have more finite control over the
	// error handling/messaging, so that's a plus.
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshalling item\n\tError: %v\n", err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("error writing item to response writer: %v", err), http.StatusInternalServerError)
		return
	}
}

// todo: WRITE TESTS FOR THESE FUNCTIONS TO MAKE SURE THEY ARE WORKING ;)
