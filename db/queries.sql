-- User Queries
-- name: CreateUser :one
INSERT INTO users (username, password, active) VALUES ($1, $2, $3) RETURNING id;

-- name: GetUserByUsername :one
SELECT id, username, password, active FROM users WHERE username = $1;

-- name: GetUsers :many
SELECT id, username, active FROM users;

-- name: GetUsersWithRoles :many
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
    u.id, u.username;


-- name: UpdateUser :exec
UPDATE users SET username = $1, password = $2, active = $3 WHERE id = $4;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;




-- UserRoles Queries
-- name: GetRoleIDByName :one
SELECT id FROM roles WHERE name = $1;

-- name: GetRolesForUser :many
SELECT r.name
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: AddUserRole :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveUserRole :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2;

-- name: RemoveRolesForUser :exec
DELETE FROM user_roles WHERE user_id = $1;

-- name: CheckUserRole :one
SELECT EXISTS (
    SELECT 1
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    WHERE ur.user_id = $1 AND r.name = $2
);




-- Note Queries
-- name: CreateNote :exec
INSERT INTO notes (user_id, content) VALUES ($1, $2);

-- name: GetNotesByUserID :many
SELECT id, user_id, content, created_at, updated_at FROM notes WHERE user_id = $1;

-- name: GetNotesByUserIDWithUserNames :many
SELECT notes.id, notes.user_id, notes.content, users.username
FROM notes
INNER JOIN users ON users.id=notes.user_id WHERE user_id = $1;

-- name: GetNotes :many
SELECT id, user_id, content, created_at, updated_at FROM notes;

-- name: GetNotesWithUserNames :many
SELECT notes.id, notes.user_id, notes.content, users.username
FROM notes
INNER JOIN users ON users.id=notes.user_id;

-- name: UpdateNote :exec
UPDATE notes SET content = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND user_id = $3;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1 AND user_id = $2;
