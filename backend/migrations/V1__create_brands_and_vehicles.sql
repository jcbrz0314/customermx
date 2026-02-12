-- V1: Create brands and vehicles tables
-- Description: Core tables for automotive brands and their vehicle models

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Brands table
CREATE TABLE brands (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now()
);

-- Vehicles table
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    model_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(brand_id, model_name)
);

-- Indexes for vehicles
CREATE INDEX idx_vehicles_brand ON vehicles(brand_id);

-- Comments for documentation
COMMENT ON TABLE brands IS 'Automotive brands (VW, Nissan, Chevrolet, etc.)';
COMMENT ON TABLE vehicles IS 'Vehicle models associated with each brand';
COMMENT ON COLUMN vehicles.brand_id IS 'Foreign key to brands table';
COMMENT ON COLUMN vehicles.model_name IS 'Model name (e.g., Jetta, Tiguan, Aveo)';
