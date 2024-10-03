DROP TABLE IF EXISTS users;

-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS update_user_updated_at ON users;

-- Drop the function if it exists
DROP FUNCTION IF EXISTS update_timestamp();
