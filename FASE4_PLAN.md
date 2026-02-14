# FASE 4 - Admin Panel y Gestión de Usuarios

**Estado:** 🚧 **PLANEADA**
**Fecha:** Febrero 2026

---

## 📋 Resumen

La Fase 4 completa el Admin Panel con gestión completa de usuarios e invitaciones. Actualmente solo existe el login, pero no hay UI para que los ADMIN puedan:
- Ver lista de usuarios
- Crear nuevos usuarios
- Editar usuarios existentes
- Asignar/cambiar roles
- Desactivar usuarios
- Gestionar invitaciones (crear, reenviar, ver pendientes)

**Backend:** Ya existe completo desde Fase 1 (User domain, Invitation domain)
**Frontend:** Crear 3 páginas nuevas + componentes

---

## 🎯 Alcance de la Implementación

### Páginas a Crear (3)

1. **Users List** (`/users`) - Lista de usuarios con filtros
2. **User Form** (`/users/new`, `/users/:id/edit`) - Crear/editar usuario
3. **Invitations** (`/invitations`) - Gestión de invitaciones

### Componentes a Crear (5)

1. **ChangeRoleDialog** - Cambiar rol de usuario
2. **DeactivateUserDialog** - Desactivar/activar usuario
3. **InviteUserDialog** - Enviar invitación
4. **ResendInvitationDialog** - Reenviar invitación
5. **UserDetailsDialog** - Ver detalles completos del usuario

### Redux State Management

- **usersSlice** - Estado de usuarios
- **invitationsSlice** - Estado de invitaciones

---

## 📦 Estructura de Archivos

```
frontend/src/
├── features/
│   ├── users/
│   │   └── usersSlice.ts               [NUEVO]
│   └── invitations/
│       └── invitationsSlice.ts         [NUEVO]
├── pages/
│   ├── Users/
│   │   ├── index.tsx                   [NUEVO] - Lista de usuarios
│   │   ├── UserForm.tsx                [NUEVO] - Crear/editar
│   │   └── UserDetail.tsx              [NUEVO] - Detalle (opcional)
│   └── Invitations/
│       └── index.tsx                   [NUEVO] - Gestión de invitaciones
├── components/
│   └── Dialogs/
│       ├── ChangeRoleDialog.tsx        [NUEVO]
│       ├── DeactivateUserDialog.tsx    [NUEVO]
│       ├── InviteUserDialog.tsx        [NUEVO]
│       └── ResendInvitationDialog.tsx  [NUEVO]
├── App.tsx                             [MODIFICAR] - Agregar rutas
└── components/Layout/index.tsx         [MODIFICAR] - Agregar items menú
```

**Total:** 11 archivos (9 nuevos, 2 modificados)

---

## 🎨 Diseño de Páginas

### 1. Users List (`/users`)

**Layout:**
```
┌─────────────────────────────────────────────────┐
│ 👥 Usuarios                    [+ Nuevo Usuario]│
├─────────────────────────────────────────────────┤
│ 🔍 Filtros                                      │
│ [Rol ▼] [Estado ▼] [Buscar...]                 │
├─────────────────────────────────────────────────┤
│ Tabla:                                          │
│ │Nombre│Email│Rol│Marca│Estado│Creado│Acciones││
│ │Juan  │j@..  │ADM│-    │Activo│01/01│ ⚙️ 🗑️  ││
│ │María │m@..  │BRA│Chev │Activo│01/02│ ⚙️ 🗑️  ││
│ │Pedro │p@..  │COO│-    │Inact │01/03│ ⚙️ ✅  ││
└─────────────────────────────────────────────────┘
```

**Funcionalidades:**
- ✅ Tabla con todos los usuarios (solo ADMIN puede ver)
- ✅ Filtros: Rol (ADMIN/COORDINATOR/BRAND), Estado (Activo/Inactivo)
- ✅ Búsqueda por nombre o email
- ✅ Chips de colores para roles
- ✅ Indicador visual de estado activo/inactivo
- ✅ Acciones por fila:
  - **Editar** (pencil icon) → navega a `/users/:id/edit`
  - **Cambiar Rol** (swap icon) → abre ChangeRoleDialog
  - **Desactivar/Activar** (toggle) → abre DeactivateUserDialog
- ✅ Botón "Nuevo Usuario" → navega a `/users/new`
- ✅ Mostrar marca asociada si role=BRAND

**Estados:**
- Loading: Spinner centrado
- Error: Alert con mensaje
- Empty: Mensaje "No hay usuarios registrados"

---

### 2. User Form (`/users/new`, `/users/:id/edit`)

**Layout - Modo Create:**
```
┌─────────────────────────────────────────────────┐
│ ← Volver   Nuevo Usuario                        │
├─────────────────────────────────────────────────┤
│ ┌─ Información Básica ────────────────────────┐ │
│ │ [Nombre completo*]                          │ │
│ │ [Email*]                                    │ │
│ │ [Contraseña*]                               │ │
│ │ [Confirmar contraseña*]                     │ │
│ └─────────────────────────────────────────────┘ │
│                                                  │
│ ┌─ Rol y Permisos ─────────────────────────────┐│
│ │ [Rol* ▼] ADMIN / COORDINATOR / BRAND        ││
│ │                                              ││
│ │ Si Rol = BRAND:                             ││
│ │ [Marca* ▼] Chevrolet, Buick, GMC, Cadillac ││
│ └─────────────────────────────────────────────┘│
│                                                  │
│                    [Cancelar] [Crear Usuario]   │
└─────────────────────────────────────────────────┘
```

**Layout - Modo Edit:**
```
┌─────────────────────────────────────────────────┐
│ ← Volver   Editar Usuario: Juan Pérez           │
├─────────────────────────────────────────────────┤
│ ┌─ Información Básica ────────────────────────┐ │
│ │ [Nombre completo*]                          │ │
│ │ [Email*] (disabled - no se puede cambiar)  │ │
│ └─────────────────────────────────────────────┘ │
│                                                  │
│ ┌─ Cambiar Contraseña (opcional) ──────────────┐│
│ │ [Nueva contraseña]                          ││
│ │ [Confirmar nueva contraseña]                ││
│ │ 💡 Dejar en blanco para mantener actual     ││
│ └─────────────────────────────────────────────┘│
│                                                  │
│ ┌─ Rol y Permisos ─────────────────────────────┐│
│ │ [Rol* ▼] ADMIN / COORDINATOR / BRAND        ││
│ │                                              ││
│ │ Si Rol = BRAND:                             ││
│ │ [Marca* ▼] Chevrolet, Buick, GMC, Cadillac ││
│ └─────────────────────────────────────────────┘│
│                                                  │
│                    [Cancelar] [Guardar Cambios] │
└─────────────────────────────────────────────────┘
```

**Validaciones:**
- ✅ Email válido y único
- ✅ Contraseña mínimo 8 caracteres (solo en create)
- ✅ Confirmar contraseña coincide
- ✅ Si role=BRAND, marca es requerida
- ✅ Si role≠BRAND, limpiar brand_id

**Funcionalidades:**
- ✅ Create: Campos vacíos, contraseña requerida
- ✅ Edit: Pre-llenar datos, email disabled, contraseña opcional
- ✅ Select de rol con 3 opciones
- ✅ Select de marca (aparece solo si role=BRAND)
- ✅ Validación en frontend antes de submit
- ✅ API: POST `/users` (create), PUT `/users/:id` (update)

---

### 3. Invitations (`/invitations`)

**Layout:**
```
┌─────────────────────────────────────────────────┐
│ 📧 Invitaciones              [+ Enviar Invitación]│
├─────────────────────────────────────────────────┤
│ 🔍 Filtros                                      │
│ [Estado ▼] Todas / Pendientes / Aceptadas      │
├─────────────────────────────────────────────────┤
│ Tabs: [Pendientes (5)] [Aceptadas] [Todas]     │
├─────────────────────────────────────────────────┤
│ Lista:                                          │
│ ┌─────────────────────────────────────────────┐│
│ │ 📧 maria@example.com                        ││
│ │ Rol: COORDINATOR                            ││
│ │ Enviada: hace 2 días                        ││
│ │ Token: abc123...                            ││
│ │              [📤 Reenviar] [🗑️ Eliminar]    ││
│ └─────────────────────────────────────────────┘│
│ ┌─────────────────────────────────────────────┐│
│ │ ✅ pedro@example.com                        ││
│ │ Rol: BRAND - Chevrolet                     ││
│ │ Aceptada: hace 1 semana                    ││
│ └─────────────────────────────────────────────┘│
└─────────────────────────────────────────────────┘
```

**Funcionalidades:**
- ✅ Lista de invitaciones con estado visual
- ✅ Tabs para filtrar: Pendientes / Aceptadas / Todas
- ✅ Cards con información de invitación:
  - Email
  - Rol (y marca si aplica)
  - Fecha de envío
  - Token (truncado)
  - Estado (pendiente/aceptada)
- ✅ Acciones por invitación pendiente:
  - **Reenviar** → ResendInvitationDialog
  - **Eliminar** → ConfirmDialog
- ✅ Invitaciones aceptadas solo lectura
- ✅ Botón "Enviar Invitación" → InviteUserDialog
- ✅ Empty states por tab

---

## 🎭 Dialogs Nuevos

### 1. ChangeRoleDialog

**Props:**
```typescript
interface ChangeRoleDialogProps {
  open: boolean;
  onClose: () => void;
  user: User;
  onSuccess: () => void;
}
```

**UI:**
```
┌────────────────────────────────────┐
│ Cambiar Rol de Usuario             │
├────────────────────────────────────┤
│ Usuario: Juan Pérez (j@example.com)│
│ Rol actual: COORDINATOR            │
│                                    │
│ [Nuevo Rol* ▼]                     │
│   • ADMIN                          │
│   • COORDINATOR                    │
│   • BRAND                          │
│                                    │
│ Si selecciona BRAND:               │
│ [Marca* ▼] Chevrolet               │
│                                    │
│ ⚠️ El usuario será notificado      │
│    del cambio de rol               │
│                                    │
│        [Cancelar] [Cambiar Rol]    │
└────────────────────────────────────┘
```

**API:** `PUT /users/:id` (actualizar role y brand_id)

---

### 2. DeactivateUserDialog

**Props:**
```typescript
interface DeactivateUserDialogProps {
  open: boolean;
  onClose: () => void;
  user: User;
  onSuccess: () => void;
}
```

**UI - Desactivar:**
```
┌────────────────────────────────────┐
│ Desactivar Usuario                 │
├────────────────────────────────────┤
│ ⚠️ ¿Está seguro de desactivar a:  │
│    Juan Pérez (j@example.com)?    │
│                                    │
│ El usuario:                        │
│ • No podrá iniciar sesión          │
│ • Perderá acceso al sistema        │
│ • Puede reactivarse después        │
│                                    │
│        [Cancelar] [Desactivar]     │
└────────────────────────────────────┘
```

**UI - Activar:**
```
┌────────────────────────────────────┐
│ Activar Usuario                    │
├────────────────────────────────────┤
│ ✅ ¿Desea reactivar a:             │
│    Juan Pérez (j@example.com)?    │
│                                    │
│ El usuario recuperará acceso       │
│ completo al sistema.               │
│                                    │
│        [Cancelar] [Activar]        │
└────────────────────────────────────┘
```

**API:** `PUT /users/:id` (actualizar is_active)

---

### 3. InviteUserDialog

**Props:**
```typescript
interface InviteUserDialogProps {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
}
```

**UI:**
```
┌────────────────────────────────────┐
│ Enviar Invitación                  │
├────────────────────────────────────┤
│ [Email*]                           │
│ ejemplo@dominio.com                │
│                                    │
│ [Rol* ▼]                           │
│   • ADMIN                          │
│   • COORDINATOR                    │
│   • BRAND                          │
│                                    │
│ Si Rol = BRAND:                    │
│ [Marca* ▼] Chevrolet               │
│                                    │
│ 💡 Se enviará un email con link    │
│    de registro al usuario          │
│                                    │
│        [Cancelar] [Enviar]         │
└────────────────────────────────────┘
```

**API:** `POST /invitations`

---

### 4. ResendInvitationDialog

**Props:**
```typescript
interface ResendInvitationDialogProps {
  open: boolean;
  onClose: () => void;
  invitation: Invitation;
  onSuccess: () => void;
}
```

**UI:**
```
┌────────────────────────────────────┐
│ Reenviar Invitación                │
├────────────────────────────────────┤
│ 📧 Email: maria@example.com        │
│ Enviada: hace 2 días               │
│                                    │
│ ¿Desea reenviar la invitación?    │
│                                    │
│ Se generará un nuevo token y       │
│ se enviará un nuevo email.         │
│                                    │
│        [Cancelar] [Reenviar]       │
└────────────────────────────────────┘
```

**API:** `POST /invitations/:id/resend`

---

## 🔄 Redux State Management

### usersSlice

```typescript
interface UsersState {
  users: User[];
  selectedUser: User | null;
  loading: boolean;
  error: string | null;
  filters: {
    role?: UserRole;
    isActive?: boolean;
    search?: string;
  };
}

// Actions
- setUsers(users: User[])
- setSelectedUser(user: User | null)
- addUser(user: User)
- updateUser(user: User)
- removeUser(userId: string)
- setLoading(boolean)
- setError(string | null)
- setFilters(filters)
```

### invitationsSlice

```typescript
interface InvitationsState {
  invitations: Invitation[];
  loading: boolean;
  error: string | null;
  filter: 'all' | 'pending' | 'accepted';
}

// Actions
- setInvitations(invitations: Invitation[])
- addInvitation(invitation: Invitation)
- removeInvitation(invitationId: string)
- setLoading(boolean)
- setError(string | null)
- setFilter(filter)
```

---

## 🛣️ Rutas Nuevas

Agregar en `App.tsx`:

```typescript
// Users routes (solo ADMIN)
<Route path="/users" element={
  <ProtectedRoute requireRole="ADMIN">
    <Layout><Users /></Layout>
  </ProtectedRoute>
} />

<Route path="/users/new" element={
  <ProtectedRoute requireRole="ADMIN">
    <Layout><UserForm /></Layout>
  </ProtectedRoute>
} />

<Route path="/users/:id/edit" element={
  <ProtectedRoute requireRole="ADMIN">
    <Layout><UserForm /></Layout>
  </ProtectedRoute>
} />

// Invitations route (solo ADMIN)
<Route path="/invitations" element={
  <ProtectedRoute requireRole="ADMIN">
    <Layout><Invitations /></Layout>
  </ProtectedRoute>
} />
```

---

## 🔐 Control de Acceso

### Permisos por Rol

**ADMIN:**
- ✅ Ver todos los usuarios
- ✅ Crear usuarios
- ✅ Editar cualquier usuario
- ✅ Cambiar roles
- ✅ Desactivar/activar usuarios
- ✅ Gestionar invitaciones

**COORDINATOR:**
- ❌ Sin acceso a gestión de usuarios
- ✅ Solo puede ver su propio perfil (futuro)

**BRAND:**
- ❌ Sin acceso a gestión de usuarios
- ✅ Solo puede ver su propio perfil (futuro)

---

## 🎨 UI/UX Guidelines

### Colores para Roles

```typescript
const roleColors: Record<UserRole, 'primary' | 'secondary' | 'success'> = {
  ADMIN: 'error',      // Rojo/naranja
  COORDINATOR: 'primary',  // Azul
  BRAND: 'success',    // Verde
};

const roleLabels: Record<UserRole, string> = {
  ADMIN: 'Administrador',
  COORDINATOR: 'Coordinador',
  BRAND: 'Marca',
};
```

### Estados Visuales

**Usuario Activo:**
- Chip verde: "Activo"
- Fila normal

**Usuario Inactivo:**
- Chip rojo: "Inactivo"
- Fila con opacity 0.6
- Texto tachado

### Empty States

**Users List vacía:**
```
┌─────────────────────────────────┐
│         👥                      │
│   No hay usuarios registrados   │
│                                 │
│   [+ Crear Primer Usuario]     │
└─────────────────────────────────┘
```

**Invitations vacías:**
```
┌─────────────────────────────────┐
│         📧                      │
│   No hay invitaciones pendientes│
│                                 │
│   [+ Enviar Invitación]        │
└─────────────────────────────────┘
```

---

## 📊 APIs Backend Existentes (Fase 1)

### Users

```
GET    /api/v1/users              - Listar usuarios
POST   /api/v1/users              - Crear usuario (ADMIN)
GET    /api/v1/users/{id}         - Obtener usuario
PUT    /api/v1/users/{id}         - Actualizar usuario (ADMIN)
DELETE /api/v1/users/{id}         - Eliminar usuario (ADMIN)
GET    /api/v1/users/role/{role}  - Usuarios por rol
```

### Invitations

```
GET    /api/v1/invitations        - Listar invitaciones (ADMIN)
POST   /api/v1/invitations        - Crear invitación (ADMIN)
GET    /api/v1/invitations/{id}   - Obtener invitación
POST   /api/v1/invitations/{id}/resend - Reenviar invitación (ADMIN)
POST   /api/v1/invitations/accept - Aceptar invitación (público)
DELETE /api/v1/invitations/{id}   - Eliminar invitación (ADMIN)
```

**Nota:** Todas las APIs ya existen, solo falta el frontend.

---

## 🔧 Orden de Implementación

### ETAPA 1: Redux Slices y API Constants
1. Crear `usersSlice.ts`
2. Crear `invitationsSlice.ts`
3. Agregar a `store/index.ts`
4. Verificar/actualizar `apiConstants.ts`

### ETAPA 2: Página Users List
5. Crear `/pages/Users/index.tsx`
6. Tabla con filtros y búsqueda
7. Integrar con usersSlice
8. Testing básico

### ETAPA 3: Página User Form
9. Crear `/pages/Users/UserForm.tsx`
10. Modo create y edit
11. Validaciones
12. Testing create/edit

### ETAPA 4: Dialogs de Users
13. Crear `ChangeRoleDialog.tsx`
14. Crear `DeactivateUserDialog.tsx`
15. Integrar en Users List
16. Testing dialogs

### ETAPA 5: Página Invitations
17. Crear `/pages/Invitations/index.tsx`
18. Lista con tabs
19. Integrar con invitationsSlice

### ETAPA 6: Dialogs de Invitations
20. Crear `InviteUserDialog.tsx`
21. Crear `ResendInvitationDialog.tsx`
22. Integrar en Invitations page

### ETAPA 7: Navegación y Rutas
23. Agregar rutas en `App.tsx`
24. Agregar items en Layout menu
25. Proteger rutas con ADMIN only

### ETAPA 8: Testing Final
26. Testing end-to-end de todos los flujos
27. Verificar permisos
28. Build y verificación

---

## ✅ Checklist de Completitud

- [ ] usersSlice con actions completas
- [ ] invitationsSlice con actions completas
- [ ] Users List con tabla y filtros
- [ ] User Form create mode
- [ ] User Form edit mode
- [ ] ChangeRoleDialog funcional
- [ ] DeactivateUserDialog funcional
- [ ] Invitations List con tabs
- [ ] InviteUserDialog funcional
- [ ] ResendInvitationDialog funcional
- [ ] Rutas protegidas con ADMIN
- [ ] Items de menú agregados
- [ ] Snackbar feedback en todas las acciones
- [ ] Validaciones frontend completas
- [ ] Error handling robusto
- [ ] Empty states en todas las listas
- [ ] Loading states consistentes
- [ ] Build exitoso

---

## 🚀 Beneficios de Fase 4

### Completitud del Sistema
✅ Admin puede gestionar usuarios sin tocar la BD
✅ Sistema de invitaciones completamente funcional
✅ Gestión de roles y permisos desde UI

### Seguridad
✅ Control de acceso granular
✅ Solo ADMIN puede crear/editar usuarios
✅ Auditoría de cambios (logs en consola)

### UX Administrativa
✅ UI consistente con resto del sistema
✅ Dialogs claros y confirmaciones
✅ Feedback inmediato con snackbars

### Preparación para Producción
✅ Sistema completo sin dependencia de SQL manual
✅ Onboarding de usuarios simplificado
✅ Gestión de usuarios escalable

---

## 📈 Métricas Estimadas

**Tiempo de desarrollo:** 10-15 horas
**Archivos nuevos:** 9
**Archivos modificados:** 2
**Total líneas de código:** ~2,000 LOC
**Componentes reutilizables:** 4 dialogs

---

## 🔜 Después de Fase 4

**Fase 5 (opcional):**
- Analytics Dashboard
- Gráficas y métricas
- Exportación de reportes a PDF
- Sistema de notificaciones
- Audit log completo
- Testing automatizado

---

**Planeado por:** Claude Code (Sonnet 4.5)
**Fecha:** Febrero 2026
**Estado:** 📋 Listo para implementar
