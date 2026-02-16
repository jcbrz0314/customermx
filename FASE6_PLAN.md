# Plan de Implementación: Fase 6 - Completar Analytics

## Contexto

La Fase 5 implementó el sistema base de analytics con visualizaciones interactivas. Los usuarios pueden ver métricas clave, gráficas y filtrar datos. Sin embargo, faltan funcionalidades importantes para hacer el sistema completamente útil:

1. **Exportación a PDF** - Las marcas necesitan compartir reportes físicos/digitales
2. **Drill-down Interactivo** - Click en gráficas para filtrar y explorar detalles
3. **Métricas Adicionales** - Más insights de negocio (eventos por tipo, distribuidores, etc.)

Esta fase completa el sistema de analytics agregando estas capacidades.

## Alcance de la Implementación

**Exportación a PDF:**
- Generar PDF de evento individual con todos los detalles
- Generar PDF de dashboard analytics con gráficas
- Templates profesionales con logo y branding
- Descarga directa desde UI
- Biblioteca: jsPDF + html2canvas

**Drill-down Interactivo:**
- Click en barra de marca → navegar a eventos de esa marca
- Click en mes de timeline → filtrar eventos de ese mes
- Click en vehículo → ver eventos donde se presentó
- Click en estado → filtrar eventos de ese estado
- Breadcrumbs para tracking de navegación

**Métricas Adicionales:**
- Eventos por tipo (Triatlón, Maratón, Golf, etc.)
- Ranking de distribuidores por rating
- Conversión leads → prospectos
- Eventos por ciudad (top 10)
- Distribución de edecanes promedio

**Fuera de alcance (Fase 7):**
- Materialized views
- Caching con Redis
- Real-time updates con WebSocket
- Dashboards personalizables

## Arquitectura de Exportación PDF

### Frontend: jsPDF + html2canvas

**jsPDF** - Generar documentos PDF
**html2canvas** - Capturar gráficas como imágenes

Flujo:
1. Usuario click "Exportar PDF"
2. Capturar gráficas con html2canvas
3. Generar PDF con jsPDF
4. Agregar texto, tablas, imágenes
5. Descargar automáticamente

### Templates

**Template 1: Reporte de Evento Individual**
- Header con logo
- Información del evento
- Métricas del reporte
- Tabla de vehículos presentados
- Tabla de coordinadores
- Comentarios

**Template 2: Dashboard Analytics**
- Header con filtros aplicados
- 4 KPIs en grid
- Gráficas capturadas como imágenes
- Timestamp de generación

## Arquitectura de Drill-down

### Frontend: Estado de Navegación

Usar query params + estado local para tracking:

```typescript
// Al hacer click en marca
const handleBrandClick = (brandId: string, brandName: string) => {
  navigate(`/events?brand=${brandId}&name=${brandName}&from=analytics`);
};

// En página Events, detectar origen
const params = new URLSearchParams(location.search);
if (params.get('from') === 'analytics') {
  // Mostrar breadcrumb: Analytics > Eventos de {name}
}
```

### Breadcrumbs Component

Componente reutilizable para mostrar ruta de navegación:
```
Home > Analytics > Eventos de Chevrolet
```

## Orden de Implementación

### ETAPA 1: Instalación de Dependencias

```bash
npm install jspdf html2canvas
npm install --save-dev @types/html2canvas
```

---

### ETAPA 2: Backend - Endpoints de Métricas Adicionales

Agregar nuevos métodos al dominio analytics.

#### Archivos a modificar:

**`/backend/internal/domain/analytics/model.go`**

Agregar nuevos structs:

```go
// Eventos por tipo
type EventTypeMetrics struct {
    EventType  string `json:"event_type"`
    EventCount int    `json:"event_count"`
    Attendees  int    `json:"attendees"`
}

// Ranking de distribuidores
type DealerMetrics struct {
    Dealer        string  `json:"dealer"`
    EventCount    int     `json:"event_count"`
    AverageRating float64 `json:"average_rating"`
    TotalAttendees int    `json:"total_attendees"`
}

// Conversión
type ConversionMetrics struct {
    TotalAttendees int     `json:"total_attendees"`
    TotalLeads     int     `json:"total_leads"`
    TotalProspects int     `json:"total_prospects"`
    LeadRate       float64 `json:"lead_rate"`      // leads / attendees
    ProspectRate   float64 `json:"prospect_rate"`  // prospects / leads
}

// Eventos por ciudad
type CityMetrics struct {
    State      string `json:"state"`
    City       string `json:"city"`
    EventCount int    `json:"event_count"`
    Attendees  int    `json:"attendees"`
}
```

**`/backend/internal/domain/analytics/repository.go`**

Agregar queries:

```go
GetMetricsByEventType(ctx context.Context, filters AnalyticsFilters) ([]EventTypeMetrics, error)
GetMetricsByDealer(ctx context.Context, filters AnalyticsFilters) ([]DealerMetrics, error)
GetConversionMetrics(ctx context.Context, filters AnalyticsFilters) (*ConversionMetrics, error)
GetMetricsByCity(ctx context.Context, filters AnalyticsFilters, limit int) ([]CityMetrics, error)
```

**Queries SQL:**

```sql
-- Por tipo de evento
SELECT
    e.event_type,
    COUNT(e.id) as event_count,
    COALESCE(SUM(er.attendees), 0) as attendees
FROM events e
LEFT JOIN event_reports er ON e.id = er.event_id
WHERE ($1::uuid IS NULL OR e.brand_id = $1)
  AND ($2::int IS NULL OR e.year = $2)
  AND e.status = 'COMPLETED'
GROUP BY e.event_type
ORDER BY event_count DESC;

-- Ranking de distribuidores
SELECT
    e.dealer,
    COUNT(e.id) as event_count,
    COALESCE(AVG(er.dealer_rating), 0) as average_rating,
    COALESCE(SUM(er.attendees), 0) as total_attendees
FROM events e
LEFT JOIN event_reports er ON e.id = er.event_id
WHERE ($1::uuid IS NULL OR e.brand_id = $1)
  AND ($2::int IS NULL OR e.year = $2)
  AND e.status = 'COMPLETED'
  AND e.dealer IS NOT NULL AND e.dealer != ''
GROUP BY e.dealer
ORDER BY average_rating DESC, event_count DESC
LIMIT 10;

-- Conversión
SELECT
    COALESCE(SUM(er.attendees), 0) as total_attendees,
    COALESCE(SUM(er.leads_collected), 0) as total_leads,
    COALESCE(SUM(er.prospects), 0) as total_prospects,
    CASE
        WHEN SUM(er.attendees) > 0 THEN
            (SUM(er.leads_collected)::float / SUM(er.attendees)::float) * 100
        ELSE 0
    END as lead_rate,
    CASE
        WHEN SUM(er.leads_collected) > 0 THEN
            (SUM(er.prospects)::float / SUM(er.leads_collected)::float) * 100
        ELSE 0
    END as prospect_rate
FROM events e
JOIN event_reports er ON e.id = er.event_id
WHERE ($1::uuid IS NULL OR e.brand_id = $1)
  AND ($2::int IS NULL OR e.year = $2)
  AND e.status = 'COMPLETED';

-- Por ciudad (top N)
SELECT
    e.state,
    e.city,
    COUNT(e.id) as event_count,
    COALESCE(SUM(er.attendees), 0) as attendees
FROM events e
LEFT JOIN event_reports er ON e.id = er.event_id
WHERE ($1::uuid IS NULL OR e.brand_id = $1)
  AND ($2::int IS NULL OR e.year = $2)
  AND e.status = 'COMPLETED'
GROUP BY e.state, e.city
ORDER BY event_count DESC
LIMIT $3;
```

**`/backend/internal/domain/analytics/service.go`**

Actualizar DashboardAnalytics:

```go
type DashboardAnalytics struct {
    Totals         TotalMetrics        `json:"totals"`
    ByBrand        []BrandMetrics      `json:"by_brand"`
    ByMonth        []MonthlyMetrics    `json:"by_month"`
    ByState        []StateMetrics      `json:"by_state"`
    TopVehicles    []VehicleMetrics    `json:"top_vehicles"`
    YearComparison []YearComparison    `json:"year_comparison"`
    // NUEVAS
    ByEventType    []EventTypeMetrics  `json:"by_event_type"`
    TopDealers     []DealerMetrics     `json:"top_dealers"`
    Conversion     *ConversionMetrics  `json:"conversion"`
    TopCities      []CityMetrics       `json:"top_cities"`
}
```

Actualizar GetDashboardAnalytics para incluir nuevas métricas.

**`/backend/internal/http/handlers/analytics.go`**

Handlers ya existentes retornarán las nuevas métricas automáticamente.

---

### ETAPA 3: Frontend - Tipos TypeScript

**`/frontend/src/types/index.ts`**

Agregar:

```typescript
export interface EventTypeMetrics {
  event_type: string;
  event_count: number;
  attendees: number;
}

export interface DealerMetrics {
  dealer: string;
  event_count: number;
  average_rating: number;
  total_attendees: number;
}

export interface ConversionMetrics {
  total_attendees: number;
  total_leads: number;
  total_prospects: number;
  lead_rate: number;    // percentage
  prospect_rate: number; // percentage
}

export interface CityMetrics {
  state: string;
  city: string;
  event_count: number;
  attendees: number;
}

// Actualizar DashboardAnalytics
export interface DashboardAnalytics {
  totals: TotalMetrics;
  by_brand: BrandMetrics[];
  by_month: MonthlyMetrics[];
  by_state: StateMetrics[];
  top_vehicles: VehicleMetrics[];
  year_comparison: YearComparison[];
  by_event_type: EventTypeMetrics[];
  top_dealers: DealerMetrics[];
  conversion: ConversionMetrics;
  top_cities: CityMetrics[];
}
```

---

### ETAPA 4: Frontend - Nuevos Componentes de Gráficas

**`/frontend/src/components/Charts/EventTypeChart.tsx`**

```tsx
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const COLORS = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6'];

export const EventTypeChart: React.FC<{ data: EventTypeMetrics[] }> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={data}
          dataKey="event_count"
          nameKey="event_type"
          cx="50%"
          cy="50%"
          outerRadius={80}
          label
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
};
```

**`/frontend/src/components/Charts/DealerRankingChart.tsx`**

```tsx
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

export const DealerRankingChart: React.FC<{ data: DealerMetrics[] }> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data} layout="vertical">
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis type="number" />
        <YAxis dataKey="dealer" type="category" width={150} />
        <Tooltip />
        <Bar dataKey="average_rating" fill="#f59e0b" name="Rating" />
      </BarChart>
    </ResponsiveContainer>
  );
};
```

**`/frontend/src/components/ConversionFunnel/index.tsx`**

```tsx
import { Card, CardContent, Typography, Box } from '@mui/material';

export const ConversionFunnel: React.FC<{ data: ConversionMetrics }> = ({ data }) => {
  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Funnel de Conversión
        </Typography>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
          {/* Level 1: Attendees */}
          <Box sx={{ bgcolor: '#2563eb', p: 2, borderRadius: 2, color: 'white' }}>
            <Typography variant="h6">{data.total_attendees.toLocaleString()}</Typography>
            <Typography variant="caption">Asistentes</Typography>
          </Box>

          {/* Level 2: Leads */}
          <Box sx={{ bgcolor: '#10b981', p: 2, borderRadius: 2, color: 'white', width: '70%' }}>
            <Typography variant="h6">{data.total_leads.toLocaleString()}</Typography>
            <Typography variant="caption">
              Leads ({data.lead_rate.toFixed(1)}%)
            </Typography>
          </Box>

          {/* Level 3: Prospects */}
          <Box sx={{ bgcolor: '#f59e0b', p: 2, borderRadius: 2, color: 'white', width: '40%' }}>
            <Typography variant="h6">{data.total_prospects.toLocaleString()}</Typography>
            <Typography variant="caption">
              Prospectos ({data.prospect_rate.toFixed(1)}%)
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};
```

---

### ETAPA 5: Frontend - Componentes de Exportación PDF

**`/frontend/src/utils/pdfExport.ts`**

```typescript
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';
import { Event, EventReport, EventVehicleWithDetails, EventCoordinatorWithUser } from '../types';

export const exportEventToPDF = async (
  event: Event,
  report: EventReport | null,
  vehicles: EventVehicleWithDetails[],
  coordinators: EventCoordinatorWithUser[]
) => {
  const pdf = new jsPDF();
  let yPos = 20;

  // Header
  pdf.setFontSize(20);
  pdf.text('Reporte de Evento', 20, yPos);
  yPos += 15;

  // Event info
  pdf.setFontSize(12);
  pdf.text(`Nombre: ${event.name}`, 20, yPos);
  yPos += 7;
  pdf.text(`Tipo: ${event.event_type}`, 20, yPos);
  yPos += 7;
  pdf.text(`Fecha: ${new Date(event.start_date).toLocaleDateString()}`, 20, yPos);
  yPos += 7;
  pdf.text(`Ubicación: ${event.city}, ${event.state}`, 20, yPos);
  yPos += 7;
  pdf.text(`Sede: ${event.venue}`, 20, yPos);
  yPos += 10;

  // Report metrics
  if (report) {
    pdf.setFontSize(14);
    pdf.text('Métricas del Evento', 20, yPos);
    yPos += 10;

    pdf.setFontSize(11);
    pdf.text(`Asistentes: ${report.attendees || 'N/A'}`, 20, yPos);
    yPos += 6;
    pdf.text(`Leads: ${report.leads_collected || 'N/A'}`, 20, yPos);
    yPos += 6;
    pdf.text(`Prospectos: ${report.prospects || 'N/A'}`, 20, yPos);
    yPos += 6;
    pdf.text(`Rating: ${report.dealer_rating || 'N/A'} / 5`, 20, yPos);
    yPos += 10;
  }

  // Vehicles table
  if (vehicles.length > 0) {
    pdf.setFontSize(14);
    pdf.text('Vehículos Presentados', 20, yPos);
    yPos += 10;

    pdf.setFontSize(10);
    vehicles.forEach((v) => {
      pdf.text(`• ${v.model_name} (${v.brand_name}): ${v.quantity}`, 25, yPos);
      yPos += 6;
    });
    yPos += 5;
  }

  // Coordinators
  if (coordinators.length > 0) {
    pdf.setFontSize(14);
    pdf.text('Coordinadores', 20, yPos);
    yPos += 10;

    pdf.setFontSize(10);
    coordinators.forEach((c) => {
      pdf.text(`• ${c.user_name} (${c.user_email})`, 25, yPos);
      yPos += 6;
    });
  }

  // Download
  pdf.save(`evento_${event.id}.pdf`);
};

export const exportDashboardToPDF = async (elementId: string, fileName: string) => {
  const element = document.getElementById(elementId);
  if (!element) return;

  // Capture as image
  const canvas = await html2canvas(element, {
    scale: 2,
    logging: false,
  });

  const imgData = canvas.toDataURL('image/png');
  const pdf = new jsPDF('p', 'mm', 'a4');

  const imgWidth = 210; // A4 width in mm
  const imgHeight = (canvas.height * imgWidth) / canvas.width;

  pdf.addImage(imgData, 'PNG', 0, 0, imgWidth, imgHeight);
  pdf.save(fileName);
};
```

**`/frontend/src/components/ExportButton/index.tsx`**

```tsx
import { Button } from '@mui/material';
import { PictureAsPdf } from '@mui/icons-material';

interface Props {
  onClick: () => void;
  loading?: boolean;
}

export const ExportButton: React.FC<Props> = ({ onClick, loading }) => {
  return (
    <Button
      variant="outlined"
      startIcon={<PictureAsPdf />}
      onClick={onClick}
      disabled={loading}
    >
      Exportar PDF
    </Button>
  );
};
```

---

### ETAPA 6: Frontend - Drill-down en Gráficas

**`/frontend/src/pages/Analytics/index.tsx`**

Agregar handlers de click:

```tsx
const handleBrandClick = (data: BrandMetrics) => {
  navigate(`/events?brand_id=${data.brand_id}&brand_name=${data.brand_name}&from=analytics`);
};

const handleMonthClick = (data: MonthlyMetrics) => {
  navigate(`/events?year=${data.year}&month=${data.month}&from=analytics`);
};

const handleVehicleClick = (data: VehicleMetrics) => {
  navigate(`/events?vehicle_id=${data.vehicle_id}&vehicle_name=${data.model_name}&from=analytics`);
};

const handleStateClick = (data: StateMetrics) => {
  navigate(`/events?state=${data.state}&from=analytics`);
};

// En los componentes de charts
<EventsByBrandChart data={analytics.by_brand} onClick={handleBrandClick} />
```

**Actualizar componentes de charts para aceptar onClick:**

```tsx
// EventsByBrandChart.tsx
interface Props {
  data: BrandMetrics[];
  onClick?: (data: BrandMetrics) => void;
}

export const EventsByBrandChart: React.FC<Props> = ({ data, onClick }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data}>
        {/* ... */}
        <Bar
          dataKey="event_count"
          fill="#2563eb"
          name="Eventos"
          onClick={onClick}
          cursor={onClick ? 'pointer' : 'default'}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};
```

---

### ETAPA 7: Frontend - Breadcrumbs Component

**`/frontend/src/components/Breadcrumbs/index.tsx`**

```tsx
import { Breadcrumbs as MUIBreadcrumbs, Link, Typography } from '@mui/material';
import { NavigateNext } from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

export const Breadcrumbs: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const params = new URLSearchParams(location.search);

  const from = params.get('from');
  if (!from) return null;

  const breadcrumbMap: Record<string, string> = {
    analytics: 'Analytics',
  };

  return (
    <MUIBreadcrumbs separator={<NavigateNext fontSize="small" />} sx={{ mb: 2 }}>
      <Link
        component="button"
        variant="body2"
        onClick={() => navigate('/dashboard')}
        sx={{ cursor: 'pointer' }}
      >
        Inicio
      </Link>
      <Link
        component="button"
        variant="body2"
        onClick={() => navigate(`/${from}`)}
        sx={{ cursor: 'pointer' }}
      >
        {breadcrumbMap[from] || from}
      </Link>
      <Typography variant="body2" color="text.primary">
        {params.get('brand_name') || params.get('vehicle_name') || 'Eventos'}
      </Typography>
    </MUIBreadcrumbs>
  );
};
```

**Agregar en `/frontend/src/pages/Events/index.tsx`:**

```tsx
import { Breadcrumbs } from '../../components/Breadcrumbs';

// En el componente
<Box sx={{ p: 3 }}>
  <Breadcrumbs />
  {/* resto del contenido */}
</Box>
```

---

### ETAPA 8: Frontend - Actualizar Página Analytics

Agregar:
1. Botón "Exportar PDF"
2. Nuevas gráficas (tipo evento, dealers, conversión, ciudades)
3. Handlers de click en todas las gráficas
4. ID en contenedor principal para captura

**`/frontend/src/pages/Analytics/index.tsx`**

```tsx
import { exportDashboardToPDF } from '../../utils/pdfExport';
import { ExportButton } from '../../components/ExportButton';
import { EventTypeChart } from '../../components/Charts/EventTypeChart';
import { DealerRankingChart } from '../../components/Charts/DealerRankingChart';
import { ConversionFunnel } from '../../components/ConversionFunnel';

// En el header
<Box sx={{ display: 'flex', gap: 2 }}>
  <ExportButton
    onClick={() => exportDashboardToPDF('analytics-dashboard', 'analytics.pdf')}
  />
</Box>

// Contenedor con ID
<Box id="analytics-dashboard">
  {/* KPIs y gráficas */}
</Box>

// Nuevas gráficas
{analytics.by_event_type.length > 0 && (
  <Grid item xs={12} lg={6}>
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Eventos por Tipo
        </Typography>
        <EventTypeChart data={analytics.by_event_type} />
      </CardContent>
    </Card>
  </Grid>
)}

{analytics.top_dealers.length > 0 && (
  <Grid item xs={12} lg={6}>
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Top 10 Distribuidores
        </Typography>
        <DealerRankingChart data={analytics.top_dealers} />
      </CardContent>
    </Card>
  </Grid>
)}

{analytics.conversion && (
  <Grid item xs={12} lg={6}>
    <ConversionFunnel data={analytics.conversion} />
  </Grid>
)}
```

---

### ETAPA 9: Frontend - Exportar Evento Individual

**`/frontend/src/pages/Events/EventDetail.tsx`**

Agregar botón de exportación:

```tsx
import { exportEventToPDF } from '../../utils/pdfExport';
import { ExportButton } from '../../components/ExportButton';

// En el header
<Box sx={{ display: 'flex', gap: 2 }}>
  <ExportButton
    onClick={() => exportEventToPDF(event, report, vehicles, coordinators)}
  />
  {/* otros botones */}
</Box>
```

---

## Archivos a Crear/Modificar

### Backend - Modificar (3 archivos):
1. `/backend/internal/domain/analytics/model.go` - Nuevos structs
2. `/backend/internal/domain/analytics/repository.go` - 4 queries nuevas
3. `/backend/internal/domain/analytics/service.go` - Actualizar DashboardAnalytics

### Frontend - Nuevos (7 archivos):
4. `/frontend/src/components/Charts/EventTypeChart.tsx`
5. `/frontend/src/components/Charts/DealerRankingChart.tsx`
6. `/frontend/src/components/ConversionFunnel/index.tsx`
7. `/frontend/src/components/ExportButton/index.tsx`
8. `/frontend/src/components/Breadcrumbs/index.tsx`
9. `/frontend/src/utils/pdfExport.ts`

### Frontend - Modificar (7 archivos):
10. `/frontend/src/types/index.ts` - Nuevos tipos
11. `/frontend/src/pages/Analytics/index.tsx` - Drill-down + nuevas gráficas + export
12. `/frontend/src/pages/Events/index.tsx` - Breadcrumbs
13. `/frontend/src/pages/Events/EventDetail.tsx` - Export button
14. `/frontend/src/components/Charts/EventsByBrandChart.tsx` - onClick prop
15. `/frontend/src/components/Charts/EventTimelineChart.tsx` - onClick prop
16. `/frontend/src/components/Charts/TopVehiclesChart.tsx` - onClick prop
17. `/frontend/src/components/Charts/StateDistributionChart.tsx` - onClick prop

**Total:** 17 archivos

---

## Tiempos Estimados

- Backend métricas adicionales: 3-4 horas
- Frontend componentes gráficas: 2-3 horas
- Sistema de exportación PDF: 4-5 horas
- Drill-down + breadcrumbs: 2-3 horas
- Testing y ajustes: 2-3 horas

**Total estimado:** 13-18 horas (~2-3 días)

---

## Verificación

### Test: Exportar Evento a PDF

```bash
# 1. Navegar a detalle de evento
# 2. Click "Exportar PDF"
# 3. Verificar descarga automática
# 4. Abrir PDF:
- ✓ Header con título
- ✓ Información del evento
- ✓ Métricas del reporte
- ✓ Tabla de vehículos
- ✓ Lista de coordinadores
```

### Test: Exportar Dashboard a PDF

```bash
# 1. En /analytics
# 2. Click "Exportar PDF"
# 3. Verificar:
- ✓ Captura de todos los KPIs
- ✓ Gráficas como imágenes
- ✓ Layout correcto en A4
- ✓ Calidad de imagen aceptable
```

### Test: Drill-down

```bash
# 1. En Analytics, click en barra de Chevrolet
# 2. Verificar:
- ✓ Navega a /events?brand_id=...&from=analytics
- ✓ Breadcrumb aparece: Inicio > Analytics > Chevrolet
- ✓ Lista filtrada a eventos de Chevrolet
- ✓ Click en breadcrumb "Analytics" vuelve a analytics

# 3. Click en mes de timeline
- ✓ Navega con filtro de mes
- ✓ Eventos del mes específico

# 4. Click en vehículo
- ✓ Filtra eventos con ese vehículo
```

---

**FIN DEL PLAN - FASE 6**
