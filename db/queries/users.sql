-- name: CreateUser :one
INSERT INTO "users"(name, email, password) VALUES
($1, $2, $3) RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM "users" WHERE email=$1;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT id FROM "users" WHERE email=$1);

-- name: CheckUserExists :one
SELECT EXISTS(SELECT id FROM "users" WHERE id=$1);