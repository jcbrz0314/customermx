-- V6: Seed test users for each role
-- Description: Create a COORDINATOR and a BRAND user for testing
-- All passwords: admin123

-- Coordinator user
INSERT INTO users (name, email, password_hash, role, brand_id, is_active)
SELECT 'Coordinador Test', 'coordinator@customermx.com',
       '$2a$10$YpAeT1IhCf0mIvFkOAnATOMjTm7lUV4yvYPAwIeMfF8vUeg761lKC',
       'COORDINATOR', NULL, TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'coordinator@customermx.com'
);

-- Brand user (Chevrolet)
INSERT INTO users (name, email, password_hash, role, brand_id, is_active)
SELECT 'Brand Chevrolet', 'chevrolet@customermx.com',
       '$2a$10$YpAeT1IhCf0mIvFkOAnATOMjTm7lUV4yvYPAwIeMfF8vUeg761lKC',
       'BRAND', (SELECT id FROM brands WHERE name = 'Chevrolet'), TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'chevrolet@customermx.com'
);
