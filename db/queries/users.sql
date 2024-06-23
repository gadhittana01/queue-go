-- name: CreateUser :one
INSERT INTO "users"(name, identity_number, email, date_of_birth) VALUES
($1, $2, $3, $4) RETURNING *;