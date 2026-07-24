package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	dbChirps, _ := cfg.db.GetAllChirps(r.Context())

	var chirps []Chirp
	for _, dbc := range dbChirps {
		var c Chirp
		c.ID = dbc.ID
		c.CreatedAt = dbc.CreatedAt
		c.UpdatedAt = dbc.UpdatedAt
		c.Body = dbc.Body
		c.UserID = dbc.UserID
		chirps = append(chirps, c)
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		http.Error(w, "failed to encode chirps", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
