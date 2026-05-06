-- V8: Add VISUALIZER role to user_role enum
-- Description: Adds the VISUALIZER role that allows read-only access to the platform

ALTER TYPE user_role ADD VALUE IF NOT EXISTS 'VISUALIZER';
