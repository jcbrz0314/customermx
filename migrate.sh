#!/bin/bash

# migrate.sh - Corre migraciones de base de datos via Docker
# Uso: ./migrate.sh

set -e

CONTAINER="customermx-postgres"
DB_USER="customermx"
DB_NAME="customermx"
MIGRATIONS_DIR="backend/migrations"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}CustomerMX - Migraciones${NC}"
echo "========================"
echo ""

# Verificar que el contenedor esta corriendo
if [ ! "$(docker ps -q -f name=$CONTAINER)" ]; then
  echo -e "${RED}Error: El contenedor '$CONTAINER' no esta corriendo.${NC}"
  echo "Ejecuta ./runenv.sh primero."
  exit 1
fi

# Crear tabla de control de migraciones si no existe
docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c "
CREATE TABLE IF NOT EXISTS schema_migrations (
    filename TEXT PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT now()
);" > /dev/null 2>&1

echo -e "${GREEN}✓ Conectado a la base de datos${NC}"
echo ""

# Copiar migraciones al contenedor
docker cp "$MIGRATIONS_DIR" "$CONTAINER:/tmp/migrations" > /dev/null 2>&1

# Buscar archivos V*__.sql y ordenarlos
MIGRATIONS=$(ls "$MIGRATIONS_DIR"/V*.sql 2>/dev/null | sort -V)

if [ -z "$MIGRATIONS" ]; then
  echo -e "${YELLOW}No hay migraciones para ejecutar.${NC}"
  exit 0
fi

APPLIED=0
SKIPPED=0

for migration_path in $MIGRATIONS; do
  filename=$(basename "$migration_path")

  # Verificar si ya fue aplicada
  already_applied=$(docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -tAc \
    "SELECT COUNT(*) FROM schema_migrations WHERE filename = '$filename';")

  if [ "$already_applied" -gt 0 ]; then
    echo -e "  ${GREEN}✓${NC} $filename (ya aplicada)"
    SKIPPED=$((SKIPPED + 1))
    continue
  fi

  # Ejecutar migracion
  echo -ne "  ${YELLOW}▸${NC} $filename ... "
  if docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -f "/tmp/migrations/$filename" > /dev/null 2>&1; then
    # Registrar migracion aplicada
    docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -c \
      "INSERT INTO schema_migrations (filename) VALUES ('$filename');" > /dev/null 2>&1
    echo -e "${GREEN}OK${NC}"
    APPLIED=$((APPLIED + 1))
  else
    echo -e "${RED}FALLO${NC}"
    echo ""
    echo -e "${RED}Error ejecutando $filename. Detalle:${NC}"
    docker exec $CONTAINER psql -U $DB_USER -d $DB_NAME -f "/tmp/migrations/$filename"
    exit 1
  fi
done

# Limpiar archivos temporales
docker exec $CONTAINER rm -rf /tmp/migrations > /dev/null 2>&1

echo ""
echo "========================"
echo -e "Aplicadas: ${GREEN}$APPLIED${NC} | Ya existentes: $SKIPPED | Total: $((APPLIED + SKIPPED))"
echo ""

if [ "$APPLIED" -gt 0 ]; then
  echo -e "${GREEN}Migraciones completadas.${NC}"
else
  echo "Base de datos al dia, nada que migrar."
fi
