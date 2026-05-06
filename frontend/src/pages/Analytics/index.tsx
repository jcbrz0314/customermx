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
import { DashboardAnalytics, Brand, EventWithBrand } from '../../types';
import { AttendanceByVenueChart } from '../../components/Charts/AttendanceByVenueChart';
import { AttendanceByYearChart } from '../../components/Charts/AttendanceByYearChart';
import { LeadsProspectsByVenueChart } from '../../components/Charts/LeadsProspectsByVenueChart';
import { TopVehiclesChart } from '../../components/Charts/TopVehiclesChart';

export const Analytics: React.FC = () => {
  const { accessToken, user } = useAppSelector((state) => state.auth);

  const [analytics, setAnalytics] = useState<DashboardAnalytics | null>(null);
  const [brands, setBrands] = useState<Brand[]>([]);
  const [availableEventTypes, setAvailableEventTypes] = useState<string[]>([]);
  const [availableOrganizers, setAvailableOrganizers] = useState<string[]>([]);
  const [availableSetupVendors, setAvailableSetupVendors] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const [brandFilter, setBrandFilter] = useState<string>('');
  const [yearFilter, setYearFilter] = useState<string>('');
  const [eventTypeFilter, setEventTypeFilter] = useState<string>('');
  const [organizerFilter, setOrganizerFilter] = useState<string>('');
  const [setupVendorFilter, setSetupVendorFilter] = useState<string>('');

  // Load brands and events on mount
  useEffect(() => {
    const fetchFilterData = async () => {
      if (!accessToken) return;

      const [brandsRes, eventsRes, vendorsRes] = await Promise.all([
        apiService.get<Brand[]>(API_ENDPOINTS.BRANDS.LIST, accessToken),
        apiService.get<EventWithBrand[]>(API_ENDPOINTS.EVENTS.LIST, accessToken),
        apiService.get<string[]>(API_ENDPOINTS.ANALYTICS.SETUP_VENDORS, accessToken),
      ]);

      if (brandsRes.data) {
        setBrands(Array.isArray(brandsRes.data) ? brandsRes.data : []);
      }

      if (eventsRes.data) {
        const events: EventWithBrand[] = Array.isArray(eventsRes.data) ? eventsRes.data : [];
        const eventTypes = [...new Set(events.map((e) => e.event_type).filter(Boolean))].sort();
        const organizers = [...new Set(events.map((e) => e.organizer).filter(Boolean))].sort();
        setAvailableEventTypes(eventTypes);
        setAvailableOrganizers(organizers);
      }

      if (vendorsRes.data) {
        setAvailableSetupVendors(Array.isArray(vendorsRes.data) ? vendorsRes.data : []);
      }
    };
    fetchFilterData();
  }, [accessToken]);

  useEffect(() => {
    fetchAnalytics();
  }, [brandFilter, yearFilter, eventTypeFilter, organizerFilter, setupVendorFilter]);

  const fetchAnalytics = async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      if (brandFilter) params.append('brand_id', brandFilter);
      if (yearFilter) params.append('year', yearFilter);
      if (eventTypeFilter) params.append('event_type', eventTypeFilter);
      if (organizerFilter) params.append('organizer', organizerFilter);
      if (setupVendorFilter) params.append('setup_vendor', setupVendorFilter);

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
            {user?.role !== 'BRAND' && (
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
            )}
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
            <FormControl sx={{ flex: 1, minWidth: 150 }}>
              <InputLabel>Tipo de Evento</InputLabel>
              <Select
                value={eventTypeFilter}
                onChange={(e) => setEventTypeFilter(e.target.value)}
                label="Tipo de Evento"
              >
                <MenuItem value="">Todos</MenuItem>
                {availableEventTypes.map((et) => (
                  <MenuItem key={et} value={et}>
                    {et}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <FormControl sx={{ flex: 1, minWidth: 150 }}>
              <InputLabel>Organizador</InputLabel>
              <Select
                value={organizerFilter}
                onChange={(e) => setOrganizerFilter(e.target.value)}
                label="Organizador"
              >
                <MenuItem value="">Todos</MenuItem>
                {availableOrganizers.map((org) => (
                  <MenuItem key={org} value={org}>
                    {org}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            {availableSetupVendors.length > 0 && (
              <FormControl sx={{ flex: 1, minWidth: 150 }}>
                <InputLabel>Empresa de Montaje</InputLabel>
                <Select
                  value={setupVendorFilter}
                  onChange={(e) => setSetupVendorFilter(e.target.value)}
                  label="Empresa de Montaje"
                >
                  <MenuItem value="">Todas</MenuItem>
                  {availableSetupVendors.map((vendor) => (
                    <MenuItem key={vendor} value={vendor}>
                      {vendor}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}
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
            Intenta cambiar el año, marca, tipo de evento u organizador.
          </Alert>
        )}
      </Stack>
    </Box>
  );
};
