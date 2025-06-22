package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/edwins-leonardi/finaid-api/internal/adapter/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Router struct {
	*gin.Engine
}

// NewRouter creates a new HTTP router
func NewRouter(config *config.HTTP) (*Router, error) {

	// Disable debug mode in production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
	v1 := router.Group("/api/v1")
	{
		hello := v1.Group("/hello")
		{
			hello.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Hello from Finaid API!",
				})
			})
		}
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
