package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/m/api"
	"github.com/google/uuid"
)

type MockItemsDB struct {
	shoppingItems *api.Item
	Err           error
}

type ItemsDB interface {
	ListShoppingItems() (*api.Item, error)
}

func (m *MockItemsDB) ListShoppingItems() (*api.Item, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if m.shoppingItems == nil {
		return nil, errors.New("no shopping found")
	}

	return m.shoppingItems, nil
}

func GetItemsHandler(db ItemsDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		item, err := db.ListShoppingItems()
		if err != nil {
			http.Error(w, "no items found", http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(item)
	}
}

func TestGetItemsHandler(t *testing.T) {
	t.Run("successful retrieval of empty item", func(t *testing.T) {
		mockDB := &MockItemsDB{
			shoppingItems: &api.Item{ID: uuid.New(), Name: "Peanut butter"},
		}
		handler := GetItemsHandler(mockDB)
		req, _ := http.NewRequest("GET", "shopping-items", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

}
