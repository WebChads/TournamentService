package usecase

import (
	"context"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	storage "github.com/WebChads/TournamentService/internal/storage/pgsql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TournamentRepository interface {
	SelectById(ctx context.Context, tournamentId uuid.UUID) (*dtos.GetTournamentResponse, error)
	SelectByName(ctx context.Context, tournamentName string) ([]dtos.GetTournamentResponse, error)
	SelectByOwn(ctx context.Context, userId uuid.UUID) ([]dtos.GetTournamentResponse, error)
	Insert(ctx context.Context, tournament dtos.CreateTournamentRequest) error
	UpdateById(ctx context.Context, req dtos.CreateTournamentRequest, tournamentId uuid.UUID) error
	DeleteById(ctx context.Context, tournamentId uuid.UUID) error
}

type BidRepository interface {
	Insert(ctx context.Context, tournament dtos.CreateBidRequest) error
}

// All service repositories
type Repositories struct {
	Tournament TournamentRepository
	Bid        BidRepository
	// ...
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Tournament: storage.NewTournamentRepository(db),
		Bid: storage.NewBidRepository(db),
		// ...
	}
}
