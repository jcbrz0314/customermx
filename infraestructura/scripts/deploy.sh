#!/bin/bash
# =============================================================
# CustomerMX - Script de deploy completo
# Ejecutar desde la carpeta raíz del proyecto: ./infraestructura/scripts/deploy.sh
# =============================================================
set -e

export AWS_ACCESS_KEY_ID=AKIAQQZQFAZZ5IAAYKP3
export AWS_SECRET_ACCESS_KEY="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"
export AWS_DEFAULT_REGION=us-west-2

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
INFRA_DIR="$REPO_ROOT/infraestructura"
KEY_FILE="$INFRA_DIR/customermx-key.pem"
SSH_OPTS="-i $KEY_FILE -o StrictHostKeyChecking=no -o ConnectTimeout=30"

echo "========================================"
echo "  CustomerMX Deploy"
echo "========================================"

# ----- Leer outputs de Terraform -----
echo ""
echo "[1/6] Leyendo outputs de Terraform..."
cd "$INFRA_DIR"

EC2_IP=$(terraform output -raw ec2_public_ip)
RDS_HOST=$(terraform output -raw rds_endpoint)
S3_BUCKET=$(terraform output -raw s3_bucket)
AMPLIFY_APP_ID=$(terraform output -raw amplify_app_id)
AMPLIFY_URL=$(terraform output -raw amplify_url)
JWT_ACCESS=$(terraform output -raw jwt_access_secret)
JWT_REFRESH=$(terraform output -raw jwt_refresh_secret)

echo "  EC2 IP:       $EC2_IP"
echo "  RDS Host:     $RDS_HOST"
echo "  S3 Bucket:    $S3_BUCKET"
echo "  Amplify URL:  $AMPLIFY_URL"

# ----- Compilar backend para Linux -----
echo ""
echo "[2/6] Compilando backend para Linux/amd64..."
cd "$REPO_ROOT/backend"

GOOS=linux GOARCH=amd64 go build -o /tmp/customermx-api ./cmd/api
GOOS=linux GOARCH=amd64 go build -o /tmp/customermx-migrate ./cmd/migrate

echo "  OK — binarios en /tmp/"

# ----- Esperar a que EC2 esté lista -----
echo ""
echo "[3/6] Esperando a que EC2 acepte conexiones SSH..."
for i in $(seq 1 20); do
  if ssh $SSH_OPTS ec2-user@$EC2_IP "echo ok" 2>/dev/null; then
    echo "  EC2 lista."
    break
  fi
  echo "  Intento $i/20 — esperando 10s..."
  sleep 10
done

# ----- Subir binarios y migrations al EC2 -----
echo ""
echo "[4/6] Subiendo archivos al EC2..."

scp $SSH_OPTS /tmp/customermx-api /tmp/customermx-migrate ec2-user@$EC2_IP:/opt/customermx/
scp $SSH_OPTS -r "$REPO_ROOT/backend/migrations/"*.sql ec2-user@$EC2_IP:/opt/customermx/migrations/

chmod +x /tmp/customermx-api /tmp/customermx-migrate 2>/dev/null || true

# Crear .env en EC2
ssh $SSH_OPTS ec2-user@$EC2_IP "cat > /opt/customermx/.env" << ENVEOF
# Server
SERVER_PORT=8080
SERVER_ENV=production

# Database
DB_HOST=${RDS_HOST}
DB_PORT=5432
DB_USER=cmxadmin
DB_PASSWORD=AlgoSeguro123!
DB_NAME=customermx
DB_SSL_MODE=disable

# JWT
JWT_ACCESS_SECRET=${JWT_ACCESS}
JWT_REFRESH_SECRET=${JWT_REFRESH}
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# S3 (AWS real — EC2 usa IAM role, no keys hardcodeadas)
AWS_REGION=us-west-2
S3_BUCKET=${S3_BUCKET}
S3_USE_PATH_STYLE=false

# Email (Google SMTP)
EMAIL_PROVIDER=smtp
EMAIL_FROM=CustomerMX <paperless@gsf-hotels.com.mx>
FRONTEND_URL=${AMPLIFY_URL}
EMAIL_LOGO_URL=

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=paperless@gsf-hotels.com.mx
SMTP_PASSWORD=pfkxqjvjlyqjopfj
SMTP_USE_TLS=true

# CORS
ALLOWED_ORIGINS=${AMPLIFY_URL},http://localhost:5173
ENVEOF

echo "  .env creado en EC2"

# Correr migraciones
echo ""
echo "  Ejecutando migraciones..."
ssh $SSH_OPTS ec2-user@$EC2_IP \
  "cd /opt/customermx && DB_HOST=$RDS_HOST DB_PORT=5432 DB_USER=cmxadmin DB_PASSWORD=AlgoSeguro123! DB_NAME=customermx ./customermx-migrate /opt/customermx/migrations"

# Habilitar y (re)iniciar servicio
ssh $SSH_OPTS ec2-user@$EC2_IP \
  "sudo systemctl enable customermx && sudo systemctl restart customermx && sleep 2 && sudo systemctl status customermx --no-pager"

echo "  Backend corriendo en http://$EC2_IP:8080"

# ----- Build del frontend -----
echo ""
echo "[5/6] Compilando frontend..."
cd "$REPO_ROOT/frontend"

VITE_API_URL="https://api.customermx.com/api/v1" npm run build

echo "  Build generado en frontend/dist/"

# ----- Deploy frontend a Amplify -----
echo ""
echo "[6/6] Desplegando frontend en Amplify..."
cd "$REPO_ROOT/frontend/dist"

zip -r /tmp/frontend-build.zip . -x "*.DS_Store"

DEPLOYMENT_JSON=$(aws amplify create-deployment \
  --app-id "$AMPLIFY_APP_ID" \
  --branch-name main \
  --region us-west-2 \
  --output json)

UPLOAD_URL=$(echo "$DEPLOYMENT_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['zipUploadUrl'])")
JOB_ID=$(echo "$DEPLOYMENT_JSON" | python3 -c "import sys,json; print(json.load(sys.stdin)['jobId'])")

echo "  Subiendo zip a Amplify..."
curl -s -T /tmp/frontend-build.zip "$UPLOAD_URL"

echo "  Iniciando deployment (job: $JOB_ID)..."
aws amplify start-deployment \
  --app-id "$AMPLIFY_APP_ID" \
  --branch-name main \
  --job-id "$JOB_ID" \
  --region us-west-2

echo ""
echo "========================================"
echo "  Deploy completado!"
echo "========================================"
echo ""
echo "  Frontend: $AMPLIFY_URL"
echo "  Backend:  http://$EC2_IP:8080"
echo "  SSH:      ssh -i infraestructura/customermx-key.pem ec2-user@$EC2_IP"
echo "  RDS:      postgresql://cmxadmin:AlgoSeguro123!@${RDS_HOST}:5432/customermx"
echo ""
echo "  El frontend de Amplify puede tardar 1-2 min en estar activo."
echo ""

# Limpiar temporales
rm -f /tmp/customermx-api /tmp/customermx-migrate /tmp/frontend-build.zip
