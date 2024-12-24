package middleware

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/zap/zaptest"
)

func TestRequestLogger(t *testing.T) {
    gin.SetMode(gin.TestMode)
    logger := zaptest.NewLogger(t)

    tests := []struct {
        name           string
        path           string
        method         string
        expectedStatus int
    }{
        {
            name:           "successful request",
            path:           "/test",
            method:         http.MethodGet,
            expectedStatus: http.StatusOK,
        },
        {
            name:           "not found request",
            path:           "/notfound",
            method:         http.MethodGet,
            expectedStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            router := gin.New()
            router.Use(RequestLogger(logger))

            // Add test endpoint
            router.GET("/test", func(c *gin.Context) {
                c.Status(http.StatusOK)
            })

            w := httptest.NewRecorder()
            req, err := http.NewRequest(tt.method, tt.path, nil)
            require.NoError(t, err)

            router.ServeHTTP(w, req)
            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}

func TestRecovery(t *testing.T) {
    gin.SetMode(gin.TestMode)
    logger := zaptest.NewLogger(t)

    tests := []struct {
        name           string
        handler       gin.HandlerFunc
        expectedStatus int
    }{
        {
            name: "panic recovery",
            handler: func(c *gin.Context) {
                panic("test panic")
            },
            expectedStatus: http.StatusInternalServerError,
        },
        {
            name: "normal handler",
            handler: func(c *gin.Context) {
                c.Status(http.StatusOK)
            },
            expectedStatus: http.StatusOK,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            router := gin.New()
            router.Use(Recovery(logger))
            router.GET("/test", tt.handler)

            w := httptest.NewRecorder()
            req, err := http.NewRequest(http.MethodGet, "/test", nil)
            require.NoError(t, err)

            router.ServeHTTP(w, req)
            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}
