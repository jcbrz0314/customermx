-- V5: Seed admin user
-- Description: Create default admin user (admin@customermx.com / admin123)

INSERT INTO users (name, email, password_hash, role, brand_id, is_active)
SELECT 'Administrator', 'admin@customermx.com',
       '$2a$10$YpAeT1IhCf0mIvFkOAnATOMjTm7lUV4yvYPAwIeMfF8vUeg761lKC',
       'ADMIN', NULL, TRUE
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'admin@customermx.com'
);
