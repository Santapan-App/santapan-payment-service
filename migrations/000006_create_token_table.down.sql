-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS update_token_updated_at ON token;

-- Drop the function if it exists
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop the table if it exists
DROP TABLE IF EXISTS token;
