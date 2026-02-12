-- V4: Seed initial brands and vehicles
-- Description: Insert initial automotive brands and their vehicle models

-- Insert brands
INSERT INTO brands (name) VALUES
    ('Chevrolet'),
    ('Buick'),
    ('GMC'),
    ('Cadillac');

-- Insert Chevrolet vehicles
INSERT INTO vehicles (brand_id, model_name)
SELECT id, model FROM brands, (VALUES
    ('Aveo'),
    ('Onix'),
    ('Spark EV'),
    ('Groove'),
    ('Captiva'),
    ('Tracker'),
    ('Trax'),
    ('Equinox EV'),
    ('Blazer EV'),
    ('Traverse'),
    ('Blazer'),
    ('Tahoe'),
    ('Suburban'),
    ('Camaro'),
    ('Corvette'),
    ('S10'),
    ('Tornado Van'),
    ('Montana'),
    ('Silverado'),
    ('Colorado'),
    ('Cheyenne'),
    ('Brightdrop')
) AS models(model)
WHERE brands.name = 'Chevrolet';

-- Insert Buick vehicles
INSERT INTO vehicles (brand_id, model_name)
SELECT id, model FROM brands, (VALUES
    ('Envista'),
    ('Encore'),
    ('Envision'),
    ('Enclave')
) AS models(model)
WHERE brands.name = 'Buick';

-- Insert GMC vehicles
INSERT INTO vehicles (brand_id, model_name)
SELECT id, model FROM brands, (VALUES
    ('Terrain'),
    ('Acadia'),
    ('Yukon'),
    ('Canyon'),
    ('Sierra'),
    ('Hummer')
) AS models(model)
WHERE brands.name = 'GMC';

-- Insert Cadillac vehicles
INSERT INTO vehicles (brand_id, model_name)
SELECT id, model FROM brands, (VALUES
    ('XT4'),
    ('XT5'),
    ('Optiq'),
    ('Lyriq'),
    ('Escalade'),
    ('Escalade V'),
    ('Escalade IQ')
) AS models(model)
WHERE brands.name = 'Cadillac';

-- Verification comment
-- This seed data includes:
-- - Chevrolet: 22 vehicle models
-- - Buick: 4 vehicle models
-- - GMC: 6 vehicle models
-- - Cadillac: 7 vehicle models
-- Total: 4 brands, 39 vehicle models
