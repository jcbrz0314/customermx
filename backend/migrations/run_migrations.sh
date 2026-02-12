#!/bin/bash

# Database Migration Script for CustomerMX
# This script runs all database migrations in order

set -e  # Exit on error

# Database configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-customermx}"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}CustomerMX Database Migration${NC}"
echo "================================"
echo "Host: $DB_HOST:$DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo ""

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: psql is not installed${NC}"
    echo "Please install PostgreSQL client tools"
    exit 1
fi

# Test database connection
echo -e "${YELLOW}Testing database connection...${NC}"
export PGPASSWORD=$DB_PASSWORD
if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c '\q' 2>/dev/null; then
    echo -e "${RED}Error: Cannot connect to database${NC}"
    echo "Make sure PostgreSQL is running (try: make docker-up)"
    exit 1
fi
echo -e "${GREEN}✓ Database connection successful${NC}"
echo ""

# Create database if it doesn't exist
echo -e "${YELLOW}Checking if database exists...${NC}"
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo -e "${GREEN}✓ Database '$DB_NAME' exists${NC}"
else
    echo -e "${YELLOW}Creating database '$DB_NAME'...${NC}"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "CREATE DATABASE $DB_NAME;"
    echo -e "${GREEN}✓ Database created${NC}"
fi
echo ""

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Run migrations in order
MIGRATIONS=(
    "V1__create_brands_and_vehicles.sql"
    "V2__create_users_and_invitations.sql"
    "V3__create_events_and_reports.sql"
    "V4__seed_brands_and_vehicles.sql"
)

echo -e "${YELLOW}Running migrations...${NC}"
echo ""

for migration in "${MIGRATIONS[@]}"; do
    migration_file="$SCRIPT_DIR/$migration"

    if [ ! -f "$migration_file" ]; then
        echo -e "${RED}Error: Migration file not found: $migration${NC}"
        exit 1
    fi

    echo -e "${YELLOW}Running: $migration${NC}"

    if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration_file" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ $migration completed${NC}"
    else
        echo -e "${RED}✗ $migration failed${NC}"
        echo "Try running it manually to see the error:"
        echo "psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migration_file"
        exit 1
    fi
    echo ""
done

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}All migrations completed successfully!${NC}"
echo ""

# Show summary
echo -e "${YELLOW}Database Summary:${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
SELECT
    'Brands' as table_name,
    COUNT(*) as count
FROM brands
UNION ALL
SELECT
    'Vehicles',
    COUNT(*)
FROM vehicles
UNION ALL
SELECT
    'Users',
    COUNT(*)
FROM users
UNION ALL
SELECT
    'Events',
    COUNT(*)
FROM events;
"

echo ""
echo -e "${GREEN}Ready to go! 🚀${NC}"
