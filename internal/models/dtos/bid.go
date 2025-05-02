package dtos

import "github.com/google/uuid"

type CreateBidRequest struct {
	AvatarUrl    string    `json:"avatar_url"`
	Age          int       `json:"age" validate:"required"`
	FullName     string    `json:"full_name" validate:"required"`
	UserId       uuid.UUID `json:"user_id"`
	TournamentId uuid.UUID `json:"tournament_id"`
}
