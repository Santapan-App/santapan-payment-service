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

-- Seed data for menu table
INSERT INTO menu (title, description, price, image_url, nutrition, features)
VALUES
    ('Avocado Salad', 'A healthy avocado salad with fresh vegetables and a side of toast', 45.900, 'http://example.com/avocado_salad.jpg', '{"calories": 350, "protein": "5g", "fat": "30g"}', '{"gluten_free": true, "vegetarian": true}'),
    ('Grilled Chicken Salad', 'A delicious grilled chicken salad', 50.000, 'http://example.com/grilled_chicken_salad.jpg', '{"calories": 400, "protein": "35g"}', '{"high_protein": true}');

-- Seed data for bundling table with a weekly plan
INSERT INTO bundling (bundling_type)
VALUES
    ('weekly', 'monthly');

-- Seed data for bundling_menu table to represent meals for each day of the week
-- Here we assume bundling_id of 1 corresponds to the 'weekly' bundling in the `bundling` table
INSERT INTO bundling_menu (bundling_id, menu_id, day_number, meal_description)
VALUES
    (1, 1, 1, 'Lunch - Avocado Salad'),
    (1, 1, 1, 'Dinner - Avocado Salad'), -- Sample duplicate to match the design
    (1, 1, 2, 'Lunch - Avocado Salad'),
    (1, 1, 2, 'Dinner - Avocado Salad'),
    (1, 1, 3, 'Lunch - Avocado Salad'),
    (1, 2, 4, 'Lunch - Grilled Chicken Salad'),
    (1, 2, 5, 'Lunch - Grilled Chicken Salad'),
    (1, 2, 6, 'Lunch - Grilled Chicken Salad'),
    (1, 1, 7, 'Lunch - Avocado Salad');

-- Add more days as needed to fill the weekly plan

-- Verify the inserted data
SELECT * FROM menu;
SELECT * FROM bundling;
SELECT * FROM bundling_menu;
