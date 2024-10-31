-- Create the article table
CREATE TABLE IF NOT EXISTS article (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,                  -- Maps to `Title` field in Go
    content TEXT NOT NULL,                        -- Maps to `Content` field in Go
    image_url VARCHAR(500),                       -- Maps to `ImageURL` field in Go
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create an index for the id column for faster lookups
CREATE INDEX idx_article_id ON article(id);

-- Create or replace the function to update the updated_at column
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop the existing trigger if it exists
DROP TRIGGER IF EXISTS update_article_updated_at ON article;

-- Create the trigger to call the function before update
CREATE TRIGGER update_article_updated_at
BEFORE UPDATE ON article
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- Insert sample data into the article table with content about the Santapan app
INSERT INTO article (title, content, image_url, created_at, updated_at)
VALUES 
    ('Welcome to Santapan: Your Guide to Healthy Eating', 
     'Santapan is designed to help users make informed choices about their diets. Discover how this app can support your journey toward balanced and nutritious meals.', 
     'https://example.com/images/santapan-intro.jpg', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),
     
    ('5 Reasons to Use Santapan for Meal Planning', 
     'Meal planning is crucial for maintaining a healthy diet. Santapan offers customizable meal plans, nutrition tracking, and recipe suggestions to simplify the process.', 
     'https://example.com/images/meal-planning.jpg', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),
     
    ('Understanding Nutritional Labels with Santapan', 
     'Deciphering nutritional labels can be challenging. Santapan provides an easy guide to understand what’s in your food, helping you make healthier choices.', 
     'https://example.com/images/nutritional-labels.jpg', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),

    ('How Santapan Supports Local Farmers and Fresh Ingredients', 
     'Santapan partners with local farmers to promote fresh and organic ingredients. Learn how the app connects users with nearby sources for healthier options.', 
     'https://example.com/images/fresh-ingredients.jpg', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP),

    ('Personalized Nutrition Insights with Santapan', 
     'Everyone’s nutritional needs are different. Santapan’s personalized insights help you understand the best dietary choices based on your health goals and preferences.', 
     'https://example.com/images/personalized-nutrition.jpg', 
     CURRENT_TIMESTAMP, 
     CURRENT_TIMESTAMP);
