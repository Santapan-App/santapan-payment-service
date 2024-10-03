-- Insert seed data into the users table
INSERT INTO users (full_name, email, password, created_at, updated_at, email_verified_at)
VALUES
('Pahala', 'pfnazhmi@gmail.com', '$2a$12$QxFnM4e4yo7oQeb.olGL5.ntXxa9Qcaa.Kn1RZZHgNWA3EJSULQDS', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('Jane Smith', 'jane.smith@example.com', 'hashed_password_2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
('Robert Brown', 'robert.brown@example.com', 'hashed_password_3', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
