package http

import (
	"net/http"
	"strconv"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

// ExpenseCategoryHandler handles HTTP requests for expense categories
type ExpenseCategoryHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type expenseCategoryHandler struct {
	service port.ExpenseCategoryService
}

// NewExpenseCategoryHandler creates a new expense category HTTP handler
func NewExpenseCategoryHandler(service port.ExpenseCategoryService) ExpenseCategoryHandler {
	return &expenseCategoryHandler{
		service: service,
	}
}

// Create handles POST /expense-categories
// @Summary Create a new expense category
// @Description Create a new expense category with the provided information
// @Tags expense-categories
// @Accept json
// @Produce json
// @Param category body domain.CreateExpenseCategoryRequest true "Expense category data"
// @Success 201 {object} ResponseData{data=domain.ExpenseCategory}
// @Failure 400 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /expense-categories [post]
func (h *expenseCategoryHandler) Create(c *gin.Context) {
	var req domain.CreateExpenseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError(c, err)
		return
	}

	category, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	rsp := newResponse(true, "Expense category created successfully", category)
	c.JSON(http.StatusCreated, rsp)
}

// List handles GET /expense-categories
// @Summary List expense categories
// @Description Get a paginated list of expense categories
// @Tags expense-categories
// @Accept json
// @Produce json
// @Param skip query int false "Number of records to skip" default(0)
// @Param limit query int false "Maximum number of records to return" default(10)
// @Success 200 {object} ResponseData{data=[]domain.ExpenseCategory}
// @Failure 500 {object} ResponseError
// @Router /expense-categories [get]
func (h *expenseCategoryHandler) List(c *gin.Context) {
	var req domain.ListExpenseCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationError(c, err)
		return
	}

	categories, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	handleSuccess(c, categories)
}

// GetByID handles GET /expense-categories/:id
// @Summary Get expense category by ID
// @Description Get a specific expense category by its ID
// @Tags expense-categories
// @Accept json
// @Produce json
// @Param id path int true "Expense category ID"
// @Success 200 {object} ResponseData{data=domain.ExpenseCategory}
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /expense-categories/{id} [get]
func (h *expenseCategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	category, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	handleSuccess(c, category)
}

// Update handles PUT /expense-categories/:id
// @Summary Update expense category
// @Description Update an existing expense category with the provided information
// @Tags expense-categories
// @Accept json
// @Produce json
// @Param id path int true "Expense category ID"
// @Param category body domain.UpdateExpenseCategoryRequest true "Updated expense category data"
// @Success 200 {object} ResponseData{data=domain.ExpenseCategory}
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /expense-categories/{id} [put]
func (h *expenseCategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	var req domain.UpdateExpenseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError(c, err)
		return
	}

	category, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	rsp := newResponse(true, "Expense category updated successfully", category)
	c.JSON(http.StatusOK, rsp)
}

// Delete handles DELETE /expense-categories/:id
// @Summary Delete expense category
// @Description Delete an expense category by its ID
// @Tags expense-categories
// @Accept json
// @Produce json
// @Param id path int true "Expense category ID"
// @Success 204 "No Content"
// @Failure 400 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /expense-categories/{id} [delete]
func (h *expenseCategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
