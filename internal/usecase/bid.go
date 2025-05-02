package usecase

import (
	"context"
	"log/slog"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
)

type BidUsecase struct {
	logger     *slog.Logger
	repository BidRepository
}

func NewBidUsecase(r BidRepository, l *slog.Logger) *BidUsecase {
	return &BidUsecase{
		logger:     l,
		repository: r,
	}
}

func (u *BidUsecase) GetByStatus(
	ctx context.Context, dto dtos.GetBidRequest,
) ([]dtos.GetBidResponse, error) {
	bids, err := u.repository.SelectByStatus(ctx, dto)
	if err != nil {
		u.logger.Error("select bid by id", slogerr.Error(err))
		return nil, err
	}

	return bids, nil
}

func (u *BidUsecase) Create(ctx context.Context, bid dtos.CreateBidRequest) error {
	err := u.repository.Insert(ctx, bid)
	if err != nil {
		u.logger.Error("create bid", slogerr.Error(err))
		return err
	}

	return nil
}
