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
func NewRouter(
	config *config.HTTP,
	personHandler PersonHandler,
	accountHandler AccountHandler,
	expenseCategoryHandler ExpenseCategoryHandler,
	expenseSubCategoryHandler ExpenseSubCategoryHandler,
	expenseHandler ExpenseHandler,
) (*Router, error) {

	// Disable debug mode in production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()

	// Set allowed origins
	if config.AllowedOrigins != "" {
		allowedOrigins := config.AllowedOrigins
		originsList := strings.Split(allowedOrigins, ",")
		slog.Info("Allowed origins", "origins", originsList)
		ginConfig.AllowOrigins = originsList
	} else {
		// Default to allowing all origins in development
		ginConfig.AllowAllOrigins = true
	}

	// Allow credentials and common headers
	ginConfig.AllowCredentials = true
	ginConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

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
		user := v1.Group("/persons")
		{
			user.GET("", personHandler.List)
			user.POST("", personHandler.Create)
			user.GET("/:id", personHandler.GetByID)
			user.PUT("/:id", personHandler.Update)
			user.DELETE("/:id", personHandler.Delete)
		}
		account := v1.Group("/accounts")
		{
			account.GET("", accountHandler.List)
			account.POST("", accountHandler.Create)
			account.GET("/:id", accountHandler.GetByID)
			account.PUT("/:id", accountHandler.Update)
			account.DELETE("/:id", accountHandler.Delete)
		}
		expenses := v1.Group("/expenses")
		{
			// Main expense routes
			expenses.GET("", expenseHandler.ListExpenses)
			expenses.POST("", expenseHandler.CreateExpense)
			expenses.GET("/:id", expenseHandler.GetExpense)
			expenses.PUT("/:id", expenseHandler.UpdateExpense)
			expenses.DELETE("/:id", expenseHandler.DeleteExpense)

			expenseCategory := expenses.Group("/categories")
			{
				expenseCategory.GET("", expenseCategoryHandler.List)
				expenseCategory.POST("", expenseCategoryHandler.Create)
				expenseCategory.GET("/:id", expenseCategoryHandler.GetByID)
				expenseCategory.PUT("/:id", expenseCategoryHandler.Update)
				expenseCategory.DELETE("/:id", expenseCategoryHandler.Delete)

				expenseSubCategory := expenseCategory.Group("/subcategories")
				{
					expenseSubCategory.GET("", expenseSubCategoryHandler.List)
					expenseSubCategory.POST("", expenseSubCategoryHandler.Create)
					expenseSubCategory.GET("/:id", expenseSubCategoryHandler.GetByID)
					expenseSubCategory.PUT("/:id", expenseSubCategoryHandler.Update)
					expenseSubCategory.DELETE("/:id", expenseSubCategoryHandler.Delete)
				}
			}
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
