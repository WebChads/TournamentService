package dtos

import "github.com/google/uuid"

type CreateTournamentRequest struct {
	Name           string    `json:"name" validate:"required"`
	TournamentDate string    `json:"tournament_date" validate:"required"`
	MatchesAmount  int       `json:"matches_amount" validate:"required"`
	UserId         uuid.UUID `json:"user_id"`
}

type GetTournamentResponse struct {
	Name           string    `json:"name"`
	TournamentDate string    `json:"tournament_date"`
	MatchesAmount  int       `json:"matches_amount"`
	UserId         uuid.UUID `json:"user_id"`
}
