CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

CREATE TABLE IF NOT EXISTS "queue" (
  "id" UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  "queue_number" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "arrival_time" timestamptz NOT NULL,
  "service_time" timestamptz,
  "total_waiting_time" interval,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW())
);

ALTER TABLE "queue" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;