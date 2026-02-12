# Database Migrations

Este directorio contiene las migraciones de base de datos para CustomerMX.

## Migraciones Disponibles

- **V1__create_brands_and_vehicles.sql** - Tablas de marcas y vehículos
- **V2__create_users_and_invitations.sql** - Tablas de usuarios e invitaciones
- **V3__create_events_and_reports.sql** - Tablas de eventos, reportes y notificaciones
- **V4__seed_brands_and_vehicles.sql** - Datos iniciales (4 marcas, 39 modelos)

## Ejecutar Migraciones

### Opción 1: Usando psql (Desarrollo)

```bash
# Asegúrate de que PostgreSQL esté corriendo
make docker-up

# Ejecutar todas las migraciones en orden
psql -h localhost -U postgres -d customermx -f migrations/V1__create_brands_and_vehicles.sql
psql -h localhost -U postgres -d customermx -f migrations/V2__create_users_and_invitations.sql
psql -h localhost -U postgres -d customermx -f migrations/V3__create_events_and_reports.sql
psql -h localhost -U postgres -d customermx -f migrations/V4__seed_brands_and_vehicles.sql
```

### Opción 2: Usando el script de migración

```bash
# Ejecutar todas las migraciones
./migrations/run_migrations.sh
```

### Opción 3: Usando golang-migrate (Recomendado para producción)

```bash
# Instalar golang-migrate
brew install golang-migrate

# Ejecutar migraciones
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/customermx?sslmode=disable" up

# Rollback
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/customermx?sslmode=disable" down
```

## Verificar Estado de la Base de Datos

```bash
# Conectar a la base de datos
psql -h localhost -U postgres -d customermx

# Ver tablas creadas
\dt

# Ver datos de brands
SELECT * FROM brands;

# Ver cantidad de vehículos por marca
SELECT b.name, COUNT(v.id) as vehicle_count
FROM brands b
LEFT JOIN vehicles v ON b.id = v.brand_id
GROUP BY b.name;
```

## Estructura de la Base de Datos

### Tablas Principales
- `brands` - Marcas automotrices
- `vehicles` - Modelos de vehículos
- `users` - Usuarios del sistema (ADMIN, COORDINATOR, BRAND)
- `invitations` - Invitaciones pendientes
- `events` - Eventos promocionales
- `event_coordinators` - Asignación de coordinadores a eventos
- `event_reports` - Reportes operativos de eventos
- `event_vehicles` - Vehículos presentados en eventos
- `notifications` - Notificaciones de usuarios

### Tipos Enum
- `user_role` - ADMIN, COORDINATOR, BRAND
- `event_status` - PLANNED, ACTIVE, COMPLETED, CLOSED

## Seed Data

Después de ejecutar V4, la base de datos contendrá:

**Chevrolet** (22 modelos):
- Aveo, Onix, Spark EV, Groove, Captiva, Tracker, Trax, Equinox EV, Blazer EV, Traverse, Blazer, Tahoe, Suburban, Camaro, Corvette, S10, Tornado Van, Montana, Silverado, Colorado, Cheyenne, Brightdrop

**Buick** (4 modelos):
- Envista, Encore, Envision, Enclave

**GMC** (6 modelos):
- Terrain, Acadia, Yukon, Canyon, Sierra, Hummer

**Cadillac** (7 modelos):
- XT4, XT5, Optiq, Lyriq, Escalade, Escalade V, Escalade IQ

**Total**: 4 marcas, 39 modelos de vehículos

## Notas Importantes

- Las migraciones deben ejecutarse en orden (V1 -> V2 -> V3 -> V4)
- V4 es idempotente - si las marcas ya existen, fallará (usa ON CONFLICT si necesitas re-ejecutar)
- Los índices están optimizados para consultas por marca, año y estado
- La tabla `users` tiene una constraint que requiere `brand_id` para usuarios con rol BRAND
