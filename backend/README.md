# CustomerMX Backend

Backend API para el sistema de gestión de eventos CustomerMX.

## Stack Tecnológico

- **Go 1.24+**
- **Chi Router** - HTTP router ligero y performante
- **PostgreSQL** - Base de datos relacional
- **JWT** - Autenticación basada en tokens
- **pgx** - Driver PostgreSQL de alto rendimiento

## Estructura del Proyecto

```
backend/
├── cmd/
│   └── api/              # Punto de entrada de la aplicación
├── internal/
│   ├── config/           # Configuración de la aplicación
│   ├── domain/           # Lógica de negocio por dominio
│   │   ├── user/
│   │   ├── brand/
│   │   ├── vehicle/
│   │   ├── event/
│   │   ├── invitation/
│   │   └── notification/
│   ├── http/             # Capa HTTP
│   │   ├── handlers/     # Controladores HTTP
│   │   ├── middleware/   # Middlewares
│   │   └── router/       # Configuración de rutas
│   └── infra/            # Infraestructura
│       ├── db/           # Conexión a base de datos
│       ├── security/     # JWT, password hashing
│       └── mail/         # Servicio de email
├── migrations/           # Migraciones de base de datos
├── .env.example          # Ejemplo de variables de entorno
├── docker-compose.yml    # Docker para desarrollo local
├── Makefile              # Comandos útiles
└── go.mod                # Dependencias Go
```

## Configuración Inicial

### 1. Clonar variables de entorno

```bash
cp .env.example .env
```

Edita `.env` y configura tus variables, especialmente:
- `JWT_ACCESS_SECRET` - Cambiar en producción
- `JWT_REFRESH_SECRET` - Cambiar en producción
- `DB_PASSWORD` - Contraseña de PostgreSQL

### 2. Instalar dependencias

```bash
make deps
```

### 3. Iniciar PostgreSQL con Docker

```bash
make docker-up
```

Esto iniciará PostgreSQL en `localhost:5432` con:
- Usuario: `postgres`
- Contraseña: `postgres`
- Base de datos: `customermx`

### 4. Ejecutar migraciones

```bash
make migrate-up
```

## Ejecutar la Aplicación

### Modo desarrollo (sin hot reload)

```bash
make run
```

### Modo desarrollo (con hot reload)

Primero instala Air:

```bash
make install-tools
```

Luego ejecuta:

```bash
make watch
```

La aplicación estará disponible en `http://localhost:8080`

### Compilar para producción

```bash
make build
```

El binario se generará en `bin/api`

## Endpoints Disponibles

### Health Check

```bash
GET /health
```

### API v1 (Próximamente)

- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/invitations/accept` - Aceptar invitación
- Y más...

## Comandos Útiles

```bash
make help              # Ver todos los comandos disponibles
make run               # Ejecutar la aplicación
make build             # Compilar la aplicación
make test              # Ejecutar tests
make test-coverage     # Tests con coverage
make clean             # Limpiar artifacts
make deps              # Descargar dependencias
make lint              # Ejecutar linter
make docker-up         # Iniciar containers Docker
make docker-down       # Detener containers Docker
make docker-logs       # Ver logs de Docker
```

## Testing

```bash
# Ejecutar todos los tests
make test

# Tests con coverage
make test-coverage
```

## Desarrollo

### Agregar una nueva migración

```bash
make migrate-create name=create_users_table
```

### Estructura de un dominio

Cada dominio (user, brand, event, etc.) sigue esta estructura:

```
domain/user/
├── model.go       # Estructuras de datos
├── repository.go  # Acceso a base de datos
├── service.go     # Lógica de negocio
└── types.go       # DTOs y tipos auxiliares
```

## Próximos Pasos - Fase 1

- [ ] Crear migraciones de base de datos (V1-V3)
- [ ] Implementar dominio de usuarios
- [ ] Implementar dominio de marcas y vehículos
- [ ] Implementar dominio de invitaciones
- [ ] Crear handlers de autenticación
- [ ] Agregar tests unitarios
- [ ] Documentación de API con Swagger

## Licencia

Privado - CustomerMX
