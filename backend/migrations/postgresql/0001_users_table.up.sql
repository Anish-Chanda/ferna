-- Enum for auth provider
CREATE TYPE user_auth_provider AS ENUM ('local', 'google', 'facebook');

-- Users table with check for password_hash when auth_provider is 'local'
CREATE TABLE users (
    id uuid PRIMARY KEY,
    avatar_url text,
    auth_provider user_auth_provider NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    email varchar(254) NOT NULL, -- Max length for email as per RFC 5321
    full_name text NOT NULL,
    password_hash text,
    timezone text NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT users_email_unique UNIQUE (email),
    -- If auth_provider = 'local', password_hash must be set
    CONSTRAINT users_local_auth_requires_password CHECK (
        auth_provider <> 'local'
        OR password_hash IS NOT NULL
    )
);