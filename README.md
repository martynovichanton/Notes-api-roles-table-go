# Notes API with Go

## Complete API with user authentication and roles with JWT

## Features

1. Goose for DB migrations.
2. SQLC for generating type safe DB operations from raw SQL queries.
3. Using the Go standard library net/http.
4. Postgresql DB.
5. JWT for user auth with access and refresh tokens.
6. Users can have multiple roles - admin or user or both.
7. Protected routes for admin only access.
8. User roles are defined in a separate table.
9. Middlwares for admin access and custom logger.
10. Preconfigured methods in rest-client folder.

## Setup
1. Set the .env and db/.env variables.
2. Set .vscode/settings.json variables for the rest client.
3. cd db
4. goose -dir migrations up
5. sqlc generate 
6. go run main.go 
7. DB cleanup - BE CAREFUL!!! - goose -dir migrations down-to 0
8. To create the first admin user - uncomment the route "POST /users" without the admin middleware in main.go and comment the line with the admin middleware.
9. Once the first admin user is created - comment the line without the middleware and uncomment the line with the middleware to protect the route.



