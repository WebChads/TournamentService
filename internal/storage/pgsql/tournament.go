package storage

import (
	"context"
	"errors"
	"time"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	"github.com/google/uuid"
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
	Id             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	TournamentDate time.Time `db:"tournament_date"`
	MatchesAmount  int       `db:"matches_amount"`
	UserId         uuid.UUID `db:"user_id"`
}

func newTournament(t dtos.CreateTournamentRequest) Tournament {
	// "02-01-2006" is just format string
	tournamentDate, _ := time.Parse("02-01-2006 15:04:05", t.TournamentDate)

	return Tournament{
		Id:             uuid.New(),
		Name:           t.Name,
		TournamentDate: tournamentDate,
		MatchesAmount:  t.MatchesAmount,
		UserId:         t.UserId,
	}
}

func (r *TournamentRepository) SelectById(ctx context.Context, tournamentId uuid.UUID) (*dtos.GetTournamentByIdResponse, error) {
	query := `
		SELECT name, tournament_date, matches_amount, user_id
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
		UserId:         tournament.UserId,
	}

	return response, nil
}

func (r *TournamentRepository) SelectByName(
	ctx context.Context, tournamentName string,
) ([]dtos.GetTournamentByIdResponse, error) {
	query := `
		SELECT name, tournament_date, matches_amount, user_id
		FROM tournaments WHERE LOWER(name) LIKE $1
	`

	var tournaments []Tournament
	err := r.db.SelectContext(ctx, &tournaments, query, "%"+tournamentName+"%")
	if err != nil {
		return nil, errors.New("failed to select tournaments: " + err.Error())
	}

	var responseTournaments []dtos.GetTournamentByIdResponse
	for _, tournament := range tournaments {
		response := dtos.GetTournamentByIdResponse{
			Name:           tournament.Name,
			TournamentDate: tournament.TournamentDate.Format("02-01-2006 15:04:05"),
			MatchesAmount:  tournament.MatchesAmount,
			UserId:         tournament.UserId,
		}

		responseTournaments = append(responseTournaments, response)
	}

	return responseTournaments, nil
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
			id,
			name,
			tournament_date,
			matches_amount,
			user_id
	    ) VALUES (:id, :name, :tournament_date, :matches_amount, :user_id)
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
	ctx context.Context, req dtos.CreateTournamentRequest, tournamentId uuid.UUID,
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
			matches_amount = :matches_amount,
			user_id = :user_id
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

func (r *TournamentRepository) DeleteById(ctx context.Context, tournamentId uuid.UUID) error {
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

	query := `DELETE FROM tournaments WHERE id = :id`

	params := map[string]any{"id": tournamentId}

	_, err = tx.NamedExecContext(ctx, query, params)
	if err != nil {
		return errors.New("failed to delete tournament: " + err.Error())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}
