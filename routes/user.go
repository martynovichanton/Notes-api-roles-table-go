package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"notes-api-go/db/database"

	"golang.org/x/crypto/bcrypt"
)

// var DB *sql.DB

type UserRoutes struct {
	Queries *database.Queries
	DB      *sql.DB
}

func (ur *UserRoutes) CreateUser(w http.ResponseWriter, r *http.Request) {
	// create user with transaction

	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	var createRequest struct {
		Username string   `json:"username"`
		Password string   `json:"password"`
		Roles    []string `json:"roles"`
		Active   bool     `json:"active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Default role is "user" if not provided
	// if user.Roles == "" {
	// 	user.Roles = "user"
	// }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	createRequest.Password = string(hashedPassword)

	user := database.CreateUserParams{
		Username: createRequest.Username,
		Password: createRequest.Password,
		Active:   createRequest.Active,
	}

	tx, err := ur.DB.Begin()
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	qtx := ur.Queries.WithTx(tx)

	userID, err := qtx.CreateUser(context.Background(), user)
	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		fmt.Println(err)
		return
	}

	// Assign roles to the user
	for _, roleName := range createRequest.Roles {
		roleID, err := qtx.GetRoleIDByName(r.Context(), roleName)
		if err != nil {
			http.Error(w, "Invalid role: "+roleName, http.StatusBadRequest)
			return
		}

		err = qtx.AddUserRole(r.Context(), database.AddUserRoleParams{
			UserID: userID,
			RoleID: roleID,
		})
		if err != nil {
			http.Error(w, "Failed to assign role: "+roleName, http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	w.WriteHeader(http.StatusCreated)
}

func (ur *UserRoutes) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := ur.Queries.GetUsers(context.Background())
	if err != nil {
		http.Error(w, "Could not fetch users", http.StatusInternalServerError)
		return
	}

	var usersWithRoles []map[string]interface{}
	for _, user := range users {
		roles, err := ur.Queries.GetRolesForUser(context.Background(), user.ID)
		if err != nil {
			http.Error(w, "User already exists", http.StatusConflict)
			fmt.Println(err)
			return
		}

		userWithRoles := map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"roles":    roles,
			"active":   user.Active,
		}
		usersWithRoles = append(usersWithRoles, userWithRoles)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersWithRoles)
}

func (ur *UserRoutes) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// update user with transaction

	// if r.Method != http.MethodPut {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	var updateRequest struct {
		ID       int64    `json:"id"`
		Username string   `json:"username"`
		Password string   `json:"password"`
		Roles    []string `json:"roles"`
		Active   bool     `json:"active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	updateRequest.Password = string(hashedPassword)

	user := database.UpdateUserParams{
		ID:       updateRequest.ID,
		Username: updateRequest.Username,
		Password: updateRequest.Password,
		Active:   updateRequest.Active,
	}

	tx, err := ur.DB.Begin()
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	qtx := ur.Queries.WithTx(tx)

	if err := qtx.UpdateUser(context.Background(), user); err != nil {
		http.Error(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	// Remove all roles from the uer
	if err := qtx.RemoveRolesForUser(context.Background(), user.ID); err != nil {
		http.Error(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	// Assign roles to the user
	for _, roleName := range updateRequest.Roles {
		roleID, err := qtx.GetRoleIDByName(r.Context(), roleName)
		if err != nil {
			http.Error(w, "Invalid role: "+roleName, http.StatusBadRequest)
			return
		}

		err = qtx.AddUserRole(r.Context(), database.AddUserRoleParams{
			UserID: updateRequest.ID,
			RoleID: roleID,
		})
		if err != nil {
			http.Error(w, "Failed to assign role: "+roleName, http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	w.WriteHeader(http.StatusOK)
}

func (ur *UserRoutes) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var deleteRequest struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := ur.Queries.DeleteUser(context.Background(), deleteRequest.ID); err != nil {
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
