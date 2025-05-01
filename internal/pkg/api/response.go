package response

import (
	"encoding/json"
	"net/http"

	"github.com/WebChads/TournamentService/internal/models/dtos"
)

func JSON(w http.ResponseWriter, statusCode int, message any) {
	response := dtos.Response{
		Status:  statusCode,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
