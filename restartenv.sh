#!/bin/bash

# restartenv.sh - Reinicia todos los servicios de CustomerMX

echo "🔄 Reiniciando CustomerMX..."
echo ""

./stopenv.sh
echo ""
./runenv.sh
