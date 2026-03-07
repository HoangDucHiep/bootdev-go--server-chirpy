package main

import (
	"net/http"
	"time"

	"ChirpyBlog/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token")
		return
	}

	// Look up refresh token in database
	tokenRecord, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Check if token is expired
	if tokenRecord.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	// Check if token is revoked
	if tokenRecord.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked")
		return
	}

	// Create new access token (1 hour expiration)
	accessToken, err := auth.MakeJWT(tokenRecord.UserID, cfg.jwtSecretKey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token")
		return
	}

	// Revoke the refresh token
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
