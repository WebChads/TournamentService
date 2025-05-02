package router

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
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

type BidUsecase interface {
	GetByStatus(ctx context.Context, dto dtos.GetBidRequest) ([]dtos.GetBidResponse, error)
	Create(ctx context.Context, dto dtos.CreateBidRequest) error
}

type BidRouter struct {
	defaultHandler *chi.Mux
	logger         *slog.Logger
	config         *config.ServerConfig
	usecase        BidUsecase
}

func NewBidRouter(r *chi.Mux, cfg *config.ServerConfig,
	log *slog.Logger, usecase BidUsecase) *BidRouter {
	router := &BidRouter{
		defaultHandler: r,
		logger:         log,
		config:         cfg,
		usecase:        usecase,
	}

	return router
}

func ConfigureBidRouter(r *BidRouter) {
	// Auth middleware
	// r.defaultHandler.Use(auth.AuthMiddleware)

	r.defaultHandler.Get("/api/v1/tournaments/{id}/bids/get-all-bids", r.GetBidByStatusHandler)

	r.defaultHandler.Post("/api/v1/tournaments/{id}/bids/create-bid", r.CreateBidHandler)
	// ...
}

func (b *BidRouter) GetBidByStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		b.logger.Error("tournament id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	// Get the bid-status from query parameters
	bidStatus := r.URL.Query().Get("bid-status")
	if bidStatus == "" {
		b.logger.Error("bid status query param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}
	bidStatusInt, err := strconv.Atoi(bidStatus)
	if err != nil {
		b.logger.Error("bid status conversion error", slogerr.Error(err))
		return
	}

	// Get user id from request context (auth middleware)
	// request.UserId = r.Context().Value("user_id").(uuid.UUID)

	// Use mock data to test just for now
	userId, err := uuid.Parse("15563154-959a-4bb0-97d0-55449a812998")
	if err != nil {
		b.logger.Error("user id conversion", slogerr.Error(err))
		return
	}

	request := dtos.GetBidRequest{
		UserId:    userId,
		BidStatus: bidStatusInt,
	}

	bids, err := b.usecase.GetByStatus(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	response.JSON(w, http.StatusOK, bids)
}

func (b *BidRouter) CreateBidHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		b.logger.Error("tournament id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	var request dtos.CreateBidRequest

	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		b.logger.Error("request body is empty", slogerr.Error(err))

		return
	}

	errs := validate.ValidateRequestBody(request)
	if errs != nil {
		response.JSON(w, http.StatusBadRequest, errs)
		return
	}

	// Get user id from request context (auth middleware)
	// request.UserId = r.Context().Value("user_id").(uuid.UUID)

	// Use mock data to test just for now
	userId := uuid.New()

	request.UserId = userId
	request.TournamentId, err = uuid.Parse(tournamentId)
	if err != nil {
		b.logger.Error("tournament id conversion", slogerr.Error(err))
		response.JSON(w, http.StatusBadRequest, "invalid tournament id")
		return
	}

	err = b.usecase.Create(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err)
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}
}
