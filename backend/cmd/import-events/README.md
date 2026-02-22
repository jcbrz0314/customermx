# Import Events - Excel to PostgreSQL

Script en Go para importar eventos históricos desde un archivo Excel gigante a la base de datos PostgreSQL.

---

## 📋 Requisitos

1. **Go 1.21+** instalado
2. **PostgreSQL** corriendo en Docker (puerto 5432)
3. **Archivo Excel** llamado `eventos.xlsx` en la raíz del proyecto
4. **Dependencias Go:**
   ```bash
   go get github.com/xuri/excelize/v2
   go get github.com/jackc/pgx/v5/pgxpool
   go get github.com/google/uuid
   ```

---

## 📁 Estructura del Excel

### Hoja "General"

**Fila 1:** Nombres de marcas (para columnas T en adelante)
**Fila 2:** Encabezados de columnas:

| Columna | Campo | Descripción |
|---------|-------|-------------|
| A | Numerador | Ignorado |
| B | Marca | Marca del evento (Chevrolet, Buick, GMC, Cadillac) |
| C | Tipo de Evento | Tipo (Triatlón, Maratón, etc.) |
| D | Organizador | Nombre del organizador |
| E | Nombre Evento | **REQUERIDO** - Nombre completo del evento |
| F | Fecha | Fecha de inicio (varios formatos aceptados) |
| G | Año | Año del evento |
| H | Dias de Evento | Duración en días |
| I | Estado | Estado de México |
| J | Ciudad | Ciudad |
| K | Sede | Nombre de la sede |
| L | Distribuidor | Nombre del distribuidor |
| M | Edecanes | Número de edecanes |
| N | Montaje | Proveedor de montaje |
| O | Promocional | Si/No (boolean) |
| P | Asistencia | Número de asistentes |
| Q | Dinamicas | Número de dinámicas/actividades |
| R | Datos Levantados | Leads capturados |
| S | Prospectos | Número de prospectos |
| T-BG | Vehículos | Cantidades por modelo (agrupados por marca en fila 1) |
| BH | Calificación | Calificación del distribuidor (1-5) |
| BI | Comentarios | Comentarios adicionales |

**Fila 3+:** Datos de eventos

---

## 🚀 Uso

### 1. Preparar el archivo Excel

Coloca tu archivo Excel en la raíz del proyecto con el nombre `eventos.xlsx`:

```bash
customermx/
├── eventos.xlsx          # Tu archivo Excel
├── backend/
│   ├── cmd/
│   │   └── import-events/
│   │       ├── main.go
│   │       └── README.md
```

### 2. Asegurar que PostgreSQL está corriendo

```bash
cd customermx
docker-compose up -d postgres
```

Verifica la conexión:
```bash
docker exec -it customermx-postgres-1 psql -U postgres -d customermx -c "SELECT COUNT(*) FROM brands;"
```

### 3. Ejecutar el script

Desde la raíz del proyecto:

```bash
cd backend/cmd/import-events
go run main.go
```

O compilar y ejecutar:

```bash
# Compilar
go build -o import-events main.go

# Ejecutar
./import-events
```

---

## 📊 Proceso de Importación

El script realiza los siguientes pasos en orden:

1. **📂 Abre el archivo Excel** `eventos.xlsx`
2. **🗄️ Conecta a PostgreSQL** y carga el catálogo existente de marcas y vehículos
3. **📋 Parsea encabezados:**
   - Fila 1: Marcas de vehículos
   - Fila 2: Modelos de vehículos
   - Detecta automáticamente nuevas marcas y vehículos
4. **📊 Parsea datos de eventos** (fila 3 en adelante)
5. **⚠️ Muestra advertencias** si encuentra marcas o vehículos nuevos
6. **📝 Genera archivos SQL** en `backend/migrations/import/`:
   - `01_new_brands.sql` - Solo si hay marcas nuevas
   - `02_new_vehicles.sql` - Solo si hay vehículos nuevos
   - `03_events.sql` - Todos los eventos
   - `04_event_reports.sql` - Reportes de eventos
   - `05_event_vehicles.sql` - Vehículos por evento
7. **💾 Ejecuta SQL en la base de datos** usando una transacción
8. **✅ Confirma éxito** con resumen de operación

---

## 🔍 Validaciones y Reglas

### Validaciones Automáticas

- ✅ **Marca del evento debe existir** - Si no existe, se crea automáticamente
- ✅ **Vehículos deben coincidir con marca del evento** - FK en base de datos lo garantiza
- ✅ **Cantidad de vehículos > 0** - Solo inserta si hay cantidad
- ✅ **Fechas inválidas** - Usa fecha actual si no se puede parsear
- ✅ **Calificación 1-5** - Valida rango antes de insertar
- ✅ **Event_id único en reportes** - Constraint de BD lo garantiza

### Reglas de Negocio

1. **Eventos históricos tienen status COMPLETED**
2. **Reportes marcados como completed = true**
3. **Campos opcionales se omiten si están vacíos** (NULL en BD)
4. **Promocional acepta:** "si", "sí", "yes", "true", "1"
5. **Inserts puros, NO idempotentes** - Cada ejecución inserta nuevos registros

### Manejo de Datos Faltantes

| Campo | Si está vacío |
|-------|---------------|
| Marca evento | ❌ Error - requerido |
| Nombre evento | ❌ Error - requerido |
| Fecha | ⚠️ Usa fecha actual |
| Año | ⚠️ Extrae de fecha |
| Días evento | Usa 1 |
| Otros campos | NULL |

---

## 📁 Archivos Generados

Los archivos SQL se generan en: `backend/migrations/import/`

```bash
backend/migrations/import/
├── 01_new_brands.sql       # Solo si hay marcas nuevas
├── 02_new_vehicles.sql     # Solo si hay vehículos nuevos
├── 03_events.sql           # Siempre
├── 04_event_reports.sql    # Siempre
└── 05_event_vehicles.sql   # Siempre
```

### Ejemplo de Salida

```sql
-- 03_events.sql
INSERT INTO events (id, brand_id, event_type, organizer, name, start_date, year, duration_days, state, city, venue, dealer, status)
VALUES ('550e8400-e29b-41d4-a716-446655440000', '7c9e6679-7425-40de-944b-e07fc1f90ae7', 'Triatlón', 'Chevrolet México', 'Triatlón Chevrolet Guadalajara 2024', '2024-06-15', 2024, 2, 'Jalisco', 'Guadalajara', 'Expo Guadalajara', 'Chevrolet Andares', 'COMPLETED');

-- 04_event_reports.sql
INSERT INTO event_reports (id, event_id, completed, attendees, leads_collected, prospects, dealer_rating)
VALUES ('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', true, 500, 120, 45, 5);

-- 05_event_vehicles.sql
INSERT INTO event_vehicles (id, event_id, vehicle_id, quantity)
VALUES ('770e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', '8c9e6679-7425-40de-944b-e07fc1f90ae8', 3);
```

---

## ⚠️ Advertencias y Logs

El script muestra logs detallados durante la ejecución:

```bash
🚀 Iniciando importación de eventos desde Excel...
📂 Abriendo archivo Excel: eventos.xlsx
🗄️  Conectando a base de datos...
📚 Cargando catálogo de marcas y vehículos existentes...
✅ Catálogo cargado: 4 marcas, 39 vehículos
📋 Parseando encabezados de Excel...
✅ Encabezados parseados: 35 columnas de vehículos
📊 Parseando datos de eventos...
✅ Datos parseados: 150 eventos, 150 reportes, 450 event_vehicles

⚠️  Se encontraron 2 vehículos nuevos que se insertarán:
   - Chevrolet: Equinox 2025
   - GMC: Yukon XL

📝 Generando archivos SQL...
✅ Archivos SQL generados en: backend/migrations/import
💾 Ejecutando SQL en base de datos...
✅ ¡Importación completada exitosamente!

📊 Resumen:
   - Marcas nuevas: 0
   - Vehículos nuevos: 2
   - Eventos insertados: 150
   - Reportes insertados: 150
   - Vehículos de eventos insertados: 450
```

---

## 🐛 Troubleshooting

### Error: "archivo Excel no encontrado"

**Problema:** El archivo `eventos.xlsx` no está en la raíz del proyecto.

**Solución:**
```bash
# Verificar ubicación
ls -la eventos.xlsx

# Debe estar en:
# /Users/josebeltran/Documents/GitHub/customermx/eventos.xlsx
```

### Error: "error al conectar a base de datos"

**Problema:** PostgreSQL no está corriendo o credenciales incorrectas.

**Solución:**
```bash
# Verificar que Docker esté corriendo
docker ps | grep postgres

# Iniciar PostgreSQL
docker-compose up -d postgres

# Verificar conexión
docker exec -it customermx-postgres-1 psql -U postgres -d customermx -c "\dt"
```

### Error: "marca no encontrada"

**Problema:** El Excel tiene una marca que no está en el catálogo y no se pudo crear.

**Solución:** El script automáticamente crea nuevas marcas. Si falla, verifica:
1. Formato del nombre de marca en Excel (sin espacios extra)
2. Logs de advertencia que el script genera

### Advertencia: "vehículo no encontrado"

**Problema:** El Excel menciona un vehículo que no existe.

**Solución:** El script automáticamente crea nuevos vehículos. Revisa los logs:
```
⚠️  Se encontraron 3 vehículos nuevos que se insertarán:
   - Chevrolet: Modelo Nuevo
```

### Error: "duplicate key value violates unique constraint"

**Problema:** Intentando insertar un evento o reporte que ya existe.

**Solución:**
- El script NO es idempotente
- Cada ejecución inserta nuevos registros
- Para re-ejecutar, primero limpia los datos:
```sql
-- ⚠️ CUIDADO: Esto borra TODOS los eventos
TRUNCATE event_vehicles, event_reports, events CASCADE;
```

---

## 🔄 Re-ejecutar Importación

Si necesitas re-ejecutar el script:

### Opción 1: Limpiar solo datos importados (recomendado)

```sql
-- Identificar eventos del import por fecha de creación
SELECT id, name, created_at FROM events
WHERE created_at > '2026-02-13'
ORDER BY created_at DESC;

-- Eliminar manualmente los IDs específicos
DELETE FROM event_vehicles WHERE event_id IN (...);
DELETE FROM event_reports WHERE event_id IN (...);
DELETE FROM events WHERE id IN (...);
```

### Opción 2: Limpiar TODOS los eventos (⚠️ DESTRUCTIVO)

```bash
docker exec -it customermx-postgres-1 psql -U postgres -d customermx -c "
TRUNCATE event_vehicles, event_reports, events CASCADE;
"
```

Luego re-ejecutar:
```bash
go run main.go
```

---

## 📝 Notas Importantes

### Performance

- **Manejo de memoria:** El script carga todo el Excel en memoria
- **Transacciones:** Usa una única transacción para atomicidad
- **Archivos grandes:** Para Excel > 100MB, considera aumentar memoria Go:
  ```bash
  GOGC=50 go run main.go
  ```

### Formato de Fechas Soportado

El script intenta parsear fechas en estos formatos:
- `2024-06-15` (YYYY-MM-DD) ← Preferido
- `15/06/2024` (DD/MM/YYYY)
- `06/15/2024` (MM/DD/YYYY)
- `2024/06/15` (YYYY/MM/DD)
- `15-06-2024` (DD-MM-YYYY)
- `06-15-2024` (MM-DD-YYYY)

### Caracteres Especiales

El script escapa automáticamente caracteres especiales en SQL:
- Comillas simples `'` → `''`
- Otros caracteres especiales se mantienen

### UUIDs

Todos los IDs se generan automáticamente usando UUID v4.

---

## 🔧 Configuración

Para cambiar la configuración, edita las constantes al inicio de `main.go`:

```go
const (
    ExcelFile    = "eventos.xlsx"           // Nombre del archivo Excel
    SheetName    = "General"                // Nombre de la hoja
    OutputDir    = "backend/migrations/import"  // Directorio de salida
    DBConnString = "postgres://postgres:postgres@localhost:5432/customermx?sslmode=disable"
)
```

---

## 📊 Ejemplo de Uso Completo

```bash
# 1. Preparar entorno
cd /Users/josebeltran/Documents/GitHub/customermx
docker-compose up -d postgres

# 2. Verificar archivo Excel
ls -la eventos.xlsx

# 3. Ejecutar import
cd backend/cmd/import-events
go run main.go

# 4. Verificar resultados
docker exec -it customermx-postgres-1 psql -U postgres -d customermx

# En psql:
SELECT COUNT(*) FROM events;
SELECT COUNT(*) FROM event_reports;
SELECT COUNT(*) FROM event_vehicles;
SELECT e.name, e.year, b.name as brand
FROM events e
JOIN brands b ON e.brand_id = b.id
ORDER BY e.start_date DESC
LIMIT 10;
```

---

## ✅ Checklist Pre-Ejecución

Antes de ejecutar el script, verifica:

- [ ] PostgreSQL está corriendo (`docker ps | grep postgres`)
- [ ] Archivo `eventos.xlsx` está en la raíz del proyecto
- [ ] El archivo Excel tiene la hoja "General"
- [ ] Los encabezados están en la fila 2
- [ ] Los datos empiezan en la fila 3
- [ ] Dependencias Go instaladas (`go mod download`)
- [ ] Tienes permisos de escritura en `backend/migrations/import/`

---

**Desarrollado por:** Claude Code (Sonnet 4.5)
**Fecha:** Febrero 2026
**Versión:** 1.0.0



  go run -C backend ./cmd/import-events/ /Users/josebeltran/Documents/GitHub/customermx/eventos.xlsx