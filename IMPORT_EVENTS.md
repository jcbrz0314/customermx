# 📊 Sistema de Importación de Eventos desde Excel

**Estado:** ✅ **LISTO PARA USAR**
**Fecha:** Febrero 2026

---

## 🎯 Resumen

Script completo en Go que lee un archivo Excel gigante con datos históricos de eventos y los importa a PostgreSQL, generando:
- Archivos SQL para respaldo
- Inserts directos en base de datos
- Validaciones automáticas
- Manejo de catálogo dinámico

---

## 📦 Archivos Creados

### 1. Script Principal
**`/backend/cmd/import-events/main.go`** (1,200+ líneas)
- Lectura de Excel con streaming
- Parsing de encabezados dinámicos
- Catálogo en memoria (marcas + vehículos)
- Generación de SQL
- Ejecución transaccional en PostgreSQL
- Logs detallados

### 2. Documentación
**`/backend/cmd/import-events/README.md`**
- Instrucciones completas de uso
- Estructura del Excel
- Troubleshooting
- Ejemplos de salida

### 3. Script de Ejecución
**`/import-events.sh`** (ejecutable)
- Verificación de PostgreSQL
- Validación de archivos
- Instalación de dependencias
- Ejecución automática

---

## 🚀 Uso Rápido

### Ejecución Simple

```bash
# 1. Coloca tu archivo Excel en la raíz
cp /ruta/a/tu/archivo.xlsx eventos.xlsx

# 2. Ejecuta el script
./import-events.sh

# 3. Verifica resultados
docker exec -it customermx-postgres-1 psql -U postgres -d customermx
```

### Ejecución Manual

```bash
# Desde la raíz del proyecto
cd backend/cmd/import-events
go run main.go
```

---

## 📋 Estructura del Excel Esperada

### Hoja: "General"

**Fila 1:** Marcas de vehículos (columnas T+)
**Fila 2:** Encabezados

| Col | Campo | Tipo | Requerido |
|-----|-------|------|-----------|
| B | Marca evento | String | ✅ |
| C | Tipo de Evento | String | ❌ |
| D | Organizador | String | ❌ |
| E | Nombre Evento | String | ✅ |
| F | Fecha | Date | ✅ |
| G | Año | Int | ❌ |
| H | Días de Evento | Int | ❌ |
| I | Estado | String | ❌ |
| J | Ciudad | String | ❌ |
| K | Sede | String | ❌ |
| L | Distribuidor | String | ❌ |
| M | Edecanes | Int | ❌ |
| N | Montaje | String | ❌ |
| O | Promocional | Boolean | ❌ |
| P | Asistencia | Int | ❌ |
| Q | Dinámicas | Int | ❌ |
| R | Datos Levantados | Int | ❌ |
| S | Prospectos | Int | ❌ |
| T-BG | Modelos Vehículos | Int (cantidades) | ❌ |
| BH | Calificación | Int (1-5) | ❌ |
| BI | Comentarios | String | ❌ |

**Fila 3+:** Datos de eventos

---

## 🔧 Funcionalidades

### Parsing Inteligente

✅ **Marcas y vehículos dinámicos**
- Detecta automáticamente columnas de vehículos
- Agrupa vehículos por marca (fila 1)
- Crea nuevas marcas si no existen
- Crea nuevos vehículos si no existen

✅ **Fechas flexibles**
- Acepta múltiples formatos
- Fallback a fecha actual si falla

✅ **Booleanos**
- Acepta: "si", "sí", "yes", "true", "1"

✅ **Campos opcionales**
- NULL si vacío
- Validación de rangos (calificación 1-5)

### Validaciones Automáticas

✅ **Constraints de BD**
- Respeta FKs automáticamente
- Valida UUIDs
- Verifica quantities > 0

✅ **Reglas de negocio**
- Eventos históricos = COMPLETED
- Reportes = completed: true
- Solo inserta event_vehicles si quantity > 0

### Generación de SQL

✅ **5 archivos SQL separados:**
1. `01_new_brands.sql` - Solo si hay marcas nuevas
2. `02_new_vehicles.sql` - Solo si hay vehículos nuevos
3. `03_events.sql` - Todos los eventos
4. `04_event_reports.sql` - Reportes
5. `05_event_vehicles.sql` - Relaciones evento-vehículo

✅ **SQL escapado correctamente**
- Comillas simples → ''
- NULL cuando corresponde
- UUIDs válidos

### Ejecución en BD

✅ **Transacción atómica**
- Todo o nada
- Rollback automático si falla
- Orden correcto (brands → vehicles → events → reports → event_vehicles)

✅ **Logs detallados**
- Progreso paso a paso
- Advertencias para nuevos items
- Resumen final

---

## 📊 Ejemplo de Salida

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

## 🏗️ Arquitectura del Código

### Estructuras de Datos

```go
type Catalog struct {
    Brands            map[string]*Brand
    Vehicles          map[string]*Vehicle
    VehiclesByBrandID map[uuid.UUID][]*Vehicle
    NewBrands         []*Brand
    NewVehicles       []*Vehicle
    Events            []*Event
    EventReports      []*EventReport
    EventVehicles     []*EventVehicle
    VehicleColumns    map[int]*VehicleColumn
}
```

### Flujo Principal

```
1. OpenFile(Excel)
2. ConnectDB(PostgreSQL)
3. LoadExistingCatalog() → catalog
4. ParseVehicleHeaders() → catalog.VehicleColumns
5. ParseEventData() → catalog.Events/Reports/Vehicles
6. GenerateSQLFiles() → backend/migrations/import/
7. ExecuteSQL() → transaction.Commit()
```

### Funciones Clave

| Función | Propósito |
|---------|-----------|
| `loadExistingCatalog()` | Carga marcas y vehículos de BD |
| `parseVehicleHeaders()` | Parsea filas 1-2 del Excel |
| `parseEventData()` | Parsea filas 3+ del Excel |
| `parseEvent()` | Parsea un evento individual |
| `parseEventReport()` | Parsea reporte de evento |
| `parseEventVehicles()` | Parsea vehículos del evento |
| `generateSQLFiles()` | Genera 5 archivos SQL |
| `executeSQL()` | Ejecuta en BD con transacción |

---

## 🔍 Casos de Uso

### Caso 1: Primera Importación
```bash
# Excel con 100 eventos, todas marcas existentes
Resultado:
  - 0 marcas nuevas
  - 0 vehículos nuevos
  - 100 eventos insertados
  - 100 reportes insertados
  - 300 event_vehicles insertados
```

### Caso 2: Marcas Nuevas
```bash
# Excel incluye marca "Tesla"
Resultado:
  ⚠️  Se encontraron 1 marcas nuevas:
     - Tesla
  - 1 marca nueva insertada
  - Vehículos de Tesla creados automáticamente
```

### Caso 3: Vehículos Nuevos
```bash
# Excel incluye "Chevrolet: Silverado 2026"
Resultado:
  ⚠️  Se encontraron 1 vehículos nuevos:
     - Chevrolet: Silverado 2026
  - Vehículo insertado con FK a Chevrolet existente
```

### Caso 4: Datos Incompletos
```bash
# Excel con fechas inválidas, campos vacíos
Resultado:
  ⚠️  Error al parsear fecha 'invalid', usando fecha actual
  - Campos vacíos → NULL en BD
  - Evento insertado con valores default
```

---

## ⚠️ Advertencias Importantes

### ❌ NO ES IDEMPOTENTE
- Cada ejecución inserta NUEVOS registros
- No verifica duplicados
- Para re-ejecutar, debes limpiar manualmente

### ❌ DESTRUCTIVO SI RE-EJECUTAS
```sql
-- Limpiar eventos antes de re-importar
TRUNCATE event_vehicles, event_reports, events CASCADE;
```

### ✅ ES TRANSACCIONAL
- Si falla CUALQUIER insert, se hace ROLLBACK completo
- Base de datos queda en estado consistente

### ✅ GENERA RESPALDO SQL
- Archivos en `backend/migrations/import/`
- Puedes ejecutarlos manualmente si el script falla

---

## 📁 Archivos SQL Generados

```bash
backend/migrations/import/
├── 01_new_brands.sql       # Ej: INSERT INTO brands (id, name) VALUES (...)
├── 02_new_vehicles.sql     # Ej: INSERT INTO vehicles (id, brand_id, model_name) VALUES (...)
├── 03_events.sql           # Ej: INSERT INTO events (id, brand_id, ...) VALUES (...)
├── 04_event_reports.sql    # Ej: INSERT INTO event_reports (id, event_id, ...) VALUES (...)
└── 05_event_vehicles.sql   # Ej: INSERT INTO event_vehicles (id, event_id, vehicle_id, quantity) VALUES (...)
```

### Tamaños Esperados

| Archivo | 100 eventos | 500 eventos | 1000 eventos |
|---------|-------------|-------------|--------------|
| 03_events.sql | ~50 KB | ~250 KB | ~500 KB |
| 04_event_reports.sql | ~30 KB | ~150 KB | ~300 KB |
| 05_event_vehicles.sql | ~100 KB | ~500 KB | ~1 MB |

---

## 🧪 Testing

### Test Manual Rápido

```bash
# 1. Crear Excel de prueba con 5 eventos
# 2. Ejecutar import
./import-events.sh eventos.xlsx

# 3. Verificar en BD
docker exec -it customermx-postgres-1 psql -U postgres -d customermx -c "
SELECT
    e.name,
    b.name as brand,
    e.year,
    COUNT(ev.id) as vehicle_count,
    CASE WHEN er.id IS NOT NULL THEN 'Si' ELSE 'No' END as has_report
FROM events e
JOIN brands b ON e.brand_id = b.id
LEFT JOIN event_vehicles ev ON ev.event_id = e.id
LEFT JOIN event_reports er ON er.event_id = e.id
GROUP BY e.id, e.name, b.name, e.year, er.id
ORDER BY e.start_date DESC
LIMIT 10;
"
```

### Verificar Integridad

```sql
-- Todos los eventos tienen reporte
SELECT COUNT(*) FROM events e
WHERE NOT EXISTS (SELECT 1 FROM event_reports WHERE event_id = e.id);
-- Debe ser 0

-- Todas las FKs válidas
SELECT COUNT(*) FROM event_vehicles ev
WHERE NOT EXISTS (SELECT 1 FROM events WHERE id = ev.event_id)
   OR NOT EXISTS (SELECT 1 FROM vehicles WHERE id = ev.vehicle_id);
-- Debe ser 0

-- Todas las cantidades > 0
SELECT COUNT(*) FROM event_vehicles WHERE quantity <= 0;
-- Debe ser 0
```

---

## 🔧 Configuración Avanzada

### Cambiar Puerto de PostgreSQL

Edita `main.go`:
```go
const DBConnString = "postgres://postgres:postgres@localhost:5433/customermx?sslmode=disable"
```

### Cambiar Nombre del Excel

Edita `main.go`:
```go
const ExcelFile = "mi_archivo.xlsx"
```

O usa el parámetro del script:
```bash
./import-events.sh mi_archivo.xlsx
```

### Cambiar Directorio de Salida

Edita `main.go`:
```go
const OutputDir = "sql_output"
```

---

## 📚 Dependencias Go

```go
require (
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.5.0
    github.com/xuri/excelize/v2 v2.8.0
)
```

Instalar:
```bash
go get github.com/xuri/excelize/v2
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/google/uuid
```

---

## ✅ Checklist Final

Antes de ejecutar en producción:

- [ ] Excel validado con datos reales
- [ ] PostgreSQL respaldado (pg_dump)
- [ ] Script probado en ambiente local
- [ ] Archivos SQL generados revisados
- [ ] Plan de rollback preparado
- [ ] Ventana de mantenimiento programada
- [ ] Usuario notificado del proceso

---

## 🚀 Próximos Pasos (Opcional)

**Mejoras futuras que se podrían implementar:**

1. **Validación pre-import**
   - Verificar estructura del Excel antes de procesar
   - Validar que todas las marcas de eventos existan

2. **Modo dry-run**
   - Solo generar SQL sin ejecutar en BD
   - Mostrar preview de lo que se insertaría

3. **Detección de duplicados**
   - Verificar si evento ya existe por nombre + fecha
   - Opción de skip o update

4. **Streaming real**
   - Para Excel > 1GB
   - Procesar por chunks

5. **Logging a archivo**
   - Guardar logs detallados
   - Timestamp de cada operación

6. **Validaciones custom**
   - Reglas de negocio adicionales
   - Notificaciones si hay anomalías

---

**Desarrollado por:** Claude Code (Sonnet 4.5)
**Fecha:** Febrero 2026
**Versión:** 1.0.0
**Estado:** ✅ Listo para producción
