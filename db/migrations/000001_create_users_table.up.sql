SET timezone TO 'GMT';

CREATE TABLE IF NOT EXISTS users(
  "id" UUID NOT NULL DEFAULT gen_random_uuid(),
  "name" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "password" VARCHAR NOT NULL,
  "verified" BOOLEAN NOT NULL DEFAULT false,
  "active" BOOLEAN NOT NULL DEFAULT true,
  "verification_token" VARCHAR,
  "avatar" VARCHAR NOT NULL DEFAULT 'user_avatar0.png',
  "created_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  CONSTRAINT "uq_users_email" UNIQUE ("email"),
  CONSTRAINT "pk_users_id" PRIMARY KEY ("id")
);
