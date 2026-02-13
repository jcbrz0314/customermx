import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  CardContent,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  CircularProgress,
  Alert,
  Chip,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
} from '@mui/material';
import { Add as AddIcon, Visibility as ViewIcon } from '@mui/icons-material';
import { format } from 'date-fns';
import { EventWithBrand, Brand, EventStatus } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

const statusColors: Record<EventStatus, 'default' | 'primary' | 'success' | 'error'> = {
  PLANNED: 'default',
  ACTIVE: 'primary',
  COMPLETED: 'success',
  CLOSED: 'error',
};

const statusLabels: Record<EventStatus, string> = {
  PLANNED: 'Planeado',
  ACTIVE: 'Activo',
  COMPLETED: 'Completado',
  CLOSED: 'Cerrado',
};

export const Events = () => {
  const navigate = useNavigate();
  const { accessToken, user } = useAppSelector((state) => state.auth);

  const [events, setEvents] = useState<EventWithBrand[]>([]);
  const [brands, setBrands] = useState<Brand[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // Filters
  const [brandFilter, setBrandFilter] = useState('');
  const [yearFilter, setYearFilter] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [stateFilter, setStateFilter] = useState('');

  const fetchBrands = async () => {
    if (!accessToken) return;

    const response = await apiService.get<Brand[]>(
      API_ENDPOINTS.BRANDS.LIST,
      accessToken
    );

    if (response.data) {
      setBrands(response.data);
    }
  };

  const fetchEvents = async () => {
    if (!accessToken) return;

    setLoading(true);
    setError('');

    // Build query string with filters
    const params = new URLSearchParams();
    if (brandFilter) params.append('brand_id', brandFilter);
    if (yearFilter) params.append('year', yearFilter);
    if (statusFilter) params.append('status', statusFilter);
    if (stateFilter) params.append('state', stateFilter);

    const endpoint = `${API_ENDPOINTS.EVENTS.LIST}${params.toString() ? `?${params.toString()}` : ''}`;

    const response = await apiService.get<EventWithBrand[]>(
      endpoint,
      accessToken
    );

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      setEvents(response.data);
    }

    setLoading(false);
  };

  useEffect(() => {
    if (accessToken) {
      fetchBrands();
    }
  }, [accessToken]);

  useEffect(() => {
    if (accessToken) {
      fetchEvents();
    }
  }, [accessToken, brandFilter, yearFilter, statusFilter, stateFilter]);

  const handleViewDetails = (eventId: string) => {
    navigate(`/events/${eventId}`);
  };

  const handleNewEvent = () => {
    navigate('/events/new');
  };

  if (loading && events.length === 0) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Eventos</Typography>
        {user?.role === 'ADMIN' && (
          <Button variant="contained" startIcon={<AddIcon />} onClick={handleNewEvent}>
            Nuevo Evento
          </Button>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      {/* Filters Card */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Filtros
          </Typography>
          <Grid container spacing={2}>
            {user?.role !== 'BRAND' && (
              <Grid xs={12} sm={6} md={3}>
                <FormControl fullWidth size="small">
                  <InputLabel>Marca</InputLabel>
                  <Select
                    value={brandFilter}
                    label="Marca"
                    onChange={(e) => setBrandFilter(e.target.value)}
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
            <Grid xs={12} sm={6} md={3}>
              <TextField
                fullWidth
                size="small"
                type="number"
                label="Año"
                value={yearFilter}
                onChange={(e) => setYearFilter(e.target.value)}
                placeholder="Ej: 2026"
              />
            </Grid>
            <Grid xs={12} sm={6} md={3}>
              <FormControl fullWidth size="small">
                <InputLabel>Estado</InputLabel>
                <Select
                  value={statusFilter}
                  label="Estado"
                  onChange={(e) => setStatusFilter(e.target.value)}
                >
                  <MenuItem value="">Todos</MenuItem>
                  <MenuItem value="PLANNED">Planeado</MenuItem>
                  <MenuItem value="ACTIVE">Activo</MenuItem>
                  <MenuItem value="COMPLETED">Completado</MenuItem>
                  <MenuItem value="CLOSED">Cerrado</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid xs={12} sm={6} md={3}>
              <TextField
                fullWidth
                size="small"
                label="Ubicación (Estado)"
                value={stateFilter}
                onChange={(e) => setStateFilter(e.target.value)}
                placeholder="Ej: Jalisco"
              />
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Events Table */}
      <Card>
        <CardContent>
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell><strong>Nombre</strong></TableCell>
                  <TableCell><strong>Tipo</strong></TableCell>
                  <TableCell><strong>Marca</strong></TableCell>
                  <TableCell><strong>Fecha</strong></TableCell>
                  <TableCell><strong>Ubicación</strong></TableCell>
                  <TableCell><strong>Estado</strong></TableCell>
                  <TableCell align="right"><strong>Acciones</strong></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {events.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} align="center">
                      No hay eventos registrados
                    </TableCell>
                  </TableRow>
                ) : (
                  events.map((event) => (
                    <TableRow key={event.id} hover>
                      <TableCell>{event.name}</TableCell>
                      <TableCell>{event.event_type}</TableCell>
                      <TableCell>
                        <Chip label={event.brand_name} color="primary" size="small" />
                      </TableCell>
                      <TableCell>
                        {format(new Date(event.start_date), 'dd/MM/yyyy')}
                      </TableCell>
                      <TableCell>{`${event.city}, ${event.state}`}</TableCell>
                      <TableCell>
                        <Chip
                          label={statusLabels[event.status]}
                          color={statusColors[event.status]}
                          size="small"
                        />
                      </TableCell>
                      <TableCell align="right">
                        <Button
                          size="small"
                          startIcon={<ViewIcon />}
                          onClick={() => handleViewDetails(event.id)}
                        >
                          Ver Detalles
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    </Box>
  );
};
