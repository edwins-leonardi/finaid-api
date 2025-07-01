package http

import (
	"net/http"
	"strconv"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	expenseService port.ExpenseService
}

// NewExpenseHandler creates a new expense handler
func NewExpenseHandler(expenseService port.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// CreateExpense godoc
// @Summary Create a new expense
// @Description Create a new expense with amount, category, subcategory, date, payee, and notes
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body domain.CreateExpenseRequest true "Expense data"
// @Success 201 {object} domain.Expense
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/expenses [post]
func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req domain.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError(c, err)
		return
	}

	expense, err := h.expenseService.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, expense)
}

// GetExpense godoc
// @Summary Get expense by ID
// @Description Get a specific expense by its ID
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Success 200 {object} domain.Expense
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/expenses/{id} [get]
func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	expense, err := h.expenseService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expense)
}

// ListExpenses godoc
// @Summary List expenses
// @Description Get a list of expenses with optional filtering and pagination
// @Tags expenses
// @Accept json
// @Produce json
// @Param skip query int false "Number of expenses to skip" default(0)
// @Param limit query int false "Maximum number of expenses to return" default(10)
// @Param category_id query int false "Filter by expense category ID"
// @Param subcategory_id query int false "Filter by expense subcategory ID"
// @Param payee_id query int false "Filter by payee (person) ID"
// @Param start_date query string false "Filter by start date (YYYY-MM-DD)"
// @Param end_date query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {array} domain.Expense
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/expenses [get]
func (h *ExpenseHandler) ListExpenses(c *gin.Context) {
	var req domain.ListExpensesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationError(c, err)
		return
	}

	expenses, err := h.expenseService.List(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// UpdateExpense godoc
// @Summary Update expense
// @Description Update an existing expense by ID
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Param expense body domain.UpdateExpenseRequest true "Updated expense data"
// @Success 200 {object} domain.Expense
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/expenses/{id} [put]
func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	var req domain.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError(c, err)
		return
	}

	expense, err := h.expenseService.Update(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, expense)
}

// DeleteExpense godoc
// @Summary Delete expense
// @Description Delete an expense by ID
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/expenses/{id} [delete]
func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	err = h.expenseService.Delete(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
