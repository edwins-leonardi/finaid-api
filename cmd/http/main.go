package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/edwins-leonardi/finaid-api/internal/adapter/config"
	"github.com/edwins-leonardi/finaid-api/internal/adapter/handler/http"
	"github.com/edwins-leonardi/finaid-api/internal/adapter/logger"
	"github.com/edwins-leonardi/finaid-api/internal/adapter/storage/postgres"
	"github.com/edwins-leonardi/finaid-api/internal/adapter/storage/postgres/repository"
	"github.com/edwins-leonardi/finaid-api/internal/core/service"
)

func main() {
	// Load App configuration
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	logger.Set(config.App)

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	// Initialize database connection
	db, err := postgres.New(context.Background(), config.DB)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Dependency injection
	// Person
	personRepo := repository.NewPersonRepository(db)
	personService := service.NewPersonService(personRepo)
	personHandler := http.NewPersonHandler(personService)

	// Account
	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo, personRepo)
	accountHandler := http.NewAccountHandler(accountService)

	// Expense Category
	expenseCategoryRepo := repository.NewExpenseCategoryRepository(db.Pool)
	expenseCategoryService := service.NewExpenseCategoryService(expenseCategoryRepo, slog.Default())
	expenseCategoryHandler := http.NewExpenseCategoryHandler(expenseCategoryService)

	// Expense SubCategory
	expenseSubCategoryRepo := repository.NewExpenseSubCategoryRepository(db.Pool)
	expenseSubCategoryService := service.NewExpenseSubCategoryService(expenseSubCategoryRepo, expenseCategoryRepo, slog.Default())
	expenseSubCategoryHandler := http.NewExpenseSubCategoryHandler(expenseSubCategoryService)

	// Init router
	router, err := http.NewRouter(
		config.HTTP,
		*personHandler,
		*accountHandler,
		expenseCategoryHandler,
		expenseSubCategoryHandler,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
