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
import { useNavigate } from 'react-router-dom';
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
import { EventTypeChart } from '../../components/Charts/EventTypeChart';
import { DealerRankingChart } from '../../components/Charts/DealerRankingChart';
import { ConversionFunnel } from '../../components/ConversionFunnel';
import { ExportButton } from '../../components/ExportButton';
import {
  AnalyticsBreadcrumbs,
  BreadcrumbFilter,
} from '../../components/AnalyticsBreadcrumbs';
import { exportDashboardToPDF } from '../../utils/pdfExport';

export const Analytics: React.FC = () => {
  const navigate = useNavigate();
  const { showSnackbar } = useSnackbar();
  const currentUser = useAppSelector((state) => state.auth.user);
  const brands = useAppSelector((state) => state.auth.brands) || [];

  const [analytics, setAnalytics] = useState<DashboardAnalytics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const currentYear = new Date().getFullYear();
  const [yearFilter, setYearFilter] = useState<number>(currentYear);
  const [brandFilter, setBrandFilter] = useState<string>('');

  // Breadcrumbs for drill-down navigation
  const [breadcrumbs, setBreadcrumbs] = useState<BreadcrumbFilter[]>([]);

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

  // Drill-down handlers
  const handleBrandClick = (brandId: string, brandName: string) => {
    setBrandFilter(brandId);
    setBreadcrumbs([...breadcrumbs, { type: 'brand', label: brandName, value: brandId }]);
    showSnackbar(`Filtrando por marca: ${brandName}`, 'info');
  };

  const handleYearClick = (year: number) => {
    setYearFilter(year);
    setBreadcrumbs([
      ...breadcrumbs,
      { type: 'year', label: year.toString(), value: year.toString() },
    ]);
    showSnackbar(`Filtrando por año: ${year}`, 'info');
  };

  const handleMonthClick = (year: number, month: number, monthName: string) => {
    // Navigate to events page with filters
    navigate(`/events?year=${year}&month=${month}`);
  };

  const handleStateClick = (state: string) => {
    // Navigate to events page with state filter
    navigate(`/events?state=${encodeURIComponent(state)}`);
  };

  const handleVehicleClick = (vehicleId: string, modelName: string) => {
    // Navigate to events page with vehicle filter (requires backend support)
    navigate(`/events?vehicle=${vehicleId}`);
  };

  const handleEventTypeClick = (eventType: string) => {
    // Navigate to events page with event type filter
    navigate(`/events?type=${encodeURIComponent(eventType)}`);
  };

  const handleDealerClick = (dealer: string) => {
    // Navigate to events page with dealer filter
    navigate(`/events?dealer=${encodeURIComponent(dealer)}`);
  };

  // Breadcrumbs handlers
  const handleRemoveFilter = (index: number) => {
    const newBreadcrumbs = breadcrumbs.filter((_, i) => i !== index);
    setBreadcrumbs(newBreadcrumbs);

    // Reset filters based on removed breadcrumb
    const removed = breadcrumbs[index];
    if (removed.type === 'brand') {
      setBrandFilter('');
    } else if (removed.type === 'year') {
      setYearFilter(currentYear);
    }
  };

  const handleClearAllFilters = () => {
    setBreadcrumbs([]);
    setBrandFilter('');
    setYearFilter(currentYear);
    showSnackbar('Filtros limpiados', 'success');
  };

  // Export handler
  const handleExport = async () => {
    if (!analytics) return;

    await exportDashboardToPDF(analytics, {
      year: yearFilter,
      brandId: brandFilter,
    });

    showSnackbar('Dashboard exportado exitosamente', 'success');
  };

  const availableYears = [2024, 2025, 2026];

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
      <Box
        sx={{
          mb: 3,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <Typography variant="h4" sx={{ fontWeight: 600 }}>
          Analytics
        </Typography>
        <ExportButton onExport={handleExport} label="Exportar Dashboard a PDF" />
      </Box>

      {/* Breadcrumbs for drill-down */}
      <AnalyticsBreadcrumbs
        filters={breadcrumbs}
        onRemoveFilter={handleRemoveFilter}
        onClearAll={handleClearAllFilters}
      />

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
                  {availableYears.map((year) => (
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
            subtitle={`Prospectos: ${analytics.totals.total_prospects.toLocaleString()}`}
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
        {/* Eventos por Marca */}
        {analytics.by_brand.length > 0 && (
          <Grid item xs={12} lg={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                  Eventos por Marca
                </Typography>
                <EventsByBrandChart
                  data={analytics.by_brand}
                  onBrandClick={currentUser?.role === 'ADMIN' ? handleBrandClick : undefined}
                />
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Timeline de Eventos */}
        {analytics.by_month.length > 0 && (
          <Grid item xs={12} lg={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                  Timeline de Eventos
                </Typography>
                <EventTimelineChart
                  data={analytics.by_month}
                  onMonthClick={handleMonthClick}
                />
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Funnel de Conversión */}
        <Grid item xs={12} lg={6}>
          <Card>
            <CardContent>
              <ConversionFunnel data={analytics.conversion} />
            </CardContent>
          </Card>
        </Grid>

        {/* Distribución por Tipo de Evento */}
        {analytics.by_event_type.length > 0 && (
          <Grid item xs={12} lg={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                  Distribución por Tipo de Evento
                </Typography>
                <EventTypeChart
                  data={analytics.by_event_type}
                  onEventTypeClick={handleEventTypeClick}
                />
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Top 10 Vehículos */}
        {analytics.top_vehicles.length > 0 && (
          <Grid item xs={12} lg={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                  Top 10 Vehículos Más Presentados
                </Typography>
                <TopVehiclesChart
                  data={analytics.top_vehicles.slice(0, 10)}
                  onVehicleClick={handleVehicleClick}
                />
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Distribución por Estado */}
        {analytics.by_state.length > 0 && (
          <Grid item xs={12} lg={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                  Distribución por Estado
                </Typography>
                <StateDistributionChart
                  data={analytics.by_state.slice(0, 6)}
                  onStateClick={handleStateClick}
                />
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Top Distribuidores */}
        {analytics.top_dealers.length > 0 && (
          <Grid item xs={12} lg={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                  Top Distribuidores por Rating
                </Typography>
                <DealerRankingChart
                  data={analytics.top_dealers}
                  onDealerClick={handleDealerClick}
                />
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Comparativa Año vs Año */}
        {analytics.year_comparison.length > 0 && (
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                  Comparativa Año vs Año
                </Typography>
                <YearComparisonChart
                  data={analytics.year_comparison}
                  onYearClick={handleYearClick}
                />
              </CardContent>
            </Card>
          </Grid>
        )}

        {/* Empty state si no hay datos */}
        {analytics.totals.total_events === 0 && (
          <Grid item xs={12}>
            <Alert severity="info">
              No hay datos disponibles para los filtros seleccionados.
              Intenta cambiar el año o la marca.
            </Alert>
          </Grid>
        )}
      </Grid>
    </Box>
  );
};
