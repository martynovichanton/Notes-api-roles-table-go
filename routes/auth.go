package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"notes-api-go/db/database"
	"notes-api-go/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthRoutes struct {
	Queries *database.Queries
}

func (ar *AuthRoutes) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	user, err := ar.Queries.GetUserByUsername(context.Background(), loginRequest.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !user.Active {
		http.Error(w, "User is inactive", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	roles, err := ar.Queries.GetRolesForUser(context.Background(), user.ID)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		fmt.Println(err)
		return
	}

	// databaseerate JWT tokens
	accessToken, err := utils.CreateToken(user.ID, user.Username, roles, time.Minute*60)
	if err != nil {
		http.Error(w, "Could not create access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.CreateToken(user.ID, user.Username, roles, time.Hour*24)
	if err != nil {
		http.Error(w, "Could not create refresh token", http.StatusInternalServerError)
		return
	}

	// Return tokens to the client
	response := map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ar *AuthRoutes) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var refreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Parse the refresh token
	claims, err := utils.VerifyToken(refreshRequest.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Fetch user from database
	user, err := ar.Queries.GetUserByUsername(context.Background(), claims.UserInfo.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !user.Active {
		http.Error(w, "User is inactive", http.StatusUnauthorized)
		return
	}

	// Generate a new access token
	accessToken, err := utils.CreateToken(claims.UserInfo.UserID, claims.UserInfo.Username, claims.UserInfo.Roles, time.Hour*24)
	if err != nil {
		http.Error(w, "Could not generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.CreateToken(claims.UserInfo.UserID, claims.UserInfo.Username, claims.UserInfo.Roles, time.Hour*24)
	if err != nil {
		http.Error(w, "Could not generate refresh token", http.StatusInternalServerError)
		return
	}

	// Return the new access token
	response := map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
