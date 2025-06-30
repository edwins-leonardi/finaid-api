package http

import (
	"net/http"
	"strconv"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

// ExpenseSubCategoryHandler handles HTTP requests for expense subcategories
type ExpenseSubCategoryHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type expenseSubCategoryHandler struct {
	service port.ExpenseSubCategoryService
}

// NewExpenseSubCategoryHandler creates a new expense subcategory HTTP handler
func NewExpenseSubCategoryHandler(service port.ExpenseSubCategoryService) ExpenseSubCategoryHandler {
	return &expenseSubCategoryHandler{
		service: service,
	}
}

// Create handles POST /expenses/categories/subcategories
// @Summary Create a new expense subcategory
// @Description Create a new expense subcategory with the provided information
// @Tags expense-subcategories
// @Accept json
// @Produce json
// @Param subcategory body domain.CreateExpenseSubCategoryRequest true "Expense subcategory data"
// @Success 201 {object} response{data=domain.ExpenseSubCategory}
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /expenses/categories/subcategories [post]
func (h *expenseSubCategoryHandler) Create(c *gin.Context) {
	var req domain.CreateExpenseSubCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError(c, err)
		return
	}

	subcategory, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	rsp := newResponse(true, "Expense subcategory created successfully", subcategory)
	c.JSON(http.StatusCreated, rsp)
}

// List handles GET /expenses/categories/subcategories
// @Summary List expense subcategories
// @Description Get a paginated list of expense subcategories, optionally filtered by expense category
// @Tags expense-subcategories
// @Accept json
// @Produce json
// @Param skip query int false "Number of records to skip" default(0)
// @Param limit query int false "Maximum number of records to return" default(10)
// @Param expense_category_id query int false "Filter by expense category ID"
// @Success 200 {object} response{data=[]domain.ExpenseSubCategory}
// @Failure 500 {object} errorResponse
// @Router /expenses/categories/subcategories [get]
func (h *expenseSubCategoryHandler) List(c *gin.Context) {
	var req domain.ListExpenseSubCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		validationError(c, err)
		return
	}

	subcategories, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	handleSuccess(c, subcategories)
}

// GetByID handles GET /expenses/categories/subcategories/:id
// @Summary Get expense subcategory by ID
// @Description Get a specific expense subcategory by its ID
// @Tags expense-subcategories
// @Accept json
// @Produce json
// @Param id path int true "Expense subcategory ID"
// @Success 200 {object} response{data=domain.ExpenseSubCategory}
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /expenses/categories/subcategories/{id} [get]
func (h *expenseSubCategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	subcategory, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	handleSuccess(c, subcategory)
}

// Update handles PUT /expenses/categories/subcategories/:id
// @Summary Update expense subcategory
// @Description Update an existing expense subcategory with the provided information
// @Tags expense-subcategories
// @Accept json
// @Produce json
// @Param id path int true "Expense subcategory ID"
// @Param subcategory body domain.UpdateExpenseSubCategoryRequest true "Updated expense subcategory data"
// @Success 200 {object} response{data=domain.ExpenseSubCategory}
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /expenses/categories/subcategories/{id} [put]
func (h *expenseSubCategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		validationError(c, err)
		return
	}

	var req domain.UpdateExpenseSubCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationError(c, err)
		return
	}

	subcategory, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	rsp := newResponse(true, "Expense subcategory updated successfully", subcategory)
	c.JSON(http.StatusOK, rsp)
}

// Delete handles DELETE /expenses/categories/subcategories/:id
// @Summary Delete expense subcategory
// @Description Delete an expense subcategory by its ID
// @Tags expense-subcategories
// @Accept json
// @Produce json
// @Param id path int true "Expense subcategory ID"
// @Success 204 "No Content"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /expenses/categories/subcategories/{id} [delete]
func (h *expenseSubCategoryHandler) Delete(c *gin.Context) {
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
