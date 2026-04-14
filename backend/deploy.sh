#!/bin/bash
# Deploy del backend de CustomerMX a EC2
set -e

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
INFRA_DIR="$REPO_ROOT/infraestructura"
SSH_KEY="$INFRA_DIR/customermx-key.pem"
REMOTE_USER="ec2-user"
REMOTE_DIR="/opt/customermx"
SERVICE_NAME="customermx"

step() { echo -e "\n${CYAN}==>${NC} $1"; }
ok()   { echo -e "${GREEN}✓${NC} $1"; }
fail() { echo -e "${RED}✗ $1${NC}"; exit 1; }

# ── Obtener IP del EC2 desde Terraform ──────────────────────────────────────
step "Obteniendo IP del EC2 desde Terraform..."
if ! command -v terraform &>/dev/null; then
  fail "terraform no está instalado"
fi

EC2_IP=$(cd "$INFRA_DIR" && terraform output -raw ec2_public_ip 2>/dev/null) \
  || fail "No se pudo obtener ec2_public_ip. ¿Corriste 'terraform apply'?"

[ -z "$EC2_IP" ] && fail "ec2_public_ip está vacío"
ok "EC2: $EC2_IP"

SSH_OPTS="-i $SSH_KEY -o StrictHostKeyChecking=no -o ConnectTimeout=10"
SSH="ssh $SSH_OPTS $REMOTE_USER@$EC2_IP"
SCP="scp $SSH_OPTS"

# ── Build del binario para Linux ─────────────────────────────────────────────
step "Compilando para linux/amd64..."
cd "$SCRIPT_DIR"
GOOS=linux GOARCH=amd64 go build -o /tmp/customermx-api ./cmd/api \
  || fail "Falló el build del API"
ok "Binary listo en /tmp/customermx-api"

# ── Subir binario al EC2 ─────────────────────────────────────────────────────
step "Subiendo binario al EC2..."
$SCP /tmp/customermx-api "$REMOTE_USER@$EC2_IP:/tmp/customermx-api" \
  || fail "Falló la copia del binario"
ok "Binario copiado"

# ── Subir migraciones ────────────────────────────────────────────────────────
step "Subiendo migraciones..."
$SCP "$SCRIPT_DIR/migrations"/*.sql "$REMOTE_USER@$EC2_IP:/tmp/" \
  || fail "Falló la copia de las migraciones"
ok "Migraciones copiadas"

# ── Aplicar en el servidor ───────────────────────────────────────────────────
step "Instalando y reiniciando en el servidor..."
$SSH bash <<REMOTE
  set -e

  # Parar el servicio
  sudo systemctl stop $SERVICE_NAME 2>/dev/null || true

  # Instalar binario
  sudo mv /tmp/customermx-api $REMOTE_DIR/customermx-api
  sudo chmod +x $REMOTE_DIR/customermx-api
  sudo chown ec2-user:ec2-user $REMOTE_DIR/customermx-api

  # Instalar migraciones
  sudo mkdir -p $REMOTE_DIR/migrations
  sudo cp /tmp/V*.sql $REMOTE_DIR/migrations/ 2>/dev/null || true
  sudo chown -R ec2-user:ec2-user $REMOTE_DIR/migrations

  # Arrancar el servicio
  sudo systemctl start $SERVICE_NAME
  sudo systemctl enable $SERVICE_NAME

  # Esperar arranque
  sleep 2
  sudo systemctl is-active $SERVICE_NAME
REMOTE

ok "Deploy completado"

# ── Verificar health ─────────────────────────────────────────────────────────
step "Verificando health del servicio..."
sleep 2
if curl -sf --max-time 5 "http://$EC2_IP:8080/health" >/dev/null 2>&1; then
  ok "Health OK — http://$EC2_IP:8080"
else
  echo -e "${YELLOW}⚠ Health check no respondió (puede tardar unos segundos en arrancar)${NC}"
  echo "   Verifica con: ./status.sh"
fi

echo -e "\n${GREEN}Deploy finalizado.${NC}"
