#!/bin/bash

# runenv.sh - Inicia todos los servicios de CustomerMX

# Detener todo primero
./stopenv.sh

echo ""
echo "🚀 Iniciando CustomerMX..."
echo ""

# Iniciar PostgreSQL si no está corriendo
echo "Verificando Base de Datos..."
if [ ! "$(docker ps -q -f name=customermx-postgres)" ]; then
  if [ "$(docker ps -aq -f status=exited -f name=customermx-postgres)" ]; then
    echo "Iniciando contenedor existente..."
    docker start customermx-postgres
  else
    echo "Creando y iniciando contenedor..."
    docker run -d \
      --name customermx-postgres \
      -e POSTGRES_USER=customermx \
      -e POSTGRES_PASSWORD=customermx123 \
      -e POSTGRES_DB=customermx \
      -p 5432:5432 \
      postgres:16-alpine
  fi
  echo "Esperando a que PostgreSQL esté listo..."
  sleep 3
fi
echo "✓ Base de Datos lista"
echo ""

# Iniciar Backend
echo "Iniciando Backend..."
cd backend
go build -o api cmd/api/main.go
if [ $? -eq 0 ]; then
  nohup ./api > ../logs/backend.log 2>&1 &
  echo "✓ Backend iniciado en puerto 8080"
  echo "  Log: logs/backend.log"
else
  echo "❌ Error compilando Backend"
  exit 1
fi
cd ..
echo ""

# Esperar un momento para que el backend inicie
sleep 2

# Iniciar Frontend
echo "Iniciando Frontend..."
mkdir -p logs
nohup npm --prefix frontend run dev > logs/frontend.log 2>&1 &
echo "✓ Frontend iniciado en puerto 5173"
echo "  Log: logs/frontend.log"
echo ""

echo "✅ CustomerMX está corriendo"
echo ""
echo "📍 URLs:"
echo "   Frontend: http://localhost:5173"
echo "   Backend:  http://localhost:8080"
echo "   Database: localhost:5432"
echo ""
echo "📋 Logs:"
echo "   Backend:  tail -f logs/backend.log"
echo "   Frontend: tail -f logs/frontend.log"
echo ""
echo "🛑 Para detener: ./stopenv.sh"
