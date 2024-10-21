package api_test

import (
	// "bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/m/api"
	"github.com/google/uuid"
)

type MockItemsDB struct {
	shoppingItems []*api.Item
	Err           error
}

type ItemsDB interface {
	ListShoppingItems() ([]*api.Item, error)
	// CreateShoppingItem() (*api.Item, error)
	// CreateShoppingItem(item *api.Item) error
}

func (m *MockItemsDB) ListShoppingItems() ([]*api.Item, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if m.shoppingItems == nil {
		return nil, errors.New("no shopping items found")
	}

	return m.shoppingItems, nil
}

func GetItemsHandler(db ItemsDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item, err := db.ListShoppingItems()
		if err != nil {
			http.Error(w, "no shopping items found", http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(item)
	}
}

// CreateItemHandler calls the CreateShoppingItem function
// func CreateItemHandler(db ItemsDB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var newItem api.Item
// 		if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
// 			http.Error(w, "invalid request payload", http.StatusBadRequest)
// 			return
// 		}
// 		newItem.ID = uuid.New()

// 		if err := db.CreateShoppingItem(&newItem); err != nil {
// 			http.Error(w, "failed to create shopping item", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(newItem)
// 	}
// }

// func TestCreateItemHandler(t *testing.T) {
// 	t.Run("successful creation of a shopping item", func(t *testing.T) {
// 		mockDB := &MockItemsDB{}
// 		handler := CreateItemHandler(mockDB)

// 		newItem := api.Item{Name: "Almond butter"}
// 		itemJSON, _ := json.Marshal(newItem)

// 		req, _ := http.NewRequest("POST", "/shopping-items", bytes.NewBuffer(itemJSON))
// 		rr := httptest.NewRecorder()
// 		handler.ServeHTTP(rr, req)

// 		if status := rr.Code; status != http.StatusCreated {
// 			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
// 		}

// 		// Read the response body
// 		bodyBytes, err := io.ReadAll(rr.Body)
// 		if err != nil {
// 			t.Fatalf("failed to read response body: %v", err)
// 		}

// 		// Unmarshal the response into returnedItem
// 		var returnedItem api.Item
// 		if err := json.Unmarshal(bodyBytes, &returnedItem); err != nil {
// 			t.Fatalf("failed to unmarshal response body: %v", err)
// 		}

// 		// Validate the returned item
// 		if returnedItem.Name != newItem.Name {
// 			t.Errorf("handler returned unexpected item name: got %v want %v", returnedItem.Name, newItem.Name)
// 		}

// 		// Check if the item was added to the mock database
// 		if len(mockDB.shoppingItems) != 1 || mockDB.shoppingItems[0].Name != newItem.Name {
// 			t.Errorf("item was not added to the mock database")
// 		}
// 	})

// 	t.Run("failure due to bad request payload", func(t *testing.T) {
// 		mockDB := &MockItemsDB{}
// 		handler := CreateItemHandler(mockDB)

// 		// Invalid JSON
// 		req, _ := http.NewRequest("POST", "/shopping-items", bytes.NewBuffer([]byte("{invalid-json")))
// 		rr := httptest.NewRecorder()
// 		handler.ServeHTTP(rr, req)

// 		if status := rr.Code; status != http.StatusBadRequest {
// 			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
// 		}
// 	})
// }

func TestGetItemsHandler(t *testing.T) {
	t.Run("successful retrieval of shopping items list", func(t *testing.T) {
		items := make([]*api.Item, 2)
		items[0] = &api.Item{ID: uuid.New(), Name: "Peanut butter"}
		items[1] = &api.Item{ID: uuid.New(), Name: "Real butter"}

		mockDB := &MockItemsDB{
			shoppingItems: items,
		}
		handler := GetItemsHandler(mockDB)
		req, _ := http.NewRequest("GET", "/shopping-items", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		// Read the response body as bytes
		bodyBytes, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		// Unmarshal the JSON response into returnedItems slice
		var returnedItems []*api.Item
		if err := json.Unmarshal(bodyBytes, &returnedItems); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}

		// Check if the returned length is correct
		if len(returnedItems) != len(items) {
			t.Errorf("handler returned wrong number of items: got %v want %v", len(returnedItems), len(items))
		}
	})
}

// this server_test file is based on an article from twilio: https://www.twilio.com/en-us/blog/how-to-test-go-http-handlers
