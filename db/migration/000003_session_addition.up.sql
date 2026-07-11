CREATE TABLE IF NOT EXISTS "sessions" (
    "id" uuid PRIMARY KEY,
    "username" varchar NOT NULL REFERENCES users(username) ON DELETE CASCADE,
    "user_agent" TEXT NOT NULL,
    "client_ip" TEXT NOT NULL,
    "refresh_token" TEXT NOT NULL,
    "is_blocked" BOOLEAN NOT NULL DEFAULT FALSE,
    "expires_at" TIMESTAMP NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);