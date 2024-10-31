CREATE TABLE menu (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    image_url VARCHAR(255),
    nutrition JSONB,       -- JSON format for flexible nutrition details
    features JSONB,        -- Features stored as JSON for flexibility
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bundling (
    id BIGSERIAL PRIMARY KEY,
    bundling_type VARCHAR(50) NOT NULL,  -- 'weekly' or 'monthly'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bundling_menu (
    id BIGSERIAL PRIMARY KEY,
    bundling_id INT REFERENCES bundling(id) ON DELETE CASCADE,
    menu_id INT REFERENCES menu(id) ON DELETE CASCADE,
    day_number INT NOT NULL,          -- Represents the day number (1-7 for weekly, 1-30 for monthly)
    meal_description TEXT,            -- Description of the meal for the day
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(bundling_id, menu_id, day_number)  -- Ensure unique combination of bundling and menu items per day
);

-- Trigger function to update `updated_at` timestamp
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to update `updated_at` column on row modifications
CREATE TRIGGER update_menu_timestamp
BEFORE UPDATE ON menu
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_bundling_timestamp
BEFORE UPDATE ON bundling
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_bundling_menu_timestamp
BEFORE UPDATE ON bundling_menu
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
