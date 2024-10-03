-- Create the token table
CREATE TABLE IF NOT EXISTS token (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE (user_id),  -- Ensure one token per user
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE  -- Foreign key reference to users table
);

-- Create an index for the user_id column for faster lookups
CREATE INDEX idx_token_user_id ON token(user_id);

-- Create or replace the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop the existing trigger if it exists
DROP TRIGGER IF EXISTS update_token_updated_at ON token;

-- Create the trigger to call the function before update
CREATE TRIGGER update_token_updated_at
BEFORE UPDATE ON token
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
