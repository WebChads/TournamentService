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

-- Up migration: creates the bids table
CREATE TABLE IF NOT EXISTS bids(
    bid_id UUID PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    age INTEGER NOT NULL,
    avatar_url VARCHAR(512),
    bid_status VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    tournament_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on full name for better query performance
CREATE INDEX idx_bids_full_name ON bids(full_name);
