package usecase

import (
	"context"
	"io"
	"log/slog"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	slogerr "github.com/WebChads/TournamentService/internal/pkg/logger"
	"github.com/go-chi/render"
	"github.com/google/uuid"
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

func (u *TournamentUsecase) GetById(ctx context.Context, tournamentId string) (*dtos.GetTournamentResponse, error) {
	id, err := uuid.Parse(tournamentId)
	if err != nil {
		u.logger.Error("uuid conversion error", slogerr.Error(err))
		return nil, err
	}

	tournament, err := u.repository.SelectById(ctx, id)
	if err != nil {
		u.logger.Error("get tournament by id", slogerr.Error(err))
		return nil, err
	}

	return tournament, nil
}

func (u *TournamentUsecase) GetByName(ctx context.Context, tournamentName string) ([]dtos.GetTournamentResponse, error) {
	tournaments, err := u.repository.SelectByName(ctx, tournamentName)
	if err != nil {
		u.logger.Error("get tournament by name", slogerr.Error(err))
		return nil, err
	}

	return tournaments, nil
}

func (u *TournamentUsecase) GetByOwn(ctx context.Context, userId uuid.UUID) ([]dtos.GetTournamentResponse, error) {
	tournaments, err := u.repository.SelectByOwn(ctx, userId)
	if err != nil {
		u.logger.Error("get tournament by name", slogerr.Error(err))
		return nil, err
	}

	return tournaments, nil
}

func (u *TournamentUsecase) Create(ctx context.Context, req dtos.CreateTournamentRequest) error {
	err := u.repository.Insert(ctx, req)
	if err != nil {
		u.logger.Error("create tournament", slogerr.Error(err))
		return err
	}

	return nil
}

func (u *TournamentUsecase) UpdateById(
	ctx context.Context, request_body io.ReadCloser, tournamentId string,
) error {
	id, err := uuid.Parse(tournamentId)
	if err != nil {
		u.logger.Error("id conversion error", slogerr.Error(err))
		return err
	}

	// Get existing tournament by id
	tournament, err := u.repository.SelectById(ctx, id)
	if err != nil {
		u.logger.Error("get tournament by id", slogerr.Error(err))
		return err
	}

	slog.Info("tournament name", slog.Any("name", tournament.Name))
	slog.Info("tournament date", slog.Any("tournament_date", tournament.TournamentDate))
	slog.Info("tournament matches", slog.Any("matches_amount", tournament.MatchesAmount))
	slog.Info("tournament user id", slog.Any("user_id", tournament.UserId))

	err = render.DecodeJSON(request_body, &tournament)
	if err != nil {
		u.logger.Error("request body is empty", slogerr.Error(err))
		return err
	}

	// Create new tournament request structure
	updatedTournament := dtos.CreateTournamentRequest{
		Name:           tournament.Name,
		TournamentDate: tournament.TournamentDate,
		MatchesAmount:  tournament.MatchesAmount,
		UserId:         tournament.UserId,
	}

	// Update tournament in the database using id
	err = u.repository.UpdateById(ctx, updatedTournament, id)
	if err != nil {
		u.logger.Error("update tournament", slogerr.Error(err))
		return err
	}

	return nil
}

func (u *TournamentUsecase) DeleteById(ctx context.Context, tournamentId string) error {
	id, err := uuid.Parse(tournamentId)
	if err != nil {
		u.logger.Error("id conversion error", slogerr.Error(err))
		return err
	}

	err = u.repository.DeleteById(ctx, id)
	if err != nil {
		u.logger.Error("delete tournament", slogerr.Error(err))
		return err
	}

	return nil
}
