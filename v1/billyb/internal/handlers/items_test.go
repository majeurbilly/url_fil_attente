package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"billyb/internal/config"
	"billyb/internal/models"
)

func setupTestRouter(t *testing.T, maxItems int) (*gin.Engine, *ItemHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	logger := zaptest.NewLogger(t)

	cfg := &config.ServerConfig{
		MaxItems: maxItems,
	}

	handler := NewItemHandler(logger, cfg)

	return router, handler
}

func TestGetItems(t *testing.T) {
	tests := []struct {
		name           string
		setupItems     []models.Item
		expectedStatus int
		expectedLen    int
	}{
		{
			name:           "empty list",
			setupItems:     []models.Item{},
			expectedStatus: http.StatusOK,
			expectedLen:    0,
		},
		{
			name: "multiple items",
			setupItems: []models.Item{
				models.NewItem("test1"),
				models.NewItem("test2"),
			},
			expectedStatus: http.StatusOK,
			expectedLen:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler := setupTestRouter(t, 10) // Set a reasonable max items limit
			handler.items = tt.setupItems

			router.GET("/items", handler.GetItems)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/items", nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string][]models.Item
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Len(t, response["items"], tt.expectedLen)
		})
	}
}

func TestAddItem(t *testing.T) {
	tests := []struct {
		name             string
		maxItems         int
		initialItems     []models.Item
		payload          interface{}
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:     "valid item",
			maxItems: 10,
			payload: map[string]string{
				"value": "test item",
			},
			expectedStatus: http.StatusCreated,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, "Item added successfully", response["message"])
				item, ok := response["item"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "test item", item["value"])
			},
		},
		{
			name:           "invalid payload - missing value",
			maxItems:       10,
			payload:        map[string]string{},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], "Field validation for 'Value' failed")
			},
		},
		{
			name:     "invalid payload - empty value",
			maxItems: 10,
			payload: map[string]string{
				"value": "",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], "Field validation for 'Value' failed")
			},
		},
		{
			name:     "maximum items limit reached",
			maxItems: 2,
			initialItems: []models.Item{
				models.NewItem("test1"),
				models.NewItem("test2"),
			},
			payload: map[string]string{
				"value": "test3",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["error"], "maximum limit of 2 items reached")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler := setupTestRouter(t, tt.maxItems)
			if tt.initialItems != nil {
				handler.items = tt.initialItems
			}

			router.POST("/items", handler.AddItem)

			payloadBytes, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(payloadBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	tests := []struct {
		name             string
		setupItems       []models.Item
		itemID           string
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "delete existing item",
			setupItems: []models.Item{
				models.NewItem("test1"),
				models.NewItem("test2"),
			},
			itemID:         "", // Will be set in the test
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "Item deleted successfully", response["message"])
				assert.NotEmpty(t, response["item_id"])
			},
		},
		{
			name:           "delete non-existent item",
			setupItems:     []models.Item{models.NewItem("test1")},
			itemID:         "non-existent-id",
			expectedStatus: http.StatusNotFound,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "item not found", response["error"])
			},
		},
		// {
		// 	name:           "delete with empty id",
		// 	setupItems:     []models.Item{models.NewItem("test1")},
		// 	itemID:         "",
		// 	expectedStatus: http.StatusBadRequest,
		// 	validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
		// 		var response map[string]interface{}
		// 		err := json.Unmarshal(w.Body.Bytes(), &response)
		// 		require.NoError(t, err)
		// 		assert.Equal(t, "item id is required", response["error"])
		// 	},
		// },  TODO: fix this case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler := setupTestRouter(t, 10)
			handler.items = tt.setupItems

			// Set up the v1 group for testing
			v1 := router.Group("/api/v1")
			v1.DELETE("/items/:id", handler.DeleteItem)

			// For the "delete existing item" test, use the ID of the first item
			itemID := tt.itemID
			if tt.name == "delete existing item" {
				itemID = tt.setupItems[0].ID
			}

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/api/v1/items/"+itemID, nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}

			// For successful deletion, verify the item was actually removed
			if w.Code == http.StatusOK {
				assert.Len(t, handler.items, len(tt.setupItems)-1)
				// Verify the deleted item is no longer in the slice
				for _, item := range handler.items {
					assert.NotEqual(t, itemID, item.ID)
				}
			}
		})
	}
}

func TestDeleteAllItems(t *testing.T) {
	tests := []struct {
		name             string
		setupItems       []models.Item
		expectedStatus   int
		validateResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "delete all items",
			setupItems: []models.Item{
				models.NewItem("test1"),
				models.NewItem("test2"),
				models.NewItem("test3"),
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "All items deleted successfully", response["message"])
				assert.Equal(t, float64(3), response["count"]) // JSON numbers are float64
			},
		},
		{
			name:           "delete all items when empty",
			setupItems:     []models.Item{},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "All items deleted successfully", response["message"])
				assert.Equal(t, float64(0), response["count"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, handler := setupTestRouter(t, 10)
			handler.items = tt.setupItems

			// Set up the v1 group for testing
			v1 := router.Group("/api/v1")
			v1.DELETE("/items", handler.DeleteAllItems)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/api/v1/items", nil)
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.validateResponse != nil {
				tt.validateResponse(t, w)
			}

			// Verify all items were actually deleted
			assert.Empty(t, handler.items)
		})
	}
}
