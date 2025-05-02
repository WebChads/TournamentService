package router

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WebChads/TournamentService/internal/config"
	"github.com/WebChads/TournamentService/internal/models/dtos"
	response "github.com/WebChads/TournamentService/internal/pkg/api"
	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
	validate "github.com/WebChads/TournamentService/internal/pkg/validator"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type TournamentUsecase interface {
	GetById(ctx context.Context, tournamentId string) (*dtos.GetTournamentResponse, error)
	GetByName(ctx context.Context, tournamentId string) ([]dtos.GetTournamentResponse, error)
	GetByOwn(ctx context.Context, userId uuid.UUID) ([]dtos.GetTournamentResponse, error)
	Create(ctx context.Context, dto dtos.CreateTournamentRequest) error
	UpdateById(ctx context.Context, request_body io.ReadCloser, tournamentId string) error
	DeleteById(ctx context.Context, tournamentId string) error
}

type TournamentRouter struct {
	defaultHandler *chi.Mux
	logger         *slog.Logger
	config         *config.ServerConfig
	usecase        TournamentUsecase
}

func NewTournamentRouter(r *chi.Mux, cfg *config.ServerConfig,
	log *slog.Logger, usecase TournamentUsecase) *TournamentRouter {
	router := &TournamentRouter{
		defaultHandler: r,
		logger:         log,
		config:         cfg,
		usecase:        usecase,
	}

	return router
}

func ConfigureTournamentRouter(r *TournamentRouter) {
	// Auth middleware
	// r.defaultHandler.Use(auth.AuthMiddleware)

	r.defaultHandler.Get("/api/v1/tournament/get-one-tournament/{id}", r.GetTournamentByIdHandler)
	r.defaultHandler.Get("/api/v1/tournament/get-tournaments/{name}", r.GetTournamentByNameHandler)
	r.defaultHandler.Get("/api/v1/tournament/get-tournaments", r.GetTournamentByOwnHandler)

	r.defaultHandler.Post("/api/v1/tournament/create-tournament", r.CreateTournamentHandler)

	r.defaultHandler.Patch("/api/v1/tournament/update-tournament/{id}", r.UpdateTournamentByIdHandler)

	r.defaultHandler.Delete("/api/v1/tournament/delete-tournament/{id}", r.DeleteTournamentByIdHandler)
	// ...
}

// GetTournamentByIdHandler godoc
// @Title GetTournamentByIdHandler
// @Summary Get tournament by tournament id
// @Tags Tournaments
// @Accept json
// @Produce json
// @Param id path string true "Tournament Id" format(uuid)
// @Success 200 {object} dtos.GetTournamentResponse "Successfully get tournament by id"
// @Failure 400 {object} dtos.Response "Tournament id param is empty"
// @Failure 400 {object} dtos.Response "No tournament with such id"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournament/get-one-tournament/{id} [get]
func (t *TournamentRouter) GetTournamentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		t.logger.Error("tournament id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	tournament, err := t.usecase.GetById(ctx, tournamentId)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	response.JSON(w, http.StatusOK, tournament)
}

// GetTournamentByNameHandler godoc
// @Title GetTournamentByNameHandler
// @Summary Get any tournaments contains a given name
// @Tags Tournaments
// @Accept json
// @Produce json
// @Param name path string true "Tournament name" format(string)
// @Success 200 {object} []dtos.GetTournamentResponse "Successfully get tournaments by name"
// @Failure 400 {object} dtos.Response "Tournament name param is empty"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournament/get-tournaments/{name} [get]
func (t *TournamentRouter) GetTournamentByNameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentName := chi.URLParam(r, "name")
	if tournamentName == "" {
		t.logger.Error("tournament name param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	tournaments, err := t.usecase.GetByName(ctx, tournamentName)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	response.JSON(w, http.StatusOK, tournaments)
}

// GetTournamentByOwnHandler godoc
// @Title GetTournamentByOwnHandler
// @Summary Get any tournaments by owner user id
// @Tags Tournaments
// @Accept json
// @Produce json
// @Success 200 {object} []dtos.GetTournamentResponse "Successfully get tournaments by owner user id"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournament/get-tournaments [get]
func (t *TournamentRouter) GetTournamentByOwnHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	// Get user id from request context (auth middleware)
	// request.UserId = r.Context().Value("user_id").(uuid.UUID)

	// Use mock data to test just for now
	userId, err := uuid.Parse("a8c75926-7a6e-4263-90f0-aecd383f8242")
	if err != nil {
		t.logger.Error("user id conversion", slogerr.Error(err))
		response.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	tournaments, err := t.usecase.GetByOwn(ctx, userId)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	response.JSON(w, http.StatusOK, tournaments)
}

// CreateTournamentHandler godoc
// @Title CreateTournamentHandler
// @Summary Create new tournament
// @Tags Tournaments
// @Accept json
// @Produce json
// @Param request body dtos.CreateTournamentRequest true "Dto with tournament information"
// @Success 200 "Successfully create new tournament"
// @Failure 400 {object} dtos.Response "Request body is empty"
// @Failure 400 {object} dtos.Response "Failed to decode request body"
// @Failure 400 {object} dtos.Response "Validation error"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournament/create-tournament [post]
func (t *TournamentRouter) CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	var request dtos.CreateTournamentRequest

	// Serialize tournament info using DTO
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		// EOF means there is no data in the request body
		if errors.Is(err, io.EOF) {
			t.logger.Error("request body is empty", slogerr.Error(err))
			response.JSON(w, http.StatusBadRequest, "request body is empty")
			return
		}

		t.logger.Error("failed to decode request body", slogerr.Error(err))
		response.JSON(w, http.StatusBadRequest, "failed to decode request body")
		return
	}

	// Get user id from request context (auth middleware)
	// request.UserId = r.Context().Value("user_id").(uuid.UUID)

	// Use mock data just for now
	request.UserId = uuid.New()

	// Validate request fields
	errs := validate.ValidateRequestBody(request)
	if errs != nil {
		response.JSON(w, http.StatusBadRequest, errs)
		return
	}

	err = t.usecase.Create(ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			response.JSON(w, http.StatusRequestTimeout, err.Error())
		} else if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}
}

// UpdateTournamentByIdHandler godoc
// @Title UpdateTournamentByIdHandler
// @Summary Update tournament by tournament id
// @Tags Tournaments
// @Accept json
// @Produce json
// @Param id path string true "Tournament id" format(uuid)
// @Param request body dtos.CreateTournamentRequest true "Dto with tournament information"
// @Success 200 "Successfully update tournament by id"
// @Failure 400 {object} dtos.Response "Tournament id param is empty"
// @Failure 400 {object} dtos.Response "No tournament with such id"
// @Failure 400 {object} dtos.Response "Request body is empty"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournament/update-tournament/{id} [patch]
func (t *TournamentRouter) UpdateTournamentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		t.logger.Error("tournament id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	err := t.usecase.UpdateById(ctx, r.Body, tournamentId)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}
}

// DeleteTournamentByIdHandler godoc
// @Title DeleteTournamentByIdHandler
// @Summary Delete tournament by tournament id
// @Tags Tournaments
// @Accept json
// @Produce json
// @Param id path string true "Tournament id" format(uuid)
// @Success 200 "Successfully delete tournament by id"
// @Failure 400 {object} dtos.Response "Tournament id param is empty"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournament/delete-tournament/{id} [delete]
func (t *TournamentRouter) DeleteTournamentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		t.logger.Error("tournament id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	err := t.usecase.DeleteById(ctx, tournamentId)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}
}
