package http

import (
	"log/slog"
	"time"

	"github.com/edwins-leonardi/finaid-api/internal/core/domain"
	"github.com/edwins-leonardi/finaid-api/internal/core/port"
	"github.com/gin-gonic/gin"
)

type PersonHandler struct {
	svc port.PersonService
}

func NewPersonHandler(svc port.PersonService) *PersonHandler {
	return &PersonHandler{
		svc: svc,
	}
}

type createRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required,email" example:"test@example.com"`
}

// Create godoc
//
//	@Summary		Create a new person
//	@Description	create a new person with the provided details
//	@Tags			Persons
//	@Accept			json
//	@Produce		json
//	@Param			createRequest	body		createRequest	true	"Create request"
//	@Success		200				{object}	personResponse	"Person created"
//	@Failure		400				{object}	errorResponse	"Validation error"
//	@Failure		401				{object}	errorResponse	"Unauthorized error"
//	@Failure		404				{object}	errorResponse	"Data not found error"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error"
//	@Router			/persons [post]
func (h *PersonHandler) Create(ctx *gin.Context) {
	slog.Info("Handling create person request")
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	person := domain.Person{
		Name:  req.Name,
		Email: req.Email,
	}

	_, err := h.svc.Create(ctx, &person)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newPersonResponse(&person)

	handleSuccess(ctx, rsp)
}

type listPersonsRequest struct {
	Skip  uint64 `form:"skip" example:"0"`
	Limit uint64 `form:"limit" example:"10"`
}

func (h *PersonHandler) List(ctx *gin.Context) {
	slog.Info("Handling list persons request")
	var req listPersonsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	// Set default value for Limit if not provided
	if req.Limit == 0 {
		req.Limit = 10
	}

	persons, err := h.svc.ListPersons(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	slog.Info("RESPONSE from LIST", "count", len(persons), "skip", req.Skip, "limit", req.Limit)
	var rsp []personResponse
	for _, person := range persons {
		rsp = append(rsp, newPersonResponse(&person))
	}

	slog.Info("Persons listed", "count", len(rsp), "skip", req.Skip, "limit", req.Limit)
	handleSuccess(ctx, rsp)
}

// personResponse represents a person response body
type personResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// newPersonResponse is a helper function to create a response body for handling user data
func newPersonResponse(person *domain.Person) personResponse {
	return personResponse{
		ID:        person.ID,
		Name:      person.Name,
		Email:     person.Email,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}
}
