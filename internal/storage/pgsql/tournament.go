package storage

import (
	"context"
	"errors"
	"time"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	"github.com/jmoiron/sqlx"
)

type TournamentRepository struct {
	db *sqlx.DB
}

func NewTournamentRepository(db *sqlx.DB) *TournamentRepository {
	return &TournamentRepository{
		db: db,
	}
}

// Database inner structure
type Tournament struct {
	Id             int       `db:"id"`
	Name           string    `db:"name"`
	TournamentDate time.Time `db:"tournament_date"`
	MatchesAmount  int       `db:"matches_amount"`
}

func newTournament(t dtos.CreateTournamentRequest) Tournament {
	// "02-01-2006" is just format string
	tournamentDate, _ := time.Parse("02-01-2006 15:04:05", t.TournamentDate)

	return Tournament{
		Name:           t.Name,
		TournamentDate: tournamentDate,
		MatchesAmount:  t.MatchesAmount,
	}
}

func (r *TournamentRepository) SelectById(ctx context.Context, tournamentId int) (*dtos.GetTournamentByIdResponse, error) {
	query := `
		SELECT name, tournament_date, matches_amount
		FROM tournaments WHERE id = :id
	`

	params := map[string]any{"id": tournamentId}

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, errors.New("failed to execute query: " + err.Error())
	}
	defer rows.Close()

	// Check if there are any rows
	if !rows.Next() {
		return nil, errors.New("no tournament with such id")
	}

	// Process row
	var tournament Tournament
	err = rows.StructScan(&tournament)
	if err != nil {
		return nil, errors.New("failed to get tournament: " + err.Error())
	}

	response := &dtos.GetTournamentByIdResponse{
		Name:           tournament.Name,
		TournamentDate: tournament.TournamentDate.Format("02-01-2006 15:04:05"),
		MatchesAmount:  tournament.MatchesAmount,
	}

	return response, nil
}

func (r *TournamentRepository) Insert(ctx context.Context, req dtos.CreateTournamentRequest) error {
	tournament := newTournament(req)

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO tournaments (
			name,
			tournament_date,
			matches_amount
	    ) VALUES (:name, :tournament_date, :matches_amount)
	`

	_, err = tx.NamedExecContext(ctx, query, tournament)
	if err != nil {
		return errors.New("failed to insert tournament: " + err.Error())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *TournamentRepository) UpdateById(
	ctx context.Context, req dtos.CreateTournamentRequest, tournamentId int,
) error {
	tournament := newTournament(req)
	tournament.Id = tournamentId

	// Start transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE tournaments SET
			name = :name,
			tournament_date = :tournament_date,
			matches_amount = :matches_amount
	    WHERE id = :id
	`

	_, err = tx.NamedExecContext(ctx, query, tournament)
	if err != nil {
		return errors.New("failed to update tournament: " + err.Error())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}
