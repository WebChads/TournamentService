package usecase

import (
	"context"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	storage "github.com/WebChads/TournamentService/internal/storage/pgsql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TournamentRepository interface {
	SelectById(ctx context.Context, tournamentId uuid.UUID) (*dtos.GetTournamentByIdResponse, error)
	SelectByName(ctx context.Context, tournamentName string) ([]dtos.GetTournamentByIdResponse, error)
	Insert(ctx context.Context, tournament dtos.CreateTournamentRequest) error
	UpdateById(ctx context.Context, req dtos.CreateTournamentRequest, tournamentId uuid.UUID) error
	DeleteById(ctx context.Context, tournamentId uuid.UUID) error
}

// All service repositories
type Repositories struct {
	Tournament TournamentRepository
	// ...
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Tournament: storage.NewTournamentRepository(db),
		// ...
	}
}
