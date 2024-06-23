CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "identity_number" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "date_of_birth" DATE NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);