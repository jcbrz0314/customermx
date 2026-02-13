# FASE 3 - Completar Interacciones de Eventos

**Estado:** 🚧 **EN PROGRESO**
**Fecha:** Febrero 2026

---

## 📋 Resumen

La Fase 3 completa las interacciones de usuario para la gestión de eventos que quedaron pendientes en Fase 2. Todas las APIs backend ya existen y funcionan correctamente. Esta fase se enfoca 100% en el frontend, implementando:

1. Dialogs/Modals para gestionar coordinadores y vehículos
2. Sistema de cambio de status de eventos
3. Formulario completo de reportes
4. Sistema de notificaciones (Snackbar/Toast)
5. Estados de carga y feedback visual

---

## 🎯 Alcance de la Implementación

### Frontend Únicamente (12 archivos nuevos)

**No se modifica backend** - Todas las APIs ya existen de Fase 2.

---

## 📦 Componentes a Crear

### 1. Sistema de Notificaciones

**`/frontend/src/components/Snackbar/index.tsx`**
- Context provider para notificaciones globales
- Hook `useSnackbar()` para mostrar mensajes
- Soporte para success, error, warning, info
- Auto-dismiss después de 5 segundos
- Posición: bottom-center

**`/frontend/src/contexts/SnackbarContext.tsx`**
- Context con estado global de notificaciones
- Queue para múltiples notificaciones

### 2. Dialog para Asignar Coordinadores

**`/frontend/src/components/Dialogs/AssignCoordinatorDialog.tsx`**

**Props:**
```typescript
interface AssignCoordinatorDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  onSuccess: () => void;
}
```

**Funcionalidad:**
- Select con lista de usuarios con rol COORDINATOR
- Filtrar usuarios que ya están asignados
- API call: `POST /api/v1/events/{eventId}/coordinators`
- Loading state en botón
- Cerrar al éxito y mostrar snackbar

**UI:**
- Title: "Asignar Coordinador"
- FormControl con Select (fullWidth)
- Botones: Cancelar (outlined), Asignar (contained)

### 3. Dialog para Agregar Vehículos

**`/frontend/src/components/Dialogs/AddVehicleDialog.tsx`**

**Props:**
```typescript
interface AddVehicleDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  eventBrandId: string;
  onSuccess: () => void;
}
```

**Funcionalidad:**
- Select con vehículos filtrados por `eventBrandId`
- TextField para cantidad (number, min: 1)
- API call: `POST /api/v1/events/{eventId}/vehicles`
- Validación: vehículo seleccionado + cantidad > 0
- Error handling para ErrBrandMismatch (400)

**UI:**
- Title: "Agregar Vehículo"
- Select de vehículos (mostrar modelo)
- TextField cantidad con helper text
- Botones: Cancelar, Agregar

### 4. Dialog para Cambiar Status

**`/frontend/src/components/Dialogs/ChangeStatusDialog.tsx`**

**Props:**
```typescript
interface ChangeStatusDialogProps {
  open: boolean;
  onClose: () => void;
  event: Event;
  onSuccess: () => void;
}
```

**Funcionalidad:**
- Select con status válidos según transición actual
- Validar transiciones permitidas en frontend
- API call: `PATCH /api/v1/events/{id}/status`
- Confirmación: "¿Está seguro de cambiar el status?"
- Mostrar current status y new status

**Transiciones válidas:**
```typescript
const validTransitions: Record<EventStatus, EventStatus[]> = {
  PLANNED: ['ACTIVE'],
  ACTIVE: ['COMPLETED'],
  COMPLETED: ['CLOSED'],
  CLOSED: [],
};
```

**UI:**
- Title: "Cambiar Status del Evento"
- Alert info con status actual
- Select con opciones válidas únicamente
- Typography con explicación de la transición
- Botones: Cancelar, Confirmar Cambio

### 5. Dialog para Editar Cantidad de Vehículo

**`/frontend/src/components/Dialogs/EditVehicleQuantityDialog.tsx`**

**Props:**
```typescript
interface EditVehicleQuantityDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  vehicleId: string;
  currentQuantity: number;
  vehicleName: string;
  onSuccess: () => void;
}
```

**Funcionalidad:**
- TextField pre-llenado con cantidad actual
- API call: `PATCH /api/v1/events/{eventId}/vehicles/{vehicleId}/quantity`
- Validación: cantidad > 0

**UI:**
- Title: "Editar Cantidad"
- Subtitle: nombre del vehículo
- TextField tipo number
- Botones: Cancelar, Guardar

### 6. Dialog para Gestionar Reporte

**`/frontend/src/components/Dialogs/EventReportDialog.tsx`**

**Props:**
```typescript
interface EventReportDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  existingReport?: EventReport;
  onSuccess: () => void;
}
```

**Funcionalidad:**
- Formulario completo con todos los campos del reporte
- Pre-llenar si `existingReport` existe (modo edit)
- API call: `POST /api/v1/events/{eventId}/report` (UPSERT)
- Campos opcionales (usar pointers en request)
- Validación: dealer_rating 1-5

**Campos del formulario:**
- hostess_count (number)
- setup_vendor (text)
- has_promotional (checkbox)
- attendees (number)
- activities_count (number)
- leads_collected (number)
- prospects (number)
- dealer_rating (rating 1-5 con estrellas)
- comments (multiline text)
- completed (checkbox)

**UI:**
- Title: "Reporte del Evento"
- Grid 2 columnas para campos
- Rating component para dealer_rating
- TextField multiline para comments
- Checkbox destacado para "Marcar como completado"
- Botones: Cancelar, Guardar Reporte

### 7. Dialog de Confirmación Genérico

**`/frontend/src/components/Dialogs/ConfirmDialog.tsx`**

**Props:**
```typescript
interface ConfirmDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  severity?: 'warning' | 'error' | 'info';
}
```

**Funcionalidad:**
- Dialog reutilizable para confirmaciones
- Usar para remover coordinadores, eliminar vehículos, etc.

---

## 🔄 Modificaciones a Páginas Existentes

### EventDetail.tsx

**Agregar estados locales:**
```typescript
const [assignCoordinatorOpen, setAssignCoordinatorOpen] = useState(false);
const [addVehicleOpen, setAddVehicleOpen] = useState(false);
const [editQuantityOpen, setEditQuantityOpen] = useState(false);
const [selectedVehicle, setSelectedVehicle] = useState<EventVehicleWithDetails | null>(null);
const [changeStatusOpen, setChangeStatusOpen] = useState(false);
const [reportDialogOpen, setReportDialogOpen] = useState(false);
const [confirmDialog, setConfirmDialog] = useState<ConfirmDialogState | null>(null);
```

**Funciones handlers:**
```typescript
const handleAssignCoordinator = () => setAssignCoordinatorOpen(true);
const handleRemoveCoordinator = (userId: string, userName: string) => {
  setConfirmDialog({
    open: true,
    title: 'Remover Coordinador',
    message: `¿Está seguro de remover a ${userName} de este evento?`,
    onConfirm: async () => {
      // API call DELETE
    },
  });
};
const handleAddVehicle = () => setAddVehicleOpen(true);
const handleEditVehicleQuantity = (vehicle: EventVehicleWithDetails) => {
  setSelectedVehicle(vehicle);
  setEditQuantityOpen(true);
};
const handleRemoveVehicle = (vehicleId: string, vehicleName: string) => {
  // Similar a remove coordinator
};
const handleChangeStatus = () => setChangeStatusOpen(true);
const handleEditReport = () => setReportDialogOpen(true);
```

**Agregar botones/iconos:**
- En lista de coordinadores: IconButton delete (solo ADMIN)
- En lista de vehículos: IconButtons edit y delete (solo ADMIN)
- Botón "Cambiar Status" en header (solo ADMIN)
- Implementar onClick en botones existentes

**Imports de dialogs:**
```typescript
import { AssignCoordinatorDialog } from '../../components/Dialogs/AssignCoordinatorDialog';
import { AddVehicleDialog } from '../../components/Dialogs/AddVehicleDialog';
import { EditVehicleQuantityDialog } from '../../components/Dialogs/EditVehicleQuantityDialog';
import { ChangeStatusDialog } from '../../components/Dialogs/ChangeStatusDialog';
import { EventReportDialog } from '../../components/Dialogs/EventReportDialog';
import { ConfirmDialog } from '../../components/Dialogs/ConfirmDialog';
import { useSnackbar } from '../../hooks/useSnackbar';
```

**Renderizar dialogs al final del componente:**
```typescript
<AssignCoordinatorDialog
  open={assignCoordinatorOpen}
  onClose={() => setAssignCoordinatorOpen(false)}
  eventId={id!}
  onSuccess={() => {
    fetchEventData();
    showSnackbar('Coordinador asignado exitosamente', 'success');
  }}
/>
// ... resto de dialogs
```

---

## 📝 Estructura de Archivos a Crear

```
frontend/src/
├── components/
│   ├── Dialogs/
│   │   ├── AssignCoordinatorDialog.tsx
│   │   ├── AddVehicleDialog.tsx
│   │   ├── EditVehicleQuantityDialog.tsx
│   │   ├── ChangeStatusDialog.tsx
│   │   ├── EventReportDialog.tsx
│   │   └── ConfirmDialog.tsx
│   └── Snackbar/
│       └── index.tsx
├── contexts/
│   └── SnackbarContext.tsx
└── hooks/
    └── useSnackbar.ts
```

**Modificar:**
- `/frontend/src/pages/Events/EventDetail.tsx`
- `/frontend/src/App.tsx` (wrap con SnackbarProvider)

**Total: 12 archivos** (9 nuevos, 2 modificados, 1 hook)

---

## 🎨 Decisiones de Diseño UI/UX

### Dialogs
- Width: `maxWidth="sm"` (600px)
- Bordes redondeados (del theme: 16px)
- Botones: Cancelar (outlined), Acción (contained)
- Loading spinner en botón de acción cuando API en progreso

### Snackbar
- Posición: `anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}`
- Auto-hide: 5000ms
- Colores según severity (success=green, error=red, etc.)
- Action: botón "X" para cerrar manualmente

### Confirmaciones
- Dialog pequeño con Alert severity="warning"
- Texto claro del impacto de la acción
- Botón confirmar en color error si es destructivo

### Feedback Visual
- Mostrar CircularProgress en botones durante API calls
- Deshabilitar botón mientras loading
- Actualizar datos después de éxito (llamar fetchEventData)
- Snackbar con mensaje descriptivo

---

## 🔗 APIs Backend Utilizadas (Ya Existen)

### Coordinadores
```
POST   /api/v1/events/{eventId}/coordinators
DELETE /api/v1/events/{eventId}/coordinators/{userId}
```

### Vehículos
```
POST   /api/v1/events/{eventId}/vehicles
PATCH  /api/v1/events/{eventId}/vehicles/{vehicleId}/quantity
DELETE /api/v1/events/{eventId}/vehicles/{vehicleId}
```

### Status
```
PATCH  /api/v1/events/{id}/status
```

### Reportes
```
POST   /api/v1/events/{eventId}/report
```

### Listados (para poblar Selects)
```
GET /api/v1/users (filtrar por role=COORDINATOR)
GET /api/v1/vehicles (filtrar por brandId)
```

---

## 🧪 Testing Manual

### Asignar Coordinador
1. Abrir detalle de evento como ADMIN
2. Click "Asignar" en card Coordinadores
3. Seleccionar coordinador del dropdown
4. Click "Asignar"
5. Verificar:
   - ✓ Loading spinner en botón
   - ✓ Dialog se cierra
   - ✓ Snackbar "Coordinador asignado exitosamente"
   - ✓ Lista se actualiza con nuevo coordinador

### Remover Coordinador
1. Click icono delete junto a coordinador
2. Dialog de confirmación aparece
3. Click "Confirmar"
4. Verificar:
   - ✓ Coordinador removido de lista
   - ✓ Snackbar de éxito

### Agregar Vehículo
1. Click "Agregar" en card Vehículos
2. Seleccionar vehículo (solo de la marca del evento)
3. Ingresar cantidad
4. Click "Agregar"
5. Verificar:
   - ✓ Vehículo aparece en lista con cantidad
   - ✓ Snackbar de éxito

### Editar Cantidad Vehículo
1. Click icono edit junto a vehículo
2. Cambiar cantidad
3. Guardar
4. Verificar cantidad actualizada

### Remover Vehículo
1. Click icono delete junto a vehículo
2. Confirmar
3. Verificar vehículo removido

### Cambiar Status
1. Como ADMIN, click botón "Cambiar Status" (agregar en header)
2. Select muestra solo transiciones válidas
3. Seleccionar nuevo status
4. Confirmar
5. Verificar:
   - ✓ Chip de status actualizado
   - ✓ Snackbar de éxito
   - ✓ No se puede retroceder (COMPLETED no puede volver a ACTIVE)

### Crear/Editar Reporte
1. Click "Editar Reporte" en card Reporte
2. Llenar campos del formulario
3. Guardar
4. Verificar:
   - ✓ Datos del reporte se muestran en card
   - ✓ UPSERT funciona (actualiza campos existentes)
   - ✓ Checkbox "Completado" funciona

---

## 🎯 Orden de Implementación

### ETAPA 1: Sistema de Notificaciones (Base)
1. Crear `SnackbarContext.tsx`
2. Crear `useSnackbar.ts` hook
3. Crear componente `Snackbar/index.tsx`
4. Modificar `App.tsx` para wrap con SnackbarProvider

### ETAPA 2: Dialog de Confirmación Genérico
5. Crear `ConfirmDialog.tsx`

### ETAPA 3: Dialogs de Coordinadores
6. Crear `AssignCoordinatorDialog.tsx`
7. Modificar `EventDetail.tsx` para agregar botón y handler
8. Testing asignar coordinador
9. Implementar remover coordinador con ConfirmDialog

### ETAPA 4: Dialogs de Vehículos
10. Crear `AddVehicleDialog.tsx`
11. Crear `EditVehicleQuantityDialog.tsx`
12. Modificar `EventDetail.tsx` para agregar botones y handlers
13. Testing agregar, editar, remover vehículos

### ETAPA 5: Cambio de Status
14. Crear `ChangeStatusDialog.tsx`
15. Agregar botón "Cambiar Status" en EventDetail header
16. Testing transiciones de status

### ETAPA 6: Gestión de Reportes
17. Crear `EventReportDialog.tsx`
18. Conectar botón "Editar Reporte" existente
19. Testing crear y editar reportes

---

## 📊 Tipos TypeScript Nuevos

**`/frontend/src/types/index.ts` - Agregar:**

```typescript
export interface AssignCoordinatorRequest {
  event_id: string;
  user_id: string;
}

export interface AddVehicleRequest {
  event_id: string;
  vehicle_id: string;
  quantity: number;
}

export interface UpdateVehicleQuantityRequest {
  quantity: number;
}

export interface ChangeStatusRequest {
  status: EventStatus;
}

export interface CreateOrUpdateReportRequest {
  hostess_count?: number;
  setup_vendor?: string;
  has_promotional?: boolean;
  attendees?: number;
  activities_count?: number;
  leads_collected?: number;
  prospects?: number;
  dealer_rating?: number;
  comments?: string;
  completed?: boolean;
}
```

---

## 🚀 Beneficios de Esta Implementación

### Usabilidad
✅ Usuarios pueden completar flujos sin salir de la página de detalle
✅ Feedback inmediato con snackbars
✅ Confirmaciones previenen errores

### UX
✅ Dialogs modales mantienen contexto
✅ Loading states claros
✅ Validaciones en frontend antes de API calls

### Mantenibilidad
✅ Componentes reutilizables (ConfirmDialog)
✅ Context pattern para notificaciones
✅ Separación de concerns

### Completitud
✅ Todas las APIs backend quedan conectadas al frontend
✅ Sistema funcional end-to-end
✅ Listo para producción

---

## 📝 Checklist de Completitud

- [ ] SnackbarContext y hook
- [ ] Componente Snackbar
- [ ] ConfirmDialog
- [ ] AssignCoordinatorDialog
- [ ] Remover coordinador con confirmación
- [ ] AddVehicleDialog
- [ ] EditVehicleQuantityDialog
- [ ] Remover vehículo con confirmación
- [ ] ChangeStatusDialog
- [ ] EventReportDialog
- [ ] Actualizar EventDetail con todos los handlers
- [ ] Testing completo de todos los flows

---

## 🔜 Próximos Pasos (Fase 4 - Opcional)

Fuera del alcance de Fase 3:
- Analytics y dashboards
- Exportación de reportes a PDF
- Sistema de búsqueda avanzada
- Filtros guardados
- Historial de cambios (audit log)
- Notificaciones push

---

**Estado:** 🚧 Listo para comenzar implementación
**Tiempo estimado:** 8-12 horas
**Archivos:** 12 (9 nuevos, 3 modificados)
