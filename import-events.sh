#!/bin/bash

# Script para importar eventos desde Excel a PostgreSQL
# Usage: ./import-events.sh [excel_file]

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  CustomerMX - Import Events from Excel${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# Check if PostgreSQL is running
echo -e "${YELLOW}[1/6]${NC} Verificando PostgreSQL..."
if ! docker ps | grep -q postgres; then
    echo -e "${RED}❌ PostgreSQL no está corriendo${NC}"
    echo -e "${YELLOW}💡 Iniciando PostgreSQL...${NC}"
    docker-compose up -d postgres
    sleep 3
fi
echo -e "${GREEN}✅ PostgreSQL está corriendo${NC}"
echo ""

# Check if Excel file exists
EXCEL_FILE="${1:-eventos.xlsx}"
echo -e "${YELLOW}[2/6]${NC} Verificando archivo Excel..."
if [ ! -f "$EXCEL_FILE" ]; then
    echo -e "${RED}❌ Archivo no encontrado: $EXCEL_FILE${NC}"
    echo -e "${YELLOW}💡 Coloca tu archivo Excel en la raíz del proyecto${NC}"
    echo -e "${YELLOW}   Uso: ./import-events.sh [ruta_al_excel]${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Archivo encontrado: $EXCEL_FILE${NC}"
echo ""

# Check database connection
echo -e "${YELLOW}[3/6]${NC} Verificando conexión a base de datos..."
if ! docker exec customermx-postgres-1 psql -U postgres -d customermx -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}❌ No se puede conectar a la base de datos${NC}"
    echo -e "${YELLOW}💡 Verifica que PostgreSQL esté corriendo correctamente${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Conexión exitosa${NC}"
echo ""

# Create output directory
echo -e "${YELLOW}[4/6]${NC} Preparando directorio de salida..."
mkdir -p backend/migrations/import
echo -e "${GREEN}✅ Directorio creado${NC}"
echo ""

# Install dependencies if needed
echo -e "${YELLOW}[5/6]${NC} Verificando dependencias Go..."
cd backend/cmd/import-events
if ! go list -m github.com/xuri/excelize/v2 > /dev/null 2>&1; then
    echo -e "${YELLOW}💡 Instalando dependencias...${NC}"
    go get github.com/xuri/excelize/v2
    go get github.com/jackc/pgx/v5/pgxpool
    go get github.com/google/uuid
fi
echo -e "${GREEN}✅ Dependencias verificadas${NC}"
echo ""

# Run import
echo -e "${YELLOW}[6/6]${NC} Ejecutando importación..."
echo -e "${BLUE}================================================${NC}"
echo ""
go run main.go
echo ""
echo -e "${BLUE}================================================${NC}"

# Show summary
echo ""
echo -e "${GREEN}✅ ¡Importación completada!${NC}"
echo ""
echo -e "${YELLOW}📁 Archivos SQL generados en:${NC}"
echo -e "   backend/migrations/import/"
echo ""
echo -e "${YELLOW}🔍 Para verificar los datos importados:${NC}"
echo -e "   ${BLUE}docker exec -it customermx-postgres-1 psql -U postgres -d customermx${NC}"
echo -e "   ${BLUE}SELECT COUNT(*) FROM events;${NC}"
echo ""
