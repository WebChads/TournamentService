package dtos

import (
	"time"
)

type CreateTournamentRequest struct {
	Name           string `json:"name" validate:"required"`
	TournamentDate string `json:"tournament_date" validate:"required"`
	MatchesAmount  int    `json:"matches_amount" validate:"required"`
}

type GetTournamentResponse struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	TournamentDate time.Time `json:"tournament_date"`
	MatchesAmount  int       `json:"matches_amount"`
}
