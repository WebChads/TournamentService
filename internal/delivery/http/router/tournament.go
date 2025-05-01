package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/WebChads/TournamentService/internal/config"
	"github.com/WebChads/TournamentService/internal/models/dtos"
	response "github.com/WebChads/TournamentService/internal/pkg/api"
	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type TournamentUsecase interface {
	Create(ctx context.Context, dto dtos.CreateTournamentRequest) error
	GetById(ctx context.Context, tournamentId string) (*dtos.GetTournamentResponse, error)
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

	r.defaultHandler.Post("/api/v1/tournament/create-tournament", r.CreateTournamentHandler)
	r.defaultHandler.Get("/api/v1/tournament/get-one-tournament/{id}", r.GetTournamentByIdHandler)
	// ...
}

func (a *TournamentRouter) GetTournamentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		a.logger.Error("tournament id param is empty")

		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	tournament, err := a.usecase.GetById(ctx, tournamentId)
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

func (a *TournamentRouter) CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	var request dtos.CreateTournamentRequest

	// Serialize tournament info using DTO
	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		// EOF means there is no data in the request body
		if errors.Is(err, io.EOF) {
			a.logger.Error("request body is empty", slogerr.Error(err))
			response.JSON(w, http.StatusBadRequest, "request body is empty")
			return
		}

		a.logger.Error("failed to decode request body", slogerr.Error(err))
		response.JSON(w, http.StatusBadRequest, "failed to decode request body")
		return
	}

	// Validate request fields
	err = validator.New().Struct(request)
	if err != nil {
		var errors []string
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				errors = append(errors, getValidationMsg(fieldErr))
			}
		}

		resp := map[string]any{"errors": errors}

		response.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = a.usecase.Create(ctx, request)
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

func getValidationMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	}

	return "validation error"
}
