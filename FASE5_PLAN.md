# Plan de Implementación: Fase 5 - Analytics y Dashboards

## Contexto

Las Fases 1-4 del sistema CustomerMX están completas con gestión completa de eventos, usuarios y datos operacionales. La base de datos contiene información valiosa sobre eventos, asistencia, vehículos presentados, ratings, etc.

La Fase 5 implementará el sistema de analytics y visualización de datos, permitiendo a los usuarios (especialmente BRAND) obtener insights valiosos sobre el desempeño de sus eventos mediante gráficas interactivas, métricas clave y comparativas.

## Alcance de la Implementación

**Backend:**
- Endpoints de analytics (`/api/v1/analytics/*`)
- Queries optimizadas con agregaciones SQL
- Métricas calculadas: totales, promedios, agrupaciones
- Filtros por marca, año, rango de fechas
- Respuestas optimizadas para gráficas

**Frontend:**
- Biblioteca de charts (Recharts)
- Dashboard mejorado con gráficas reales
- Página Analytics dedicada con múltiples visualizaciones
- Filtros de tiempo interactivos
- Cards con estadísticas clave (KPIs)
- Componentes de gráficas reutilizables

**Tipos de Visualizaciones:**
1. Eventos por marca (bar chart)
2. Eventos por mes/timeline (line chart)
3. Ratings promedio por marca (bar chart con colores)
4. Top 10 vehículos más presentados (horizontal bar)
5. Distribución por estado (pie chart)
6. Asistencia total y promedio (KPI cards)
7. Leads y prospectos (KPI cards)
8. Comparativa año vs año

**Fuera de alcance (fases posteriores):**
- Materialized views (optimización futura)
- Exportación a PDF de gráficas
- Dashboards personalizables por usuario
- Alertas automáticas

## Arquitectura de Analytics

### Backend: Endpoints de Agregación

Los endpoints de analytics retornan datos pre-agregados para evitar procesamiento pesado en frontend:

```go
// Estructura de respuesta típica
type AnalyticsResponse struct {
    Totals    TotalMetrics       `json:"totals"`
    ByBrand   []BrandMetrics     `json:"by_brand"`
    ByMonth   []MonthlyMetrics   `json:"by_month"`
    ByState   []StateMetrics     `json:"by_state"`
    TopVehicles []VehicleMetrics `json:"top_vehicles"`
}
```

### Frontend: Recharts + Cards

Recharts es la biblioteca elegida porque:
- ✅ Integración nativa con React
- ✅ Componentes declarativos
- ✅ Responsive por defecto
- ✅ Animaciones suaves
- ✅ Tooltips interactivos
- ✅ TypeScript support

## Orden de Implementación

### ETAPA 1: Backend - Endpoints de Analytics

Crear endpoints de agregación optimizados.

#### Archivos a crear:

**`/backend/internal/domain/analytics/model.go`**

```go
package analytics

import "time"

// Métricas totales
type TotalMetrics struct {
    TotalEvents      int     `json:"total_events"`
    TotalAttendees   int     `json:"total_attendees"`
    TotalLeads       int     `json:"total_leads"`
    TotalProspects   int     `json:"total_prospects"`
    AverageAttendees float64 `json:"average_attendees"`
    AverageRating    float64 `json:"average_rating"`
}

// Métricas por marca
type BrandMetrics struct {
    BrandID          string  `json:"brand_id"`
    BrandName        string  `json:"brand_name"`
    EventCount       int     `json:"event_count"`
    TotalAttendees   int     `json:"total_attendees"`
    TotalLeads       int     `json:"total_leads"`
    AverageRating    float64 `json:"average_rating"`
}

// Métricas mensuales
type MonthlyMetrics struct {
    Year         int    `json:"year"`
    Month        int    `json:"month"`
    MonthName    string `json:"month_name"`
    EventCount   int    `json:"event_count"`
    Attendees    int    `json:"attendees"`
}

// Métricas por estado
type StateMetrics struct {
    State        string `json:"state"`
    EventCount   int    `json:"event_count"`
    Attendees    int    `json:"attendees"`
}

// Métricas de vehículos
type VehicleMetrics struct {
    VehicleID    string `json:"vehicle_id"`
    ModelName    string `json:"model_name"`
    BrandName    string `json:"brand_name"`
    TimesPresented int  `json:"times_presented"`
    TotalQuantity  int  `json:"total_quantity"`
}

// Comparativa año vs año
type YearComparison struct {
    Year           int     `json:"year"`
    EventCount     int     `json:"event_count"`
    TotalAttendees int     `json:"total_attendees"`
    AverageRating  float64 `json:"average_rating"`
}

// Request con filtros
type AnalyticsFilters struct {
    BrandID   *string    `json:"brand_id"`
    Year      *int       `json:"year"`
    StartDate *time.Time `json:"start_date"`
    EndDate   *time.Time `json:"end_date"`
}

// Response principal
type DashboardAnalytics struct {
    Totals       TotalMetrics      `json:"totals"`
    ByBrand      []BrandMetrics    `json:"by_brand"`
    ByMonth      []MonthlyMetrics  `json:"by_month"`
    ByState      []StateMetrics    `json:"by_state"`
    TopVehicles  []VehicleMetrics  `json:"top_vehicles"`
    YearComparison []YearComparison `json:"year_comparison"`
}
```

**`/backend/internal/domain/analytics/repository.go`**

Interface con métodos de agregación:

```go
type Repository interface {
    GetTotalMetrics(ctx context.Context, filters AnalyticsFilters) (*TotalMetrics, error)
    GetMetricsByBrand(ctx context.Context, filters AnalyticsFilters) ([]BrandMetrics, error)
    GetMetricsByMonth(ctx context.Context, filters AnalyticsFilters) ([]MonthlyMetrics, error)
    GetMetricsByState(ctx context.Context, filters AnalyticsFilters) ([]StateMetrics, error)
    GetTopVehicles(ctx context.Context, filters AnalyticsFilters, limit int) ([]VehicleMetrics, error)
    GetYearComparison(ctx context.Context, brandID *string) ([]YearComparison, error)
}
```

**Queries SQL críticas:**

```sql
-- Métricas totales
SELECT
    COUNT(e.id) as total_events,
    COALESCE(SUM(er.attendees), 0) as total_attendees,
    COALESCE(SUM(er.leads_collected), 0) as total_leads,
    COALESCE(SUM(er.prospects), 0) as total_prospects,
    COALESCE(AVG(er.attendees), 0) as average_attendees,
    COALESCE(AVG(er.dealer_rating), 0) as average_rating
FROM events e
LEFT JOIN event_reports er ON e.id = er.event_id
WHERE ($1::uuid IS NULL OR e.brand_id = $1)
  AND ($2::int IS NULL OR e.year = $2)
  AND e.status = 'COMPLETED';

-- Métricas por marca
SELECT
    b.id as brand_id,
    b.name as brand_name,
    COUNT(e.id) as event_count,
    COALESCE(SUM(er.attendees), 0) as total_attendees,
    COALESCE(SUM(er.leads_collected), 0) as total_leads,
    COALESCE(AVG(er.dealer_rating), 0) as average_rating
FROM brands b
LEFT JOIN events e ON b.id = e.brand_id
LEFT JOIN event_reports er ON e.id = er.event_id
WHERE e.status = 'COMPLETED'
  AND ($1::int IS NULL OR e.year = $1)
GROUP BY b.id, b.name
ORDER BY event_count DESC;

-- Métricas mensuales (timeline)
SELECT
    EXTRACT(YEAR FROM e.start_date) as year,
    EXTRACT(MONTH FROM e.start_date) as month,
    TO_CHAR(e.start_date, 'Month') as month_name,
    COUNT(e.id) as event_count,
    COALESCE(SUM(er.attendees), 0) as attendees
FROM events e
LEFT JOIN event_reports er ON e.id = er.event_id
WHERE ($1::uuid IS NULL OR e.brand_id = $1)
  AND ($2::int IS NULL OR e.year = $2)
  AND e.status = 'COMPLETED'
GROUP BY year, month, month_name
ORDER BY year, month;

-- Top vehículos
SELECT
    v.id as vehicle_id,
    v.model_name,
    b.name as brand_name,
    COUNT(DISTINCT ev.event_id) as times_presented,
    SUM(ev.quantity) as total_quantity
FROM vehicles v
JOIN brands b ON v.brand_id = b.id
JOIN event_vehicles ev ON v.id = ev.vehicle_id
JOIN events e ON ev.event_id = e.id
WHERE e.status = 'COMPLETED'
  AND ($1::uuid IS NULL OR e.brand_id = $1)
  AND ($2::int IS NULL OR e.year = $2)
GROUP BY v.id, v.model_name, b.name
ORDER BY total_quantity DESC
LIMIT $3;

-- Comparativa año vs año
SELECT
    e.year,
    COUNT(e.id) as event_count,
    COALESCE(SUM(er.attendees), 0) as total_attendees,
    COALESCE(AVG(er.dealer_rating), 0) as average_rating
FROM events e
LEFT JOIN event_reports er ON e.id = er.event_id
WHERE ($1::uuid IS NULL OR e.brand_id = $1)
  AND e.status = 'COMPLETED'
GROUP BY e.year
ORDER BY e.year;
```

**`/backend/internal/domain/analytics/service.go`**

```go
type Service interface {
    GetDashboardAnalytics(ctx context.Context, filters AnalyticsFilters) (*DashboardAnalytics, error)
    GetEventsByBrand(ctx context.Context, filters AnalyticsFilters) ([]BrandMetrics, error)
    GetEventTimeline(ctx context.Context, filters AnalyticsFilters) ([]MonthlyMetrics, error)
}

type analyticsService struct {
    repo Repository
}

func (s *analyticsService) GetDashboardAnalytics(ctx context.Context, filters AnalyticsFilters) (*DashboardAnalytics, error) {
    // Llamadas paralelas a todas las métricas
    // Combinar resultados en DashboardAnalytics
}
```

**`/backend/internal/http/handlers/analytics.go`**

```go
type AnalyticsHandler struct {
    analyticsService analytics.Service
}

// GET /api/v1/analytics/dashboard
func (h *AnalyticsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
    // Parse query params: brand_id, year, start_date, end_date
    // Call service
    // Return DashboardAnalytics
}

// GET /api/v1/analytics/events/by-brand
func (h *AnalyticsHandler) GetEventsByBrand(w http.ResponseWriter, r *http.Request) {
    // Return []BrandMetrics
}

// GET /api/v1/analytics/events/timeline
func (h *AnalyticsHandler) GetEventTimeline(w http.ResponseWriter, r *http.Request) {
    // Return []MonthlyMetrics
}
```

**Registrar rutas en router:**

```go
// Initialize analytics domain
analyticsRepo := analytics.NewRepository(dbConn.Pool)
analyticsService := analytics.NewService(analyticsRepo)
analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

// Register routes (protegidas)
r.Get("/analytics/dashboard", analyticsHandler.GetDashboard)
r.Get("/analytics/events/by-brand", analyticsHandler.GetEventsByBrand)
r.Get("/analytics/events/timeline", analyticsHandler.GetEventTimeline)
```

**Filtrado automático por rol:**

```go
func (h *AnalyticsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
    claims, _ := middleware.GetUserFromContext(r.Context())

    filters := analytics.AnalyticsFilters{}

    // BRAND users: solo su marca
    if claims.Role == "BRAND" && claims.BrandID != nil {
        filters.BrandID = claims.BrandID
    }

    // Parse otros filtros de query params...

    result, err := h.analyticsService.GetDashboardAnalytics(r.Context(), filters)
    // ...
}
```

---

### ETAPA 2: Frontend - Instalación de Recharts

```bash
npm install recharts
npm install --save-dev @types/recharts
```

---

### ETAPA 3: Frontend - Tipos TypeScript

**`/frontend/src/types/index.ts`** (agregar):

```typescript
// Analytics types
export interface TotalMetrics {
  total_events: number;
  total_attendees: number;
  total_leads: number;
  total_prospects: number;
  average_attendees: number;
  average_rating: number;
}

export interface BrandMetrics {
  brand_id: string;
  brand_name: string;
  event_count: number;
  total_attendees: number;
  total_leads: number;
  average_rating: number;
}

export interface MonthlyMetrics {
  year: number;
  month: number;
  month_name: string;
  event_count: number;
  attendees: number;
}

export interface StateMetrics {
  state: string;
  event_count: number;
  attendees: number;
}

export interface VehicleMetrics {
  vehicle_id: string;
  model_name: string;
  brand_name: string;
  times_presented: number;
  total_quantity: number;
}

export interface YearComparison {
  year: number;
  event_count: number;
  total_attendees: number;
  average_rating: number;
}

export interface DashboardAnalytics {
  totals: TotalMetrics;
  by_brand: BrandMetrics[];
  by_month: MonthlyMetrics[];
  by_state: StateMetrics[];
  top_vehicles: VehicleMetrics[];
  year_comparison: YearComparison[];
}

export interface AnalyticsFilters {
  brand_id?: string;
  year?: number;
  start_date?: string;
  end_date?: string;
}
```

---

### ETAPA 4: Frontend - API Constants

**`/frontend/src/services/api/apiConstants.ts`** (agregar):

```typescript
ANALYTICS: {
  DASHBOARD: '/analytics/dashboard',
  BY_BRAND: '/analytics/events/by-brand',
  TIMELINE: '/analytics/events/timeline',
}
```

---

### ETAPA 5: Frontend - Componentes de Gráficas Reutilizables

**`/frontend/src/components/Charts/EventsByBrandChart.tsx`**

```tsx
import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { BrandMetrics } from '../../types';

interface Props {
  data: BrandMetrics[];
}

export const EventsByBrandChart: React.FC<Props> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="brand_name" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Bar dataKey="event_count" fill="#2563eb" name="Eventos" />
      </BarChart>
    </ResponsiveContainer>
  );
};
```

**`/frontend/src/components/Charts/EventTimelineChart.tsx`**

```tsx
import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { MonthlyMetrics } from '../../types';

interface Props {
  data: MonthlyMetrics[];
}

export const EventTimelineChart: React.FC<Props> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="month_name" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Line type="monotone" dataKey="event_count" stroke="#2563eb" name="Eventos" />
        <Line type="monotone" dataKey="attendees" stroke="#10b981" name="Asistentes" />
      </LineChart>
    </ResponsiveContainer>
  );
};
```

**`/frontend/src/components/Charts/TopVehiclesChart.tsx`**

```tsx
import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { VehicleMetrics } from '../../types';

interface Props {
  data: VehicleMetrics[];
}

export const TopVehiclesChart: React.FC<Props> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart data={data} layout="vertical">
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis type="number" />
        <YAxis dataKey="model_name" type="category" width={120} />
        <Tooltip />
        <Bar dataKey="total_quantity" fill="#8b5cf6" name="Cantidad Total" />
      </BarChart>
    </ResponsiveContainer>
  );
};
```

**`/frontend/src/components/Charts/StateDistributionChart.tsx`**

```tsx
import React from 'react';
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { StateMetrics } from '../../types';

interface Props {
  data: StateMetrics[];
}

const COLORS = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899'];

export const StateDistributionChart: React.FC<Props> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={data}
          dataKey="event_count"
          nameKey="state"
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

**`/frontend/src/components/Charts/YearComparisonChart.tsx`**

```tsx
import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { YearComparison } from '../../types';

interface Props {
  data: YearComparison[];
}

export const YearComparisonChart: React.FC<Props> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="year" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Bar dataKey="event_count" fill="#2563eb" name="Eventos" />
        <Bar dataKey="total_attendees" fill="#10b981" name="Asistentes" />
      </BarChart>
    </ResponsiveContainer>
  );
};
```

---

### ETAPA 6: Frontend - KPI Cards Component

**`/frontend/src/components/KPICard/index.tsx`**

```tsx
import React from 'react';
import { Card, CardContent, Typography, Box } from '@mui/material';
import { TrendingUp, TrendingDown } from '@mui/icons-material';

interface Props {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: React.ReactElement;
  trend?: {
    direction: 'up' | 'down';
    value: string;
  };
  color?: string;
}

export const KPICard: React.FC<Props> = ({
  title,
  value,
  subtitle,
  icon,
  trend,
  color = '#2563eb',
}) => {
  return (
    <Card>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <Box>
            <Typography color="text.secondary" variant="body2" gutterBottom>
              {title}
            </Typography>
            <Typography variant="h4" sx={{ fontWeight: 600, color }}>
              {value}
            </Typography>
            {subtitle && (
              <Typography variant="caption" color="text.secondary">
                {subtitle}
              </Typography>
            )}
            {trend && (
              <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                {trend.direction === 'up' ? (
                  <TrendingUp sx={{ fontSize: 16, color: 'success.main', mr: 0.5 }} />
                ) : (
                  <TrendingDown sx={{ fontSize: 16, color: 'error.main', mr: 0.5 }} />
                )}
                <Typography variant="caption" color={trend.direction === 'up' ? 'success.main' : 'error.main'}>
                  {trend.value}
                </Typography>
              </Box>
            )}
          </Box>
          {icon && (
            <Box
              sx={{
                bgcolor: `${color}15`,
                borderRadius: 2,
                p: 1.5,
                color,
              }}
            >
              {icon}
            </Box>
          )}
        </Box>
      </CardContent>
    </Card>
  );
};
```

---

### ETAPA 7: Frontend - Página Analytics Completa

**`/frontend/src/pages/Analytics/index.tsx`**

```tsx
import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Alert,
} from '@mui/material';
import {
  Event as EventIcon,
  People as PeopleIcon,
  ContactPage as ContactPageIcon,
  Star as StarIcon,
} from '@mui/icons-material';
import { useAppSelector } from '../../hooks/useRedux';
import { useSnackbar } from '../../hooks/useSnackbar';
import { apiService } from '../../services/api/apiService';
import { API_ENDPOINTS } from '../../services/api/apiConstants';
import { DashboardAnalytics } from '../../types';
import { KPICard } from '../../components/KPICard';
import { EventsByBrandChart } from '../../components/Charts/EventsByBrandChart';
import { EventTimelineChart } from '../../components/Charts/EventTimelineChart';
import { TopVehiclesChart } from '../../components/Charts/TopVehiclesChart';
import { StateDistributionChart } from '../../components/Charts/StateDistributionChart';
import { YearComparisonChart } from '../../components/Charts/YearComparisonChart';

export const Analytics: React.FC = () => {
  const { showSnackbar } = useSnackbar();
  const currentUser = useAppSelector((state) => state.auth.user);
  const brands = useAppSelector((state) => state.auth.brands) || [];

  const [analytics, setAnalytics] = useState<DashboardAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const [yearFilter, setYearFilter] = useState<number>(new Date().getFullYear());
  const [brandFilter, setBrandFilter] = useState<string>('');

  useEffect(() => {
    fetchAnalytics();
  }, [yearFilter, brandFilter]);

  const fetchAnalytics = async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      if (yearFilter) params.append('year', yearFilter.toString());
      if (brandFilter) params.append('brand_id', brandFilter);

      const response = await apiService.get(
        `${API_ENDPOINTS.ANALYTICS.DASHBOARD}?${params.toString()}`
      );
      setAnalytics(response.data);
    } catch (err: any) {
      setError(err.message || 'Error al cargar analytics');
      showSnackbar('Error al cargar analytics', 'error');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error || !analytics) {
    return (
      <Box sx={{ p: 3 }}>
        <Alert severity="error">{error || 'No se pudieron cargar los datos'}</Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" sx={{ fontWeight: 600 }}>
          Analytics
        </Typography>
      </Box>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Año</InputLabel>
                <Select
                  value={yearFilter}
                  onChange={(e) => setYearFilter(e.target.value as number)}
                  label="Año"
                >
                  {[2024, 2025, 2026].map((year) => (
                    <MenuItem key={year} value={year}>
                      {year}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            {currentUser?.role === 'ADMIN' && (
              <Grid item xs={12} sm={6}>
                <FormControl fullWidth>
                  <InputLabel>Marca</InputLabel>
                  <Select
                    value={brandFilter}
                    onChange={(e) => setBrandFilter(e.target.value)}
                    label="Marca"
                  >
                    <MenuItem value="">Todas</MenuItem>
                    {brands.map((brand) => (
                      <MenuItem key={brand.id} value={brand.id}>
                        {brand.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            )}
          </Grid>
        </CardContent>
      </Card>

      {/* KPIs */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <KPICard
            title="Total Eventos"
            value={analytics.totals.total_events}
            icon={<EventIcon />}
            color="#2563eb"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <KPICard
            title="Total Asistentes"
            value={analytics.totals.total_attendees.toLocaleString()}
            subtitle={`Promedio: ${Math.round(analytics.totals.average_attendees)}`}
            icon={<PeopleIcon />}
            color="#10b981"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <KPICard
            title="Leads Capturados"
            value={analytics.totals.total_leads.toLocaleString()}
            icon={<ContactPageIcon />}
            color="#f59e0b"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <KPICard
            title="Rating Promedio"
            value={analytics.totals.average_rating.toFixed(1)}
            subtitle="de 5.0"
            icon={<StarIcon />}
            color="#ef4444"
          />
        </Grid>
      </Grid>

      {/* Charts */}
      <Grid container spacing={3}>
        <Grid item xs={12} lg={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Eventos por Marca
              </Typography>
              <EventsByBrandChart data={analytics.by_brand} />
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} lg={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Timeline de Eventos
              </Typography>
              <EventTimelineChart data={analytics.by_month} />
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} lg={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Top 10 Vehículos
              </Typography>
              <TopVehiclesChart data={analytics.top_vehicles.slice(0, 10)} />
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} lg={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Distribución por Estado
              </Typography>
              <StateDistributionChart data={analytics.by_state.slice(0, 6)} />
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Comparativa Año vs Año
              </Typography>
              <YearComparisonChart data={analytics.year_comparison} />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};
```

---

### ETAPA 8: Frontend - Mejorar Dashboard Existente

**`/frontend/src/pages/Dashboard/index.tsx`** (reemplazar stats estáticos con analytics reales)

Usar los mismos KPICards y algunas gráficas resumidas.

---

### ETAPA 9: Frontend - Rutas y Navegación

**`/frontend/src/App.tsx`** (agregar):

```tsx
import { Analytics } from './pages/Analytics';

// En Routes:
<Route
  path="/analytics"
  element={
    <ProtectedRoute>
      <Layout>
        <Analytics />
      </Layout>
    </ProtectedRoute>
  }
/>
```

**`/frontend/src/components/Layout/index.tsx`** (agregar item):

```tsx
import { Analytics as AnalyticsIcon } from '@mui/icons-material';

// En baseMenuItems:
{ text: 'Analytics', icon: <AnalyticsIcon />, path: '/analytics', roles: ['ADMIN', 'BRAND'] },
```

---

## Archivos a Crear/Modificar

### Backend - Nuevos (4 archivos):
1. `/backend/internal/domain/analytics/model.go`
2. `/backend/internal/domain/analytics/repository.go`
3. `/backend/internal/domain/analytics/service.go`
4. `/backend/internal/http/handlers/analytics.go`

### Backend - Modificar (1 archivo):
5. `/backend/internal/http/router/router.go`

### Frontend - Nuevos (11 archivos):
6. `/frontend/src/components/Charts/EventsByBrandChart.tsx`
7. `/frontend/src/components/Charts/EventTimelineChart.tsx`
8. `/frontend/src/components/Charts/TopVehiclesChart.tsx`
9. `/frontend/src/components/Charts/StateDistributionChart.tsx`
10. `/frontend/src/components/Charts/YearComparisonChart.tsx`
11. `/frontend/src/components/KPICard/index.tsx`
12. `/frontend/src/pages/Analytics/index.tsx`

### Frontend - Modificar (4 archivos):
13. `/frontend/src/types/index.ts`
14. `/frontend/src/services/api/apiConstants.ts`
15. `/frontend/src/App.tsx`
16. `/frontend/src/components/Layout/index.tsx`

**Total:** 16 archivos

---

## Decisiones de Diseño

### 1. Agregación en Backend vs Frontend
**Decisión:** Backend hace todas las agregaciones

**Razón:**
- Performance: PostgreSQL es más rápido que JavaScript
- Reduce payload de red (envía solo datos agregados)
- Cálculos consistentes (no dependen del cliente)

### 2. Recharts sobre Chart.js
**Decisión:** Usar Recharts

**Razón:**
- Componentes React nativos (no necesita refs)
- API declarativa (más legible)
- TypeScript support out of the box
- Responsive por defecto
- Menor bundle size

### 3. Filtros en Query Params
**Decisión:** Filtros como query params, no en body

**Razón:**
- GET request (cacheable)
- URLs compartibles
- Browser back/forward funciona
- Standard REST practice

### 4. KPI Cards Reutilizables
**Decisión:** Componente KPICard genérico

**Razón:**
- DRY (no repetir código)
- Consistencia visual
- Fácil agregar nuevos KPIs
- Props flexibles (icon, trend, color)

---

## Verificación de Analytics

### Test Backend (curl)

```bash
# 1. Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@customermx.com","password":"admin123"}' \
  | jq -r '.data.access_token')

# 2. Get dashboard analytics
curl "http://localhost:8080/api/v1/analytics/dashboard?year=2026" \
  -H "Authorization: Bearer $TOKEN" | jq

# Verificar:
# - totals con números reales
# - by_brand array con todas las marcas
# - by_month array con 12 meses
# - top_vehicles array
# - year_comparison array

# 3. Filtrar por marca
curl "http://localhost:8080/api/v1/analytics/dashboard?brand_id=CHEVROLET_UUID" \
  -H "Authorization: Bearer $TOKEN" | jq

# Verificar:
# - Solo datos de Chevrolet
```

### Test Frontend (navegador)

1. Login como ADMIN
2. Click "Analytics" en menú
3. Verificar:
   - ✓ 4 KPI cards con números reales
   - ✓ Gráfica "Eventos por Marca" con barras
   - ✓ Timeline con líneas
   - ✓ Top vehículos horizontal
   - ✓ Pie chart de estados
   - ✓ Comparativa años
4. Cambiar año en filtro → gráficas se actualizan
5. Login como BRAND → solo ve su marca
6. Responsive: probar en móvil

---

## Métricas Esperadas

**Backend:**
- 4 archivos nuevos (~600 líneas)
- 3 endpoints de analytics
- 6 queries SQL optimizadas

**Frontend:**
- 11 componentes nuevos (~1,200 líneas)
- 1 página completa
- 5 tipos de gráficas
- 1 componente KPI reutilizable

**Performance:**
- Analytics endpoint: < 500ms
- Render gráficas: < 100ms
- Bundle increment: ~50KB (Recharts)

---

## Próximas Optimizaciones (Fase 6)

1. **Materialized Views**
   - Pre-calcular analytics
   - Refresh automático
   - Queries instantáneas

2. **Caching**
   - Redis para analytics
   - TTL de 5 minutos
   - Invalidación manual

3. **Drill-down**
   - Click en gráfica → filtrar
   - Detalles por marca/vehículo
   - Breadcrumbs de navegación

4. **Exportación**
   - PDF con gráficas
   - Excel con datos raw
   - Compartir por email

---

**FIN DEL PLAN - FASE 5**
