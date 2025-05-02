package dtos

import "github.com/google/uuid"

type CreateBidRequest struct {
	AvatarUrl    string    `json:"avatar_url"`
	Age          int       `json:"age" validate:"required"`
	FullName     string    `json:"full_name" validate:"required"`
	UserId       uuid.UUID `json:"user_id"`
	TournamentId uuid.UUID `json:"tournament_id"`
}

type GetBidRequest struct {
	BidStatus int       `json:"bid_status"`
	UserId    uuid.UUID `json:"user_id"`
}

type GetBidResponse struct {
	FullName  string `json:"full_name"`
	Age       int    `json:"age"`
	AvatarUrl string `json:"avatar_url"`
	BidStatus int    `json:"bid_status"`
}
