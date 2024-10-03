-- Create the OTP type ENUM
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'otp_type') THEN
        CREATE TYPE otp_type AS ENUM ('sms', 'email', 'whatsapp');
    END IF;
END;
$$;

-- Create the otp table
CREATE TABLE IF NOT EXISTS otp (
    id BIGSERIAL PRIMARY KEY,
    code CHAR(6) NOT NULL,  -- Change code to fixed length char(6)
    retry SMALLINT NOT NULL,
    type otp_type NOT NULL,  -- Change type to ENUM type
    user_id BIGINT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id),  -- Secondary key for user_id and device_id
    FOREIGN KEY (user_id) REFERENCES users(id)  -- Foreign key relationship
);

-- Create or replace the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop the existing trigger if it exists (for safety)
DROP TRIGGER IF EXISTS update_token_updated_at ON otp;

-- Create the trigger to call the function before update
CREATE TRIGGER update_token_updated_at
BEFORE UPDATE ON otp
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();