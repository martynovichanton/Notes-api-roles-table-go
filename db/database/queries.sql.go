// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package database

import (
	"context"
	"database/sql"
)

const addUserRole = `-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`

type AddUserRoleParams struct {
	UserID int64
	RoleID int64
}

func (q *Queries) AddUserRole(ctx context.Context, arg AddUserRoleParams) error {
	_, err := q.db.ExecContext(ctx, addUserRole, arg.UserID, arg.RoleID)
	return err
}

const checkUserRole = `-- name: CheckUserRole :one
SELECT EXISTS (
    SELECT 1
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE ur.user_id = $1 AND r.name = $2
)
`

type CheckUserRoleParams struct {
	UserID int64
	Name   string
}

func (q *Queries) CheckUserRole(ctx context.Context, arg CheckUserRoleParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkUserRole, arg.UserID, arg.Name)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createNote = `-- name: CreateNote :exec
INSERT INTO notes (user_id, content) VALUES ($1, $2)
`

type CreateNoteParams struct {
	UserID  int64
	Content string
}

// Note Queries
func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) error {
	_, err := q.db.ExecContext(ctx, createNote, arg.UserID, arg.Content)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password, active) VALUES ($1, $2, $3) RETURNING id
`

type CreateUserParams struct {
	Username string
	Password string
	Active   bool
}

// User Queries
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password, arg.Active)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteNote = `-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1 AND user_id = $2
`

type DeleteNoteParams struct {
	ID     int64
	UserID int64
}

func (q *Queries) DeleteNote(ctx context.Context, arg DeleteNoteParams) error {
	_, err := q.db.ExecContext(ctx, deleteNote, arg.ID, arg.UserID)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getNotes = `-- name: GetNotes :many
SELECT id, user_id, content, created_at, updated_at FROM notes
`

func (q *Queries) GetNotes(ctx context.Context) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, getNotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotesByUserID = `-- name: GetNotesByUserID :many
SELECT id, user_id, content, created_at, updated_at FROM notes WHERE user_id = $1
`

func (q *Queries) GetNotesByUserID(ctx context.Context, userID int64) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, getNotesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotesByUserIDWithUserNames = `-- name: GetNotesByUserIDWithUserNames :many
SELECT notes.id, notes.user_id, notes.content, users.username
FROM notes
INNER JOIN users ON users.id=notes.user_id WHERE user_id = $1
`

type GetNotesByUserIDWithUserNamesRow struct {
	ID       int64
	UserID   int64
	Content  string
	Username string
}

func (q *Queries) GetNotesByUserIDWithUserNames(ctx context.Context, userID int64) ([]GetNotesByUserIDWithUserNamesRow, error) {
	rows, err := q.db.QueryContext(ctx, getNotesByUserIDWithUserNames, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNotesByUserIDWithUserNamesRow
	for rows.Next() {
		var i GetNotesByUserIDWithUserNamesRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotesWithUserNames = `-- name: GetNotesWithUserNames :many
SELECT notes.id, notes.user_id, notes.content, users.username
FROM notes
INNER JOIN users ON users.id=notes.user_id
`

type GetNotesWithUserNamesRow struct {
	ID       int64
	UserID   int64
	Content  string
	Username string
}

func (q *Queries) GetNotesWithUserNames(ctx context.Context) ([]GetNotesWithUserNamesRow, error) {
	rows, err := q.db.QueryContext(ctx, getNotesWithUserNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNotesWithUserNamesRow
	for rows.Next() {
		var i GetNotesWithUserNamesRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Content,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoleIDByName = `-- name: GetRoleIDByName :one
SELECT id FROM roles WHERE name = $1
`

// UserRoles Queries
func (q *Queries) GetRoleIDByName(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getRoleIDByName, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getRolesForUser = `-- name: GetRolesForUser :many
SELECT r.name
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1
`

func (q *Queries) GetRolesForUser(ctx context.Context, userID int64) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getRolesForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, password, active FROM users WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Active,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, username, active FROM users
`

type GetUsersRow struct {
	ID       int64
	Username string
	Active   bool
}

func (q *Queries) GetUsers(ctx context.Context) ([]GetUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersRow
	for rows.Next() {
		var i GetUsersRow
		if err := rows.Scan(&i.ID, &i.Username, &i.Active); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersWithRoles = `-- name: GetUsersWithRoles :many
SELECT 
    u.id AS user_id, 
    u.username, 
    COALESCE(string_agg(r.name, ','), '') AS roles
FROM 
    users u
LEFT JOIN 
    user_roles ur ON u.id = ur.user_id
LEFT JOIN 
    roles r ON ur.role_id = r.id
GROUP BY 
    u.id, u.username
`

type GetUsersWithRolesRow struct {
	UserID   int64
	Username string
	Roles    sql.NullString
}

func (q *Queries) GetUsersWithRoles(ctx context.Context) ([]GetUsersWithRolesRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersWithRoles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersWithRolesRow
	for rows.Next() {
		var i GetUsersWithRolesRow
		if err := rows.Scan(&i.UserID, &i.Username, &i.Roles); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeRolesForUser = `-- name: RemoveRolesForUser :exec
DELETE FROM user_roles WHERE user_id = $1
`

func (q *Queries) RemoveRolesForUser(ctx context.Context, userID int64) error {
	_, err := q.db.ExecContext(ctx, removeRolesForUser, userID)
	return err
}

const removeUserRole = `-- name: RemoveUserRole :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2
`

type RemoveUserRoleParams struct {
	UserID int64
	RoleID int64
}

func (q *Queries) RemoveUserRole(ctx context.Context, arg RemoveUserRoleParams) error {
	_, err := q.db.ExecContext(ctx, removeUserRole, arg.UserID, arg.RoleID)
	return err
}

const updateNote = `-- name: UpdateNote :exec
UPDATE notes SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND user_id = $3
`

type UpdateNoteParams struct {
	Content string
	ID      int64
	UserID  int64
}

func (q *Queries) UpdateNote(ctx context.Context, arg UpdateNoteParams) error {
	_, err := q.db.ExecContext(ctx, updateNote, arg.Content, arg.ID, arg.UserID)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET username = $1, password = $2, active = $3 WHERE id = $4
`

type UpdateUserParams struct {
	Username string
	Password string
	Active   bool
	ID       int64
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Username,
		arg.Password,
		arg.Active,
		arg.ID,
	)
	return err
}
