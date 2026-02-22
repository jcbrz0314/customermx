import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Alert,
  Stack,
} from '@mui/material';
import { useAppSelector } from '../../hooks/useRedux';
import { apiService } from '../../services/api/apiService';
import { API_ENDPOINTS } from '../../services/api/apiConstants';
import { DashboardAnalytics, Brand } from '../../types';
import { AttendanceByVenueChart } from '../../components/Charts/AttendanceByVenueChart';
import { AttendanceByYearChart } from '../../components/Charts/AttendanceByYearChart';
import { LeadsProspectsByVenueChart } from '../../components/Charts/LeadsProspectsByVenueChart';
import { TopVehiclesChart } from '../../components/Charts/TopVehiclesChart';

export const Analytics: React.FC = () => {
  const accessToken = useAppSelector((state) => state.auth.accessToken);

  const [analytics, setAnalytics] = useState<DashboardAnalytics | null>(null);
  const [brands, setBrands] = useState<Brand[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters - brand is primary
  const [brandFilter, setBrandFilter] = useState<string>('');
  const [yearFilter, setYearFilter] = useState<string>('');

  // Load brands on mount
  useEffect(() => {
    const fetchBrands = async () => {
      if (!accessToken) return;
      const response = await apiService.get<Brand[]>(API_ENDPOINTS.BRANDS.LIST, accessToken);
      if (response.data) {
        setBrands(Array.isArray(response.data) ? response.data : []);
      }
    };
    fetchBrands();
  }, [accessToken]);

  useEffect(() => {
    fetchAnalytics();
  }, [brandFilter, yearFilter]);

  const fetchAnalytics = async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      if (brandFilter) params.append('brand_id', brandFilter);
      if (yearFilter) params.append('year', yearFilter);

      const response = await apiService.get(
        `${API_ENDPOINTS.ANALYTICS.DASHBOARD}?${params.toString()}`,
        accessToken || undefined
      );
      setAnalytics(response.data);
    } catch (err: any) {
      setError(err.message || 'Error al cargar analytics');
    } finally {
      setLoading(false);
    }
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
    <Box>
      {/* Header */}
      <Typography variant="h4" sx={{ fontWeight: 700, mb: 3 }}>
        Analytics
      </Typography>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
            <FormControl sx={{ flex: 1, minWidth: 150 }}>
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
            <FormControl sx={{ flex: 1, minWidth: 150 }}>
              <InputLabel>Año</InputLabel>
              <Select
                value={yearFilter}
                onChange={(e) => setYearFilter(e.target.value as string)}
                label="Año"
              >
                <MenuItem value="">Todos</MenuItem>
                {availableYears.map((year) => (
                  <MenuItem key={year} value={year.toString()}>
                    {year}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
        </CardContent>
      </Card>

      {/* Charts - stacked vertically */}
      <Stack spacing={3}>
        {/* Asistencia por Sede */}
        {analytics.by_venue && analytics.by_venue.length > 0 && (
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                Asistencia por Sede
              </Typography>
              <AttendanceByVenueChart data={analytics.by_venue} />
            </CardContent>
          </Card>
        )}

        {/* Asistencia por Año */}
        {analytics.year_comparison && analytics.year_comparison.length > 0 && (
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                Asistencia por Año
              </Typography>
              <AttendanceByYearChart data={analytics.year_comparison} />
            </CardContent>
          </Card>
        )}

        {/* Datos Levantados y Prospectos por Sede */}
        {analytics.by_venue && analytics.by_venue.length > 0 && (
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                Datos Levantados y Prospectos por Sede
              </Typography>
              <LeadsProspectsByVenueChart data={analytics.by_venue} />
            </CardContent>
          </Card>
        )}

        {/* Unidades Exhibidas */}
        {analytics.top_vehicles && analytics.top_vehicles.length > 0 && (
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
                Unidades Exhibidas
              </Typography>
              <TopVehiclesChart data={analytics.top_vehicles} />
            </CardContent>
          </Card>
        )}

        {/* Empty state */}
        {analytics.totals.total_events === 0 && (
          <Alert severity="info">
            No hay datos disponibles para los filtros seleccionados.
            Intenta cambiar el año o la marca.
          </Alert>
        )}
      </Stack>
    </Box>
  );
};
