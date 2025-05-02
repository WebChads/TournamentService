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
		BidId:        uuid.New(),
		FullName:     b.FullName,
		Age:          b.Age,
		AvatarUrl:    b.AvatarUrl,
		BidStatus:    3,
		UserId:       b.UserId,
		TournamentId: b.TournamentId,
	}
}

type UpdateBid struct {
	BidId     uuid.UUID `db:"bid_id"`
	BidStatus int       `db:"bid_status"`
}

func (r *BidRepository) SelectByStatus(
	ctx context.Context, req dtos.GetBidRequest,
) ([]dtos.GetBidResponse, error) {
	query := `
		SELECT bid_id, full_name, age, avatar_url, bid_status
		FROM bids WHERE user_id = $1 AND bid_status = $2
	`

	var bids []Bid
	err := r.db.SelectContext(ctx, &bids, query, req.UserId, req.BidStatus)
	if err != nil {
		return nil, errors.New("failed to select bids: " + err.Error())
	}

	var responseBids []dtos.GetBidResponse
	for _, bid := range bids {
		response := dtos.GetBidResponse{
			BidId:     bid.BidId,
			FullName:  bid.FullName,
			Age:       bid.Age,
			AvatarUrl: bid.AvatarUrl,
			BidStatus: bid.BidStatus,
		}

		responseBids = append(responseBids, response)
	}

	return responseBids, nil
}

func (r *BidRepository) Insert(ctx context.Context, req dtos.CreateBidRequest) error {
	bid := NewBid(req)

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

	_, err = tx.NamedExecContext(ctx, query, bid)
	if err != nil {
		return errors.New("failed to insert bid: " + err.Error())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *BidRepository) UpdateStatus(
	ctx context.Context, dto dtos.UpdateBidStatusRequest,
) error {
	bid := UpdateBid{
		BidId: dto.BidId,
		BidStatus: dto.BidStatus,
	}

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
		UPDATE bids SET
			bid_status = :bid_status
		WHERE bid_id = :bid_id
	`
	_, err = tx.NamedExecContext(ctx, query, bid)
	if err != nil {
		return errors.New("failed to update bid status: " + err.Error())
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}
