#!/bin/bash
# Estado del backend de CustomerMX en EC2
set -e

# Colores
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
INFRA_DIR="$REPO_ROOT/infraestructura"
SSH_KEY="$INFRA_DIR/customermx-key.pem"
REMOTE_USER="ec2-user"
SERVICE_NAME="customermx"
LOG_LINES="${1:-30}"   # Número de líneas de log (arg opcional, default 30)

fail() { echo -e "${RED}✗ $1${NC}"; exit 1; }

# ── Obtener IP ───────────────────────────────────────────────────────────────
if ! command -v terraform &>/dev/null; then
  fail "terraform no está instalado"
fi

EC2_IP=$(cd "$INFRA_DIR" && terraform output -raw ec2_public_ip 2>/dev/null) \
  || fail "No se pudo obtener ec2_public_ip. ¿Corriste 'terraform apply'?"

[ -z "$EC2_IP" ] && fail "ec2_public_ip está vacío"

SSH_OPTS="-i $SSH_KEY -o StrictHostKeyChecking=no -o ConnectTimeout=10"
SSH="ssh $SSH_OPTS $REMOTE_USER@$EC2_IP"

echo -e "${BOLD}CustomerMX Backend Status${NC}"
echo "═══════════════════════════════════════"
echo -e "  Servidor : ${CYAN}$EC2_IP${NC}"
echo "═══════════════════════════════════════"

# ── Systemd status ───────────────────────────────────────────────────────────
echo -e "\n${BOLD}Servicio systemd:${NC}"
$SSH bash <<REMOTE
  STATUS=\$(systemctl is-active $SERVICE_NAME 2>/dev/null || echo "unknown")

  if [ "\$STATUS" = "active" ]; then
    echo -e "  Estado  : \033[0;32m● activo (running)\033[0m"
  elif [ "\$STATUS" = "failed" ]; then
    echo -e "  Estado  : \033[0;31m✗ fallado\033[0m"
  else
    echo -e "  Estado  : \033[1;33m○ \$STATUS\033[0m"
  fi

  # PID, uptime, memoria
  systemctl show $SERVICE_NAME --no-pager \
    --property=MainPID,ActiveEnterTimestamp,MemoryCurrent \
    2>/dev/null | while IFS='=' read key val; do
    case \$key in
      MainPID)             [ "\$val" != "0" ] && echo "  PID     : \$val" ;;
      ActiveEnterTimestamp) echo "  Desde   : \$val" ;;
      MemoryCurrent)
        if [[ "\$val" =~ ^[0-9]+$ ]] && [ "\$val" -gt 0 ]; then
          mb=\$(( val / 1024 / 1024 ))
          echo "  Memoria : \${mb} MB"
        fi
        ;;
    esac
  done
REMOTE

# ── Health check HTTP ────────────────────────────────────────────────────────
echo -e "\n${BOLD}Health check:${NC}"
if curl -sf --max-time 5 "http://$EC2_IP:8080/health" >/dev/null 2>&1; then
  BODY=$(curl -sf --max-time 5 "http://$EC2_IP:8080/health" 2>/dev/null || echo "")
  echo -e "  HTTP    : ${GREEN}OK${NC} — http://$EC2_IP:8080/health"
  [ -n "$BODY" ] && echo "  Body    : $BODY"
else
  echo -e "  HTTP    : ${RED}Sin respuesta${NC} — http://$EC2_IP:8080/health"
fi

# HTTPS
if curl -sf --max-time 5 "https://api.customermx.com/health" >/dev/null 2>&1; then
  echo -e "  HTTPS   : ${GREEN}OK${NC} — https://api.customermx.com/health"
else
  echo -e "  HTTPS   : ${YELLOW}No disponible${NC} — https://api.customermx.com/health"
fi

# ── CORS — ALLOWED_ORIGINS en el servidor ───────────────────────────────────
echo -e "\n${BOLD}CORS (ALLOWED_ORIGINS en EC2):${NC}"
EXPECTED_ORIGINS="https://customermx.com,https://www.customermx.com,https://api.customermx.com,http://localhost:5173"
$SSH bash <<REMOTE
  ENV_FILE="/opt/customermx/.env"
  if [ ! -f "\$ENV_FILE" ]; then
    echo -e "  \033[0;31m✗ No se encontró \$ENV_FILE\033[0m"
  else
    ORIGINS=\$(grep '^ALLOWED_ORIGINS=' "\$ENV_FILE" | cut -d'=' -f2-)
    if [ -z "\$ORIGINS" ]; then
      echo -e "  \033[0;31m✗ ALLOWED_ORIGINS no está definido en .env\033[0m"
    else
      echo -e "  Valor   : \$ORIGINS"
      # Verificar cada origen esperado
      ALL_OK=true
      for origin in https://customermx.com https://www.customermx.com https://api.customermx.com; do
        if echo "\$ORIGINS" | grep -q "\$origin"; then
          echo -e "  \033[0;32m✓\033[0m \$origin"
        else
          echo -e "  \033[0;31m✗ FALTA: \$origin\033[0m"
          ALL_OK=false
        fi
      done
      \$ALL_OK && echo -e "  \033[0;32mTodos los orígenes de producción están presentes.\033[0m"
    fi
  fi
REMOTE

# ── Últimos logs ─────────────────────────────────────────────────────────────
echo -e "\n${BOLD}Últimos $LOG_LINES logs:${NC}"
echo "───────────────────────────────────────"
$SSH "sudo journalctl -u $SERVICE_NAME -n $LOG_LINES --no-pager --output=short 2>/dev/null"

echo -e "\n${BOLD}Uso:${NC}  ./status.sh [N_lineas_log]   (default: 30)"
