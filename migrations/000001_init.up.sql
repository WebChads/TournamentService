-- Up migration: creates the tournaments table
CREATE TABLE IF NOT EXISTS tournaments(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tournament_date TIMESTAMP NOT NULL,
    matches_amount INTEGER NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on name for better query performance
CREATE INDEX idx_tournaments_name ON tournaments(name);