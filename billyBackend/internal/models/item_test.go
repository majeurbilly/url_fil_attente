package models

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
    tests := []struct {
        name     string
        value    string
        validate func(*testing.T, Item)
    }{
        {
            name:  "creates item with value",
            value: "test item",
            validate: func(t *testing.T, item Item) {
                assert.Equal(t, "test item", item.Value)
                assert.NotEmpty(t, item.ID)
                assert.WithinDuration(t, time.Now(), item.CreatedAt, 2*time.Second)
            },
        },
        {
            name:  "creates item with empty value",
            value: "",
            validate: func(t *testing.T, item Item) {
                assert.Empty(t, item.Value)
                assert.NotEmpty(t, item.ID)
                assert.WithinDuration(t, time.Now(), item.CreatedAt, 2*time.Second)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            item := NewItem(tt.value)
            tt.validate(t, item)
        })
    }
}
