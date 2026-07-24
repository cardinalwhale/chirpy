package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGetOne(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		http.Error(w, "failed to parse UUID", http.StatusInternalServerError)
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		http.Error(w, "chirp not found", http.StatusNotFound)
	}
	var chirpRes Chirp
	chirpRes = Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, 200, chirpRes)
}
