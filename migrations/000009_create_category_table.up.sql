-- Create the banner table
CREATE TABLE IF NOT EXISTS banner (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,               -- Maps to `Title` field in Go
    image_url VARCHAR(255),                     -- Maps to `ImageURL` field in Go
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create an index for the id column for faster lookups
CREATE INDEX idx_banner_id ON banner(id);

-- Create or replace the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop the existing trigger if it exists
DROP TRIGGER IF EXISTS update_banner_updated_at ON banner;

-- Create the trigger to call the function before update
CREATE TRIGGER update_banner_updated_at
BEFORE UPDATE ON banner
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
