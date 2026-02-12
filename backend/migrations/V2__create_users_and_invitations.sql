-- V2: Create users and invitations tables
-- Description: User management, roles, and invitation system

-- User role enum
CREATE TYPE user_role AS ENUM ('ADMIN', 'COORDINATOR', 'BRAND');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL,
    brand_id UUID REFERENCES brands(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    -- Constraint: BRAND users must have a brand_id, others must not
    CHECK (
        (role = 'BRAND' AND brand_id IS NOT NULL)
        OR
        (role <> 'BRAND' AND brand_id IS NULL)
    )
);

-- Invitations table
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL,
    role user_role NOT NULL,
    brand_id UUID REFERENCES brands(id),
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    -- Constraint: BRAND invitations must have a brand_id
    CHECK (
        (role = 'BRAND' AND brand_id IS NOT NULL)
        OR
        (role <> 'BRAND' AND brand_id IS NULL)
    )
);

-- Indexes
CREATE INDEX idx_users_brand ON users(brand_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_invitations_email ON invitations(email);
CREATE INDEX idx_invitations_token ON invitations(token);

-- Comments for documentation
COMMENT ON TABLE users IS 'System users with role-based access control';
COMMENT ON COLUMN users.role IS 'User role: ADMIN (full access), COORDINATOR (event management), BRAND (view own brand data)';
COMMENT ON COLUMN users.brand_id IS 'Required for BRAND role users, NULL for ADMIN and COORDINATOR';
COMMENT ON TABLE invitations IS 'Pending user invitations sent by administrators';
COMMENT ON COLUMN invitations.token IS 'Unique token sent via email for accepting invitation';
COMMENT ON COLUMN invitations.expires_at IS 'Expiration timestamp for the invitation';
