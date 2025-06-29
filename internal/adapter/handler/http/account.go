package http

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	svc port.AccountService
}

func NewAccountHandler(svc port.AccountService) *AccountHandler {
	return &AccountHandler{
		svc: svc,
	}
}

type createAccountRequest struct {
	Name           string  `json:"name" binding:"required" example:"Main Checking Account"`
	Currency       string  `json:"currency" binding:"required" example:"USD"`
	AccountType    string  `json:"account_type" binding:"required" example:"checking"`
	InitialBalance float64 `json:"initial_balance" example:"1000.50"`
	PrimaryOwnerID uint64  `json:"primary_owner_id" binding:"required" example:"1"`
	SecondOwnerID  *uint64 `json:"second_owner_id,omitempty" example:"2"`
}

// Create godoc
//
//	@Summary		Create a new account
//	@Description	create a new account with the provided details
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			createAccountRequest	body		createAccountRequest	true	"Create account request"
//	@Success		200						{object}	accountResponse			"Account created"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/accounts [post]
func (h *AccountHandler) Create(ctx *gin.Context) {
	slog.Info("Handling create account request")
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	account := domain.Account{
		Name:           req.Name,
		Currency:       req.Currency,
		AccountType:    req.AccountType,
		InitialBalance: req.InitialBalance,
		PrimaryOwnerID: req.PrimaryOwnerID,
		SecondOwnerID:  req.SecondOwnerID,
	}

	_, err := h.svc.Create(ctx, &account)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAccountResponse(&account)
	handleSuccess(ctx, rsp)
}

type listAccountsRequest struct {
	Skip  uint64 `form:"skip" example:"0"`
	Limit uint64 `form:"limit" example:"10"`
}

// List godoc
//
//	@Summary		List accounts
//	@Description	get a list of accounts with pagination
//	@Tags			Accounts
//	@Produce		json
//	@Param			skip	query		int	false	"Number of accounts to skip"
//	@Param			limit	query		int	false	"Number of accounts to return"
//	@Success		200		{array}		accountResponse	"Accounts listed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		401		{object}	errorResponse	"Unauthorized error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/accounts [get]
func (h *AccountHandler) List(ctx *gin.Context) {
	slog.Info("Handling list accounts request")
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	// Set default value for Limit if not provided
	if req.Limit == 0 {
		req.Limit = 10
	}

	accounts, err := h.svc.ListAccounts(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	slog.Info("RESPONSE from LIST", "count", len(accounts), "skip", req.Skip, "limit", req.Limit)
	var rsp []accountResponse
	for _, account := range accounts {
		rsp = append(rsp, newAccountResponse(&account))
	}

	slog.Info("Accounts listed", "count", len(rsp), "skip", req.Skip, "limit", req.Limit)
	handleSuccess(ctx, rsp)
}

// GetByID godoc
//
//	@Summary		Get an account by ID
//	@Description	get an account by their unique identifier
//	@Tags			Accounts
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	accountResponse	"Account found"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/accounts/{id} [get]
func (h *AccountHandler) GetByID(ctx *gin.Context) {
	slog.Info("Handling get account by ID request")

	// Get account ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		validationError(ctx, err)
		return
	}

	account, err := h.svc.GetAccount(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAccountResponse(account)
	handleSuccess(ctx, rsp)
}

type updateAccountRequest struct {
	Name           string  `json:"name" binding:"required" example:"Updated Checking Account"`
	Currency       string  `json:"currency" binding:"required" example:"USD"`
	AccountType    string  `json:"account_type" binding:"required" example:"savings"`
	InitialBalance float64 `json:"initial_balance" example:"2000.75"`
	PrimaryOwnerID uint64  `json:"primary_owner_id" binding:"required" example:"1"`
	SecondOwnerID  *uint64 `json:"second_owner_id,omitempty" example:"2"`
}

// Update godoc
//
//	@Summary		Update an account
//	@Description	update an existing account with the provided details
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id					path		int						true	"Account ID"
//	@Param			updateAccountRequest	body		updateAccountRequest	true	"Update account request"
//	@Success		200					{object}	accountResponse			"Account updated"
//	@Failure		400					{object}	errorResponse			"Validation error"
//	@Failure		401					{object}	errorResponse			"Unauthorized error"
//	@Failure		404					{object}	errorResponse			"Data not found error"
//	@Failure		409					{object}	errorResponse			"Data conflict error"
//	@Failure		500					{object}	errorResponse			"Internal server error"
//	@Router			/accounts/{id} [put]
func (h *AccountHandler) Update(ctx *gin.Context) {
	slog.Info("Handling update account request")

	// Get account ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		validationError(ctx, err)
		return
	}

	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	account := domain.Account{
		ID:             id,
		Name:           req.Name,
		Currency:       req.Currency,
		AccountType:    req.AccountType,
		InitialBalance: req.InitialBalance,
		PrimaryOwnerID: req.PrimaryOwnerID,
		SecondOwnerID:  req.SecondOwnerID,
	}

	updatedAccount, err := h.svc.UpdateAccount(ctx, &account)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newAccountResponse(updatedAccount)
	handleSuccess(ctx, rsp)
}

// Delete godoc
//
//	@Summary		Delete an account
//	@Description	delete an existing account by ID
//	@Tags			Accounts
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	map[string]string	"Account deleted successfully"
//	@Failure		400	{object}	errorResponse		"Validation error"
//	@Failure		401	{object}	errorResponse		"Unauthorized error"
//	@Failure		404	{object}	errorResponse		"Data not found error"
//	@Failure		500	{object}	errorResponse		"Internal server error"
//	@Router			/accounts/{id} [delete]
func (h *AccountHandler) Delete(ctx *gin.Context) {
	slog.Info("Handling delete account request")

	// Get account ID from URL parameter
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		validationError(ctx, err)
		return
	}

	err = h.svc.DeleteAccount(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, gin.H{"message": "Account deleted successfully"})
}

// accountResponse represents an account response body
type accountResponse struct {
	ID             uint64    `json:"id" example:"1"`
	Name           string    `json:"name" example:"Main Checking Account"`
	Currency       string    `json:"currency" example:"USD"`
	AccountType    string    `json:"account_type" example:"checking"`
	InitialBalance float64   `json:"initial_balance" example:"1000.50"`
	PrimaryOwnerID uint64    `json:"primary_owner_id" example:"1"`
	SecondOwnerID  *uint64   `json:"second_owner_id,omitempty" example:"2"`
	CreatedAt      time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt      time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// newAccountResponse is a helper function to create a response body for handling account data
func newAccountResponse(account *domain.Account) accountResponse {
	return accountResponse{
		ID:             account.ID,
		Name:           account.Name,
		Currency:       account.Currency,
		AccountType:    account.AccountType,
		InitialBalance: account.InitialBalance,
		PrimaryOwnerID: account.PrimaryOwnerID,
		SecondOwnerID:  account.SecondOwnerID,
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
	}
}
