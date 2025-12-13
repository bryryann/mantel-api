-- This migration creates the 'users' table if it does not already exist.
-- The table includes the following columns:
-- 1. id: An integer that is auto-generated as the primary key.
-- 2. created_at: A timestamp with timezone, defaulting to the current time.
-- 3. username: A unique text field that cannot be null.
-- 4. email: A unique case-insensitive text field that cannot be null.
-- 5. password_hash: A bytea field for storing hashed passwords, cannot be null.
-- 6. version: An integer field with a default value of 1.
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    username text UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    version integer NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX users_email_lower_unique
ON users (LOWER(email));
