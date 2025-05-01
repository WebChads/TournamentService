package usecase

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
)

type TournamentUsecase struct {
	logger     *slog.Logger
	repository TournamentRepository
}

func NewTournamentUsecase(r TournamentRepository, l *slog.Logger) *TournamentUsecase {
	return &TournamentUsecase{
		logger:     l,
		repository: r,
	}
}

func (a *TournamentUsecase) GetById(ctx context.Context, tournamentId string) (*dtos.GetTournamentResponse, error) {
	id, err := strconv.Atoi(tournamentId)
	if err != nil {
		a.logger.Error("id conversion error", slogerr.Error(err))
		return nil, err
	}

	tournament, err := a.repository.Select(ctx, id)
	if err != nil {
		a.logger.Error("get tournament", slogerr.Error(err))
		return nil, err
	}

	return tournament, nil
}

func (a *TournamentUsecase) Create(ctx context.Context, req dtos.CreateTournamentRequest) error {
	err := a.repository.Insert(ctx, req)
	if err != nil {
		a.logger.Error("create tournament", slogerr.Error(err))
		return err
	}

	return nil
}
