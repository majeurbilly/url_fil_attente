package models

import (
    "time"

    "github.com/google/uuid"
)

type Item struct {
    ID        string    `json:"id"`
    Value     string    `json:"value"`
    CreatedAt time.Time `json:"created_at"`
}

func NewItem(value string) Item {
    return Item{
        ID:        uuid.New().String(),
        Value:     value,
        CreatedAt: time.Now(),
    }
}
