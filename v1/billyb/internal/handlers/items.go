package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"billyb/internal/config"
	"billyb/internal/models"
)

type ItemHandler struct {
	items  []models.Item
	logger *zap.Logger
	config *config.ServerConfig
}

func NewItemHandler(logger *zap.Logger, cfg *config.ServerConfig) *ItemHandler {
	return &ItemHandler{
		items:  make([]models.Item, 0),
		logger: logger,
		config: cfg,
	}
}

func (h *ItemHandler) GetItems(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "GetItems")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		h.logger.Error("request timeout", zap.Error(ctx.Err()))
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
		return
	default:
		h.logger.Info("retrieving items", zap.Int("count", len(h.items)))
		c.JSON(http.StatusOK, gin.H{
			"items": h.items,
		})
	}
}

func (h *ItemHandler) AddItem(c *gin.Context) {
	_, span := otel.Tracer("").Start(c.Request.Context(), "AddItem")
	defer span.End()

	// Check if we're at the item limit before processing the request
	if len(h.items) >= h.config.MaxItems {
		h.logger.Warn("cannot add item: maximum items limit reached",
			zap.Int("current_items", len(h.items)),
			zap.Int("max_items", h.config.MaxItems),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": fmt.Sprintf("cannot add more items: maximum limit of %d items reached", h.config.MaxItems),
		})
		return
	}

	var input struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Error("invalid request payload",
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := models.NewItem(input.Value)
	h.items = append(h.items, item)

	h.logger.Info("item added successfully",
		zap.String("item_id", item.ID),
		zap.String("value", item.Value),
		zap.String("client_ip", c.ClientIP()),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Item added successfully",
		"item":    item,
	})
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
    _, span := otel.Tracer("").Start(c.Request.Context(), "DeleteItem")
    defer span.End()

    id := c.Param("id")
    if id == "" {
        h.logger.Error("missing item id",
            zap.String("client_ip", c.ClientIP()),
        )
        c.JSON(http.StatusBadRequest, gin.H{"error": "item id is required"})
        return
    }

    for i, item := range h.items {
        if item.ID == id {
            h.items = append(h.items[:i], h.items[i+1:]...)
            h.logger.Info("item deleted successfully",
                zap.String("item_id", id),
                zap.String("client_ip", c.ClientIP()),
            )
            c.JSON(http.StatusOK, gin.H{
                "message": "Item deleted successfully",
                "item_id": id,
            })
            return
        }
    }

    h.logger.Warn("item not found",
        zap.String("item_id", id),
        zap.String("client_ip", c.ClientIP()),
    )
    c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
}

func (h *ItemHandler) DeleteAllItems(c *gin.Context) {
    _, span := otel.Tracer("").Start(c.Request.Context(), "DeleteAllItems")
    defer span.End()

    itemCount := len(h.items)
    h.items = make([]models.Item, 0)

    h.logger.Info("all items deleted",
        zap.Int("items_deleted", itemCount),
        zap.String("client_ip", c.ClientIP()),
    )
    
    c.JSON(http.StatusOK, gin.H{
        "message": "All items deleted successfully",
        "count":   itemCount,
    })
}