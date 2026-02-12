# CustomerMX - Fase 1 Completada ✅

**Fecha de Completación:** 12 de Febrero, 2026
**Status:** ✅ COMPLETADA (Backend)

---

## Objetivo de la Fase 1

Establecer la base y estructura completa del backend CustomerMX, incluyendo:
- Base de datos PostgreSQL con migraciones
- Modelos de dominio en Go
- Sistema de autenticación JWT
- Middleware de permisos
- Endpoints básicos de API REST

---

## 1. Infraestructura y Configuración

### 1.1 Stack Tecnológico Implementado

- **Lenguaje:** Go 1.24+
- **Router:** Chi v5 (HTTP router ligero y performante)
- **Base de Datos:** PostgreSQL 15 (Docker)
- **Driver DB:** pgx v5 (driver PostgreSQL de alto rendimiento)
- **Autenticación:** JWT (golang-jwt/jwt v5)
- **Password:** bcrypt (golang.org/x/crypto)
- **CORS:** go-chi/cors
- **Env Management:** godotenv

### 1.2 Estructura del Proyecto

```
backend/
├── cmd/
│   ├── api/              # Servidor HTTP principal
│   │   └── main.go
│   ├── migrate/          # Script de migraciones
│   │   └── main.go
│   └── seed/             # Script para crear admin inicial
│       └── main.go
├── internal/
│   ├── config/           # Configuración de la aplicación
│   │   └── config.go
│   ├── domain/           # Lógica de negocio por dominio
│   │   ├── user/         # ✅ Completo
│   │   │   ├── model.go
│   │   │   ├── errors.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── brand/        # ✅ Completo
│   │   │   ├── model.go
│   │   │   ├── errors.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── vehicle/      # ✅ Completo
│   │   │   ├── model.go
│   │   │   ├── errors.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── invitation/   # ✅ Completo
│   │   │   ├── model.go
│   │   │   ├── errors.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── event/        # ⏳ Pendiente (Fase 2)
│   │   └── notification/ # ⏳ Pendiente (Fase 4)
│   ├── http/
│   │   ├── handlers/     # Controladores HTTP
│   │   │   ├── response.go
│   │   │   ├── auth.go
│   │   │   ├── user.go
│   │   │   ├── brand.go
│   │   │   ├── vehicle.go
│   │   │   └── invitation.go
│   │   ├── middleware/   # Middlewares
│   │   │   └── auth.go
│   │   └── router/       # Configuración de rutas
│   │       └── router.go
│   └── infra/            # Infraestructura
│       ├── db/
│       │   └── connection.go
│       ├── security/
│       │   ├── jwt.go
│       │   └── password.go
│       └── mail/
│           └── service.go
├── migrations/           # Migraciones SQL
│   ├── V1__create_brands_and_vehicles.sql
│   ├── V2__create_users_and_invitations.sql
│   ├── V3__create_events_and_reports.sql
│   ├── V4__seed_brands_and_vehicles.sql
│   ├── README.md
│   └── run_migrations.sh
├── .env                  # Variables de entorno (git ignored)
├── .env.example          # Ejemplo de configuración
├── .gitignore
├── .air.toml             # Configuración hot reload
├── docker-compose.yml    # PostgreSQL para desarrollo
├── Makefile              # Comandos útiles
├── README.md
├── test_api.sh           # Script de pruebas API
└── go.mod
```

---

## 2. Base de Datos PostgreSQL

### 2.1 Migraciones Implementadas

#### V1 - Brands y Vehicles
- Tabla `brands` - Marcas automotrices
- Tabla `vehicles` - Modelos de vehículos asociados a marcas
- Constraint: Modelo único por marca
- Índice en `vehicles.brand_id`

#### V2 - Users e Invitations
- Enum `user_role` (ADMIN, COORDINATOR, BRAND)
- Tabla `users` - Usuarios del sistema con roles
- Tabla `invitations` - Invitaciones pendientes
- Constraints:
  - Usuarios BRAND deben tener `brand_id`
  - Usuarios no-BRAND no pueden tener `brand_id`
- Índices en email, role, brand_id

#### V3 - Events, Reports y Notifications
- Enum `event_status` (PLANNED, ACTIVE, COMPLETED, CLOSED)
- Tabla `events` - Eventos promocionales
- Tabla `event_coordinators` - Asignación coordinadores a eventos
- Tabla `event_reports` - Reportes operativos de eventos
- Tabla `event_vehicles` - Vehículos presentados en eventos
- Tabla `notifications` - Sistema de notificaciones
- Índices optimizados para consultas por marca, año, estado

#### V4 - Seed Data
- 4 Marcas: Chevrolet, Buick, GMC, Cadillac
- 39 Modelos de vehículos distribuidos por marca
  - Chevrolet: 22 modelos
  - Buick: 4 modelos
  - GMC: 6 modelos
  - Cadillac: 7 modelos

### 2.2 Tablas Creadas

Total: **9 tablas**

| Tabla | Descripción | Registros Iniciales |
|-------|-------------|---------------------|
| `brands` | Marcas automotrices | 4 |
| `vehicles` | Modelos de vehículos | 39 |
| `users` | Usuarios del sistema | 1 (admin) |
| `invitations` | Invitaciones pendientes | 0 |
| `events` | Eventos promocionales | 0 |
| `event_coordinators` | Asignación coordinadores | 0 |
| `event_reports` | Reportes de eventos | 0 |
| `event_vehicles` | Vehículos en eventos | 0 |
| `notifications` | Notificaciones usuarios | 0 |

---

## 3. Dominios Implementados

### 3.1 User Domain

**Responsabilidad:** Gestión de usuarios y autenticación

**Modelos:**
- `User` - Usuario del sistema
- `CreateUserRequest` - Solicitud de creación
- `UpdateUserRequest` - Solicitud de actualización
- `LoginRequest` - Credenciales de login
- `LoginResponse` - Respuesta con tokens JWT
- `UserResponse` - Usuario sin datos sensibles

**Roles:**
- `ADMIN` - Acceso completo al sistema
- `COORDINATOR` - Gestión de eventos asignados
- `BRAND` - Visualización de datos de su marca

**Funcionalidades:**
- ✅ Login con email/password
- ✅ Generación de tokens JWT (access + refresh)
- ✅ CRUD completo de usuarios
- ✅ Listado por rol
- ✅ Listado por marca
- ✅ Activación de cuentas
- ✅ Validación de usuarios activos/inactivos

**Reglas de Negocio:**
- Usuarios BRAND deben tener `brand_id`
- Usuarios ADMIN y COORDINATOR no tienen `brand_id`
- Passwords hasheados con bcrypt
- Email único en el sistema

### 3.2 Brand Domain

**Responsabilidad:** Gestión de marcas automotrices

**Modelos:**
- `Brand` - Marca automotriz
- `CreateBrandRequest` - Solicitud de creación
- `UpdateBrandRequest` - Solicitud de actualización

**Funcionalidades:**
- ✅ CRUD completo de marcas
- ✅ Búsqueda por nombre
- ✅ Listado ordenado alfabéticamente

**Reglas de Negocio:**
- Nombres de marca únicos
- Eliminación en cascada de vehículos asociados

### 3.3 Vehicle Domain

**Responsabilidad:** Gestión de modelos de vehículos

**Modelos:**
- `Vehicle` - Vehículo
- `VehicleWithBrand` - Vehículo con info de marca (JOIN)
- `CreateVehicleRequest` - Solicitud de creación
- `UpdateVehicleRequest` - Solicitud de actualización

**Funcionalidades:**
- ✅ CRUD completo de vehículos
- ✅ Listado con información de marca
- ✅ Listado filtrado por marca
- ✅ Búsqueda por ID

**Reglas de Negocio:**
- Modelo único por marca (no duplicados)
- Debe estar asociado a una marca existente

### 3.4 Invitation Domain

**Responsabilidad:** Sistema de invitaciones de usuarios

**Modelos:**
- `Invitation` - Invitación
- `CreateInvitationRequest` - Solicitud de creación
- `AcceptInvitationRequest` - Solicitud de aceptación

**Funcionalidades:**
- ✅ Creación de invitaciones con token único
- ✅ Envío de email (mock en desarrollo)
- ✅ Aceptación de invitaciones
- ✅ Creación automática de usuario al aceptar
- ✅ Reenvío de invitaciones
- ✅ Expiración de 7 días
- ✅ Listado de invitaciones pendientes

**Reglas de Negocio:**
- Token seguro de 64 caracteres (hex)
- No se puede invitar a emails ya registrados
- Invitaciones BRAND requieren `brand_id`
- Una vez aceptada, no se puede re-aceptar
- Validación de expiración

---

## 4. API REST Implementada

### 4.1 Endpoints Públicos (2)

**Sin autenticación requerida**

```http
POST /api/v1/auth/login
POST /api/v1/invitations/accept
```

### 4.2 Endpoints Protegidos (34)

**Requieren JWT token en header `Authorization: Bearer <token>`**

#### Auth (3 endpoints)
```http
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
GET    /api/v1/auth/me
```

#### Users (6 endpoints)
```http
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}
DELETE /api/v1/users/{id}
GET    /api/v1/users/role/{role}
```

#### Brands (5 endpoints)
```http
GET    /api/v1/brands
POST   /api/v1/brands
GET    /api/v1/brands/{id}
PUT    /api/v1/brands/{id}
DELETE /api/v1/brands/{id}
```

#### Vehicles (6 endpoints)
```http
GET    /api/v1/vehicles
POST   /api/v1/vehicles
GET    /api/v1/vehicles/{id}
PUT    /api/v1/vehicles/{id}
DELETE /api/v1/vehicles/{id}
GET    /api/v1/brands/{brandId}/vehicles
```

#### Invitations (5 endpoints)
```http
GET    /api/v1/invitations
POST   /api/v1/invitations
GET    /api/v1/invitations/{id}
POST   /api/v1/invitations/{id}/resend
DELETE /api/v1/invitations/{id}
```

#### Health Check (1 endpoint)
```http
GET    /health
```

**Total: 36 endpoints**

### 4.3 Formato de Respuestas

#### Éxito
```json
{
  "data": { ... },
  "message": "optional message"
}
```

#### Error
```json
{
  "error": "HTTP Status Text",
  "message": "Detailed error message"
}
```

---

## 5. Seguridad Implementada

### 5.1 Autenticación JWT

**Configuración:**
- Algoritmo: HS256
- Access Token TTL: 15 minutos (configurable)
- Refresh Token TTL: 7 días (configurable)

**Payload de Claims:**
```json
{
  "sub": "user-uuid",
  "role": "ADMIN|COORDINATOR|BRAND",
  "brand_id": "brand-uuid (opcional)",
  "is_active": true,
  "iat": 1234567890,
  "exp": 1234567890
}
```

**Flujo de Autenticación:**
1. Login con email/password
2. Validación de credenciales + bcrypt
3. Generación de token pair (access + refresh)
4. Cliente guarda tokens
5. Requests con `Authorization: Bearer <access_token>`
6. Refresh cuando access token expira

### 5.2 Password Hashing

- Algoritmo: bcrypt
- Cost factor: DefaultCost (10)
- Passwords nunca expuestos en JSON
- Verificación segura con timing attack protection

### 5.3 Middleware de Autenticación

**Ubicación:** `internal/http/middleware/auth.go`

**Funcionalidad:**
- Extracción de token del header Authorization
- Validación de formato "Bearer <token>"
- Verificación de firma JWT
- Validación de expiración
- Verificación de usuario activo
- Inyección de claims en contexto

**Uso:**
```go
r.Use(middleware.AuthMiddleware(cfg))
```

### 5.4 CORS

**Configuración:**
- Orígenes permitidos: `http://localhost:5173`, `http://localhost:3000`
- Métodos: GET, POST, PUT, DELETE, OPTIONS
- Headers permitidos: Accept, Authorization, Content-Type
- Credentials: true
- Max Age: 300s

---

## 6. Comandos y Scripts

### 6.1 Makefile

```bash
# Docker
make docker-up        # Iniciar PostgreSQL
make docker-down      # Detener PostgreSQL
make docker-logs      # Ver logs de Docker

# Migraciones
make migrate-up       # Ejecutar migraciones
make migrate-status   # Ver tablas creadas
make migrate-verify   # Verificar seed data
make seed             # Crear usuario admin inicial

# Base de datos
make db-shell         # Abrir shell PostgreSQL
make db-reset         # Resetear BD (¡DESTRUCTIVO!)

# Desarrollo
make run              # Ejecutar API
make build            # Compilar binario
make watch            # Hot reload (requiere Air)

# Testing
make test             # Ejecutar tests
make test-coverage    # Tests con coverage

# Otros
make deps             # Descargar dependencias
make clean            # Limpiar artifacts
make lint             # Ejecutar linter
make install-tools    # Instalar herramientas dev
```

### 6.2 Scripts Personalizados

**test_api.sh** - Script de pruebas completo
```bash
./test_api.sh
```

Prueba:
1. Health endpoint
2. Login
3. /auth/me con token
4. Listado de marcas

---

## 7. Configuración

### 7.1 Variables de Entorno (.env)

```bash
# Server
SERVER_PORT=8080
SERVER_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=customermx
DB_SSL_MODE=disable

# JWT
JWT_ACCESS_SECRET=your-super-secret-access-key
JWT_REFRESH_SECRET=your-super-secret-refresh-key
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# Email
EMAIL_PROVIDER=mock
EMAIL_FROM=noreply@customermx.com
FRONTEND_URL=http://localhost:5173

# SMTP (opcional)
SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_USE_TLS=true
```

### 7.2 Docker Compose

PostgreSQL en contenedor:
```yaml
services:
  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: customermx
```

---

## 8. Credenciales de Acceso

### Usuario Administrador Inicial

```
Email:    admin@customermx.com
Password: admin123
Role:     ADMIN
```

**Generado automáticamente con:** `make seed`

---

## 9. Pruebas Exitosas

### 9.1 Health Check
```bash
curl http://localhost:8080/health
```
**Status:** ✅ 200 OK

### 9.2 Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@customermx.com","password":"admin123"}'
```
**Status:** ✅ 200 OK
**Response:** User + access_token + refresh_token

### 9.3 Usuario Actual
```bash
curl http://localhost:8080/api/v1/auth/me \
  -H 'Authorization: Bearer <token>'
```
**Status:** ✅ 200 OK

### 9.4 Listar Marcas
```bash
curl http://localhost:8080/api/v1/brands \
  -H 'Authorization: Bearer <token>'
```
**Status:** ✅ 200 OK
**Response:** 4 marcas

---

## 10. Problemas Resueltos

### 10.1 Escaneo de Enums PostgreSQL

**Problema:** Error 500 al hacer login debido a que pgx no podía escanear el enum `user_role` directamente a un tipo custom de Go.

**Solución:** Escanear primero a `string` y luego convertir al tipo `Role`:
```go
var roleStr string
err := r.pool.QueryRow(ctx, query, email).Scan(
    &user.ID, &user.Name, &user.Email, &user.PasswordHash,
    &roleStr, // Escanear a string primero
    &user.BrandID, &user.IsActive,
    &user.CreatedAt, &user.UpdatedAt,
)
user.Role = Role(roleStr) // Convertir después
```

**Archivos afectados:**
- `internal/domain/user/repository.go` (todos los métodos de lectura)
- `internal/domain/invitation/repository.go` (métodos que leen users)

---

## 11. Decisiones de Diseño

### 11.1 Arquitectura Hexagonal / Clean Architecture

- **Domain** - Lógica de negocio pura
- **Repository** - Acceso a datos (PostgreSQL)
- **Service** - Orquestación y lógica de aplicación
- **Handler** - Capa HTTP (presentación)

### 11.2 Dependency Injection

Todos los servicios y repos se inicializan en el router y se inyectan a los handlers.

### 11.3 Error Handling

- Errores de dominio definidos como constantes
- Errores de infraestructura wrapeados con contexto
- Respuestas HTTP consistentes con status codes apropiados

### 11.4 Validación

- Validación en capa de dominio (métodos Validate())
- No se confía en validación de cliente
- Mensajes de error claros y específicos

---

## 12. Métricas del Proyecto

### 12.1 Líneas de Código

- **Dominio:** ~1,200 líneas
- **Handlers:** ~800 líneas
- **Infraestructura:** ~400 líneas
- **Migraciones:** ~300 líneas SQL
- **Total:** ~2,700 líneas

### 12.2 Archivos Creados

- Go: 28 archivos
- SQL: 4 archivos de migración
- Config: 6 archivos (Makefile, docker-compose, etc.)
- **Total:** 38 archivos

### 12.3 Tiempo de Desarrollo

- Fase 1: 1 sesión
- Debugging: Issue del enum resuelto

---

## 13. Dependencias de Go

```go
require (
    github.com/go-chi/chi/v5 v5.2.5
    github.com/go-chi/cors v1.2.2
    github.com/golang-jwt/jwt/v5 v5.3.1
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.8.0
    github.com/joho/godotenv v1.5.1
    golang.org/x/crypto v0.48.0
)
```

---

## 14. Próximos Pasos (Fase 2)

### Prioridad Alta
- [ ] **Dominio Event** - CRUD de eventos
- [ ] **Asignación de Coordinadores** - Relación events <-> coordinators
- [ ] **Vehículos en Eventos** - Relación events <-> vehicles
- [ ] **Reportes de Eventos** - event_reports CRUD

### Prioridad Media
- [ ] **Middleware de Permisos por Rol** - Admin-only endpoints
- [ ] **Brand Isolation** - Usuarios BRAND solo ven sus datos
- [ ] **Coordinator Isolation** - Solo eventos asignados

### Prioridad Baja
- [ ] **Tests Unitarios** - Cobertura >80%
- [ ] **Tests de Integración** - E2E testing
- [ ] **Documentación OpenAPI** - Swagger spec

---

## 15. Notas Finales

### 15.1 Estado del Proyecto

✅ **Fase 1 Backend - COMPLETADA**

El backend está completamente funcional y listo para:
- Desarrollo del frontend
- Implementación de la Fase 2 (Events)
- Integración continua

### 15.2 Calidad del Código

- ✅ Compila sin warnings
- ✅ Sin código duplicado significativo
- ✅ Nombres descriptivos y consistentes
- ✅ Separación clara de responsabilidades
- ✅ Error handling robusto

### 15.3 Lecciones Aprendidas

1. **Enums de PostgreSQL:** Siempre escanear a string primero cuando uses tipos custom en Go
2. **Validación Temprana:** La validación en el dominio previene errores downstream
3. **Dependency Injection:** Facilita testing y mantenibilidad
4. **Migraciones Versionadas:** Critical para mantener consistencia de BD

---

## 16. Referencias

- **Repositorio:** `/Users/josebeltran/Documents/GitHub/customermx`
- **Documentación Principal:** `README.md`
- **Migraciones:** `backend/migrations/`
- **Postman Collection:** ⏳ Pendiente

---

**Documento Creado:** 12 de Febrero, 2026
**Última Actualización:** 12 de Febrero, 2026
**Versión:** 1.0
**Estado:** ✅ FASE 1 COMPLETADA
