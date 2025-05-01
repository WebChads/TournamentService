package dtos

type CreateTournamentRequest struct {
	Name           string `json:"name" validate:"required"`
	TournamentDate string `json:"tournament_date" validate:"required"`
	MatchesAmount  int    `json:"matches_amount" validate:"required"`
}

type GetTournamentByIdResponse struct {
	Name           string `json:"name"`
	TournamentDate string `json:"tournament_date"`
	MatchesAmount  int    `json:"matches_amount"`
}
