package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"billyb/internal/config"
	"billyb/internal/handlers"
)

func setupTestServer(t *testing.T) (*gin.Engine, *config.ServerConfig) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	logger := zaptest.NewLogger(t)

	serverConfig := &config.ServerConfig{
		Port:     8080,
		Timeout:  30,
		MaxItems: 5, // Setting a reasonable limit for testing
	}

	itemHandler := handlers.NewItemHandler(logger, serverConfig)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/items", itemHandler.GetItems)
		v1.POST("/items", itemHandler.AddItem)
	}

	return router, serverConfig
}

func TestAPIFlow(t *testing.T) {
	router, cfg := setupTestServer(t)

	// Test adding items up to limit
	itemsToAdd := make([]string, cfg.MaxItems)
	for i := 0; i < cfg.MaxItems; i++ {
		itemsToAdd[i] = fmt.Sprintf("item%d", i+1)
	}

	// Add items up to the limit
	for _, item := range itemsToAdd {
		w := httptest.NewRecorder()
		payload := map[string]string{"value": item}
		payloadBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(payloadBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Item added successfully", response["message"])
	}

	// Try to add one more item - should fail
	w := httptest.NewRecorder()
	payload := map[string]string{"value": "should-fail"}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/items", bytes.NewBuffer(payloadBytes))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var errorResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)
	assert.Contains(t, errorResponse["error"],
		fmt.Sprintf("maximum limit of %d items reached", cfg.MaxItems))

	// Test getting all items
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/api/v1/items", nil)
	require.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify we have exactly MaxItems items
	assert.Len(t, response["items"], cfg.MaxItems)

	// Verify items are in the correct order
	for i, item := range itemsToAdd {
		assert.Equal(t, item, response["items"][i]["value"])
	}
}
