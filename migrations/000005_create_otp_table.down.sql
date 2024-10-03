-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS update_token_updated_at ON otp;

-- Drop the function if it exists
DROP FUNCTION IF EXISTS update_timestamp();

-- Drop the otp table if it exists
DROP TABLE IF EXISTS otp;

-- Drop the ENUM type if it exists
DROP TYPE IF EXISTS otp_type;
