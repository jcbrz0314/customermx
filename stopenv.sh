#!/bin/bash

# stopenv.sh - Detiene todos los servicios de CustomerMX

echo "🛑 Deteniendo CustomerMX..."

# Detener Frontend (Vite)
echo "Deteniendo Frontend..."
FRONTEND_PID=$(lsof -ti:5173,5174)
if [ -n "$FRONTEND_PID" ]; then
  kill -9 $FRONTEND_PID 2>/dev/null
  echo "✓ Frontend detenido"
else
  echo "✓ Frontend no está corriendo"
fi

# Detener Backend (Go API)
echo "Deteniendo Backend..."
BACKEND_PID=$(lsof -ti:8080)
if [ -n "$BACKEND_PID" ]; then
  kill -9 $BACKEND_PID 2>/dev/null
  echo "✓ Backend detenido"
else
  echo "✓ Backend no está corriendo"
fi

# Limpiar procesos de npm/vite que puedan quedar
pkill -f "vite" 2>/dev/null
pkill -f "npm.*dev" 2>/dev/null

echo "✅ Todos los servicios detenidos"
