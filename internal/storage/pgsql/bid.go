package storage

import (
	"context"
	"errors"

	"github.com/WebChads/TournamentService/internal/models/dtos"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BidRepository struct {
	db *sqlx.DB
}

func NewBidRepository(db *sqlx.DB) *BidRepository {
	return &BidRepository{
		db: db,
	}
}

type Bid struct {
	BidId        uuid.UUID `db:"bid_id"`
	FullName     string    `db:"full_name"`
	Age          int       `db:"age"`
	AvatarUrl    string    `db:"avatar_url"`
	BidStatus    int       `db:"bid_status"`
	UserId       uuid.UUID `db:"user_id"`
	TournamentId uuid.UUID `db:"tournament_id"`
}

func NewBid(b dtos.CreateBidRequest) Bid {
	return Bid{
		BidId: uuid.New(),
		FullName: b.FullName,
		Age: b.Age,
		AvatarUrl: b.AvatarUrl,
		BidStatus: 3,
		UserId: b.UserId,
		TournamentId: b.TournamentId,
	}
}

func (r *BidRepository) Insert(ctx context.Context, req dtos.CreateBidRequest) error {
	tournament := NewBid(req)

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
		INSERT INTO bids (
			bid_id,
			full_name,
			age,
			avatar_url,
			bid_status,
			user_id,
			tournament_id
	    ) VALUES (:bid_id, :full_name, :age, :avatar_url, :bid_status, :user_id, :tournament_id)
	`

	_, err = tx.NamedExecContext(ctx, query, tournament)
	if err != nil {
		return errors.New("failed to insert bid: " + err.Error())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}
