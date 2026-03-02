#!/bin/bash

# localstack.sh - Crea o inicia el contenedor LocalStack de CustomerMX (puerto 4666)

CONTAINER_NAME="customermx-localstack"
INIT_DIR="$(cd "$(dirname "$0")" && pwd)/backend/localstack-init"

start_localstack() {
  if [ "$(docker ps -q -f name=${CONTAINER_NAME})" ]; then
    echo "✓ LocalStack ya está corriendo"
    return 0
  fi

  if [ "$(docker ps -aq -f status=exited -f name=${CONTAINER_NAME})" ]; then
    echo "Iniciando contenedor existente..."
    docker start "${CONTAINER_NAME}"
  else
    echo "Creando y iniciando contenedor..."
    docker run -d \
      --name "${CONTAINER_NAME}" \
      -p 4666:4566 \
      -e SERVICES=s3 \
      -e DEFAULT_REGION=us-east-1 \
      -e PERSISTENCE=0 \
      -v "${INIT_DIR}:/etc/localstack/init/ready.d" \
      localstack/localstack:3
  fi

  echo "Esperando a que LocalStack esté listo..."
  for i in $(seq 1 15); do
    if curl -s http://localhost:4666/_localstack/health | grep -q '"s3"'; then
      echo "✓ LocalStack listo (S3 disponible en http://localhost:4666)"
      return 0
    fi
    sleep 1
  done

  echo "⚠️  LocalStack inició pero S3 puede tardar unos segundos más en estar disponible"
}

start_localstack
