# CustomerMX — Documentación Completa

---

## 1. Visión general

**CustomerMX** es un sistema de gestión de eventos enfocado en la promoción de vehículos de distintas **marcas automotrices**. Permite:

- planificar eventos,
- registrar información operativa durante su ejecución,
- consultar resultados y métricas posteriores.

El sistema está diseñado para manejar múltiples marcas, sus vehículos y la operación de eventos asociados.

### Stack tecnológico

- **Frontend:** React + TypeScript
- **Backend:** Go
- **Base de datos:** PostgreSQL

### Scripts de gestión del entorno

Para facilitar el desarrollo, el proyecto incluye scripts de gestión:

```bash
# Iniciar todos los servicios (BD, Backend, Frontend)
./runenv.sh

# Detener todos los servicios
./stopenv.sh

# Reiniciar todos los servicios
./restartenv.sh
```

**Nota:** `runenv.sh` automáticamente detiene los servicios antes de iniciarlos para evitar conflictos.

Los logs se guardan en `logs/`:
- `logs/backend.log` - Logs del backend
- `logs/frontend.log` - Logs del frontend

URLs de acceso:
- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **Base de datos:** localhost:5432

---

## 2. Roles del sistema

El sistema cuenta con tres roles principales:

1. **Administrador General**
2. **Coordinador**
3. **Marca Automotriz** (rol: BRAND)

---

## 3. Restricción de asociación — Rol Marca Automotriz

- Todo usuario con rol **Marca** debe estar asociado a una marca específica.
- Una marca puede tener múltiples usuarios.
- Cada usuario Marca pertenece a **una sola marca**.
- Solo puede visualizar información de su marca.

Ejemplo:

- VW → 3 usuarios Marca  
- Nissan → 2 usuarios Marca

---

## 4. Funcionalidades por rol

### 4.1 Administrador General

- Invitar usuarios al sistema (Coordinador o Marca)  
- Activar / desactivar usuarios  
- Registrar marcas y vehículos  
- Crear, editar, eliminar eventos  
- Asignar coordinadores a eventos  
- Consultar eventos pasados y futuros  
- Captura de información durante y post evento

### 4.2 Coordinador

- Recibir invitación por correo  
- Recibir notificación de asignación a eventos  
- Consultar eventos asignados  
- Registrar información durante el evento

### 4.3 Marca Automotriz

- Recibir invitación por correo  
- Consultar eventos donde participaron sus vehículos  
- Consultar eventos pasados y futuros  
- Visualizar gráficas y filtros dinámicos  
- Recibir notificación cuando el evento esté completamente registrado  
- Exportar información a PDF

---

## 5. Datos del evento

### 5.1 Programación del evento

Campos:

- Marca (brand_id)  
- Tipo de evento (Triatlón, Fútbol, Golf…)  
- Organizador  
- Nombre del evento  
- Fecha de inicio  
- Año  
- Duración (días)  
- Estado  
- Ciudad  
- Sede  
- Distribuidor

---

### 5.2 Información durante / post evento

- Número de edecanes  
- Responsable del montaje  
- Promocional (Sí / No)  
- Número de asistentes  
- Número de dinámicas (puede ser N/A)  
- Datos levantados  
- Prospectos  
- Calificación del distribuidor (1–5)  
- Comentarios

---

### 5.3 Vehículos presentados

Cada evento está ligado a **una marca**.

- Solo vehículos de esa marca  
- Registrar modelo y cantidad

Ejemplo VW:

- Jetta → 3  
- Tiguan → 2  
- Taos → 1

---

### 5.4 Datos predefinidos (seed data)

**Marcas y vehículos iniciales:**

- **Chevrolet:** Aveo, Onix, Spark EV, Groove, Captiva, Tracker, Trax, Equinox EV, Blazer EV, Traverse, Blazer, Tahoe, Suburban, Camaro, Corvette, S10, Tornado Van, Montana, Silverado, Colorado, Cheyenne, Brightdrop  
- **Buick:** Envista, Encore, Envision, Enclave  
- **GMC:** Terrain, Acadia, Yukon, Canyon, Sierra, Hummer  
- **Cadillac:** XT4, XT5, Optiq, Lyriq, Escalade, Escalade V, Escalade IQ

---

## 6. Estructura de Base de Datos (PostgreSQL)

### 6.1 Marcas

```sql
CREATE TABLE brands (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now()
);
```

### 6.2 Vehículos

```sql
CREATE TABLE vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    model_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(brand_id, model_name)
);

CREATE INDEX idx_vehicles_brand ON vehicles(brand_id);
```

### 6.3 Usuarios

```sql
CREATE TYPE user_role AS ENUM ('ADMIN', 'COORDINATOR', 'BRAND');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL,
    brand_id UUID REFERENCES brands(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    CHECK (
        (role = 'BRAND' AND brand_id IS NOT NULL)
        OR
        (role <> 'BRAND' AND brand_id IS NULL)
    )
);

CREATE INDEX idx_users_brand ON users(brand_id);
```

### 6.4 Invitaciones

```sql
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL,
    role user_role NOT NULL,
    brand_id UUID REFERENCES brands(id),
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_invitations_email ON invitations(email);
```

### 6.5 Eventos

```sql
CREATE TYPE event_status AS ENUM ('PLANNED','ACTIVE','COMPLETED','CLOSED');

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand_id UUID NOT NULL REFERENCES brands(id),
    event_type TEXT NOT NULL,
    organizer TEXT NOT NULL,
    name TEXT NOT NULL,
    start_date DATE NOT NULL,
    year INT NOT NULL,
    duration_days INT NOT NULL,
    state TEXT NOT NULL,
    city TEXT NOT NULL,
    venue TEXT NOT NULL,
    dealer TEXT NOT NULL,
    status event_status DEFAULT 'PLANNED',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

### Índices críticos

```sql
CREATE INDEX idx_events_brand ON events(brand_id);
CREATE INDEX idx_events_year ON events(year);
CREATE INDEX idx_events_date ON events(start_date);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_state_city ON events(state, city);
CREATE INDEX idx_events_brand_year ON events(brand_id, year);
CREATE INDEX idx_events_brand_date ON events(brand_id, start_date);
CREATE INDEX idx_events_brand_status ON events(brand_id, status);
```

### 6.6 Coordinadores de eventos

```sql
CREATE TABLE event_coordinators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(event_id, user_id)
);

CREATE INDEX idx_event_coord_event ON event_coordinators(event_id);
CREATE INDEX idx_event_coord_user ON event_coordinators(user_id);
```

### 6.7 Reportes de eventos

```sql
CREATE TABLE event_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL UNIQUE REFERENCES events(id) ON DELETE CASCADE,
    hostess_count INT,
    setup_vendor TEXT,
    has_promotional BOOLEAN,
    attendees INT,
    activities_count INT,
    leads_collected INT,
    prospects INT,
    dealer_rating INT CHECK (dealer_rating BETWEEN 1 AND 5),
    comments TEXT,
    completed BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_event_reports_event ON event_reports(event_id);
CREATE INDEX idx_event_reports_rating ON event_reports(dealer_rating);
```

### 6.8 Vehículos presentados en eventos

```sql
CREATE TABLE event_vehicles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    UNIQUE(event_id, vehicle_id)
);

CREATE INDEX idx_event_vehicles_event ON event_vehicles(event_id);
CREATE INDEX idx_event_vehicles_vehicle ON event_vehicles(vehicle_id);
```

### 6.9 Notificaciones

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    payload JSONB,
    sent_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
```

---

## 7. APIs

[Se mantienen todos los endpoints con `brand_id` en lugar de `company_id`]

> TODO: materialized views para dashboards por marca y año

---

## 8. Middleware de permisos

- JWT auth + role  
- Brand isolation (solo datos de su marca)  
- Coordinator → solo eventos asignados  
- Admin → acceso total  
- Endpoints sensibles retornan 401/403/404 según contexto

> TODO: materialized views respeten permisos JWT y brand_id

---

## 9. Seguridad JWT

- Payload: `{sub, role, brand_id, is_active, iat, exp}`  
- Firma: HS256 / EdDSA  
- Expiración corta (15–60 min)  
- Refresh token opcional  
- Validación DB obligatoria  
- HTTPS obligatorio  
- Logout → revocar / incrementar token_version

> TODO: dashboards/materialized views deben validar JWT y brand_id

---

## 10. Reglas de negocio críticas

- Eventos siempre ligados a **una marca**  
- Vehículos solo de esa marca  
- Coordinadores solo en eventos asignados  
- Usuarios marca solo ven datos de su marca  
- Todos los filtros y dashboards optimizados para lectura rápida

---

## 11. Implementación — Roadmap por fases

### 11.1 Fase 1 — Base y estructura
- Configuración repositorios  
- Flyway y PostgreSQL  
- Migraciones V1–V3 con seed de marcas y vehículos  
- Modelos en Go y TypeScript  
- Autenticación JWT  
- Middleware base de permisos  
- Endpoints login y aceptación de invitación

### 11.2 Fase 2 — Gestión de eventos
- CRUD eventos  
- Asignación coordinadores  
- Registro vehículos y reportes  
- Validación de permisos por rol y brand_id  
- Middleware de visibilidad de eventos

### 11.3 Fase 3 — Analytics y exportación
- Endpoints dashboards `/analytics/events`  
- Exportación PDF  
- Materialized views iniciales (TODO)  
- Filtros dinámicos y gráficas frontend

### 11.4 Fase 4 — Notificaciones y flujo de invitaciones
- Envío correos de invitación  
- Notificaciones de asignación y reporte completado  
- Endpoint historial de notificaciones

### 11.5 Fase 5 — Optimización y seguridad
- Optimización consultas e índices  
- Revisión middleware y permisos  
- Seguridad JWT y refresh  
- Auditoría de accesos  
- Pruebas de stress en consultas y dashboards

### 11.6 Fase 6 — Producción y mantenimiento
- CI/CD backend/frontend  
- Migraciones Flyway producción  
- Backups periódicos  
- Documentación completa  
- Monitoreo de errores y métricas de uso

---

> TODO: crear migraciones para materialized views optimizadas para dashboards por marca y año.
