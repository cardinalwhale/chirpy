package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cardinalwhale/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get the user", err)
		return
	}

	check, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem with checking the hash", err)
		return
	}

	if !check {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds >= 3600 {
		params.ExpiresInSeconds = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(params.ExpiresInSeconds)*time.Second)

	respondWithJSON(w, http.StatusOK,
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	)
}
