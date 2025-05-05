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
	UpdateStatus(ctx context.Context, dto dtos.UpdateBidStatusRequest) error
	DeleteById(ctx context.Context, bidId uuid.UUID) error
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

	r.defaultHandler.Patch("/api/v1/tournaments/{id}/bids/{bid_id}/change-bid-status", r.UpdateBidStatusHandler)

	r.defaultHandler.Delete("/api/v1/tournaments/{id}/bids/{bid_id}/cancel-bid", r.CancelBidHandler)
}

// GetBidByStatusHandler godoc
// @Title GetBidByStatusHandler
// @Summary Get bid by specific bid status
// @Tags Bids
// @Accept json
// @Produce json
// @Param id path string true "Tournament Id" format(uuid)
// @Param bid-status query string true "Bid status" format(int)
// @Success 200 {object} []dtos.GetBidResponse "Successfully get list of needed bids"
// @Failure 400 {object} dtos.Response "tournament id param is empty"
// @Failure 400 {object} dtos.Response "bid status query param is empty"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournaments/{id}/bids/get-all-bids [get]
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
		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user id from request context (auth middleware)
	// request.UserId = r.Context().Value("user_id").(uuid.UUID)

	// Use mock data to test just for now
	userId, err := uuid.Parse("71ec685c-c366-40af-a158-4d9dc4c45892")
	if err != nil {
		b.logger.Error("user id conversion", slogerr.Error(err))
		response.JSON(w, http.StatusInternalServerError, err.Error())
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

// CreateBidHandler godoc
// @Title CreateBidHandler
// @Summary Create new bid
// @Tags Bids
// @Accept json
// @Produce json
// @Param id path string true "Tournament Id" format(uuid)
// @Param request body dtos.CreateBidRequest true "Dto with bid information"
// @Success 200 "Successfully create new bid"
// @Failure 400 {object} dtos.Response "Tournament id param is empty"
// @Failure 400 {object} dtos.Response "Request body is empty"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournaments/{id}/bids/create-bid [post]
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
		response.JSON(w, http.StatusBadRequest, "request body is empty")
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

// UpdateBidStatusHandler godoc
// @Title UpdateBidStatusHandler
// @Summary Update bid status to another one
// @Tags Bids
// @Accept json
// @Produce json
// @Param id path string true "Tournament Id" format(uuid)
// @Param bid_id path string true "Bid Id" format(uuid)
// @Success 200 "Successfully update specific bid"
// @Failure 400 {object} dtos.Response "Tournament id param is empty"
// @Failure 400 {object} dtos.Response "Bid id param is empty"
// @Failure 400 {object} dtos.Response "Request body is empty"
// @Failure 400 {object} dtos.Response "Validation error"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournaments/{id}/bids/{bid_id}/change-bid-status [patch]
func (b *BidRouter) UpdateBidStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		b.logger.Error("tournament id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	bidId := chi.URLParam(r, "bid_id")
	if bidId == "" {
		b.logger.Error("bid id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	var request dtos.UpdateBidStatusRequest

	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		b.logger.Error("request body is empty", slogerr.Error(err))
		response.JSON(w, http.StatusBadRequest, "request body is empty")
		return
	}

	errs := validate.ValidateRequestBody(request)
	if errs != nil {
		response.JSON(w, http.StatusBadRequest, errs)
		return
	}

	request.BidId, err = uuid.Parse(bidId)
	if err != nil {
		b.logger.Error("bid id conversion", slogerr.Error(err))
		response.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = b.usecase.UpdateStatus(ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}
}

// CancelBidHandler godoc
// @Title CancelBidHandler
// @Summary Delete bit from the bid list
// @Tags Bids
// @Accept json
// @Produce json
// @Param id path string true "Tournament Id" format(uuid)
// @Param bid_id path string true "Bid Id" format(uuid)
// @Success 200 "Successfully delete specific bid"
// @Failure 400 {object} dtos.Response "Tournament id param is empty"
// @Failure 400 {object} dtos.Response "Bid id param is empty"
// @Failure 500 {object} dtos.Response "Happened internal error"
// @Router /api/v1/tournaments/{id}/bids/{bid_id}/cancel-bid [delete]
func (b *BidRouter) CancelBidHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
	defer cancel()

	tournamentId := chi.URLParam(r, "id")
	if tournamentId == "" {
		b.logger.Error("tournament id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	bidId := chi.URLParam(r, "bid_id")
	if bidId == "" {
		b.logger.Error("bid id param is empty")
		response.JSON(w, http.StatusBadRequest, "invalid request")
		return
	}

	bidUUID, err := uuid.Parse(bidId)
	if err != nil {
		b.logger.Error("bid id conversion", slogerr.Error(err))
		response.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = b.usecase.DeleteById(ctx, bidUUID)
	if err != nil {
		if strings.Contains(err.Error(), "failed") {
			response.JSON(w, http.StatusInternalServerError, err.Error())
		} else {
			response.JSON(w, http.StatusBadRequest, err.Error())
		}

		return
	}
}
