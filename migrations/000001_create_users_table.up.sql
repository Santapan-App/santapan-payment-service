CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,                          -- int64 in Go maps to BIGINT in SQL
    full_name VARCHAR(255) NOT NULL,                  -- Map to `FirstName` field in Go
    email VARCHAR(255) NOT NULL UNIQUE,                -- Map to `Email Address` field in Go, with UNIQUE constraint
    password VARCHAR(255) NOT NULL,                  -- Map to `Password` field in Go
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Maps to `CreatedAt`
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Maps to `UpdatedAt`
    deleted_at TIMESTAMP WITH TIME ZONE NULL,          -- Nullable for soft deletion, maps to `DeletedAt`
    email_verified_at TIMESTAMP WITH TIME ZONE         -- (Missing in Go struct, add it if needed)
);

-- Create an index for the email field
CREATE INDEX idx_users_email ON users (email);

-- Create trigger function to update `updated_at`
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger to call the function before updates
CREATE TRIGGER update_user_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
