import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  CardContent,
  TextField,
  Typography,
  Alert,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Divider,
  Stack,
} from '@mui/material';
import { Save as SaveIcon, Cancel as CancelIcon } from '@mui/icons-material';
import { Brand, CreateEventRequest, EventWithBrand, UpdateEventRequest } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

export const EventForm = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { accessToken } = useAppSelector((state) => state.auth);
  const isEditMode = Boolean(id);

  const [brands, setBrands] = useState<Brand[]>([]);
  const [loading, setLoading] = useState(false);
  const [fetchingEvent, setFetchingEvent] = useState(isEditMode);
  const [error, setError] = useState('');

  // Form fields
  const [brandId, setBrandId] = useState('');
  const [eventType, setEventType] = useState('');
  const [organizer, setOrganizer] = useState('');
  const [name, setName] = useState('');
  const [startDate, setStartDate] = useState('');
  const [year, setYear] = useState('');
  const [durationDays, setDurationDays] = useState('');
  const [state, setState] = useState('');
  const [city, setCity] = useState('');
  const [venue, setVenue] = useState('');
  const [dealer, setDealer] = useState('');

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

  const fetchEvent = async () => {
    if (!id || !accessToken) return;

    setFetchingEvent(true);
    const response = await apiService.get<EventWithBrand>(
      API_ENDPOINTS.EVENTS.BY_ID(id),
      accessToken
    );

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      const event = response.data;
      setBrandId(event.brand_id);
      setEventType(event.event_type);
      setOrganizer(event.organizer);
      setName(event.name);
      setStartDate(event.start_date.split('T')[0]);
      setYear(event.year.toString());
      setDurationDays(event.duration_days.toString());
      setState(event.state);
      setCity(event.city);
      setVenue(event.venue);
      setDealer(event.dealer);
    }

    setFetchingEvent(false);
  };

  useEffect(() => {
    if (accessToken) {
      fetchBrands();
      if (isEditMode) {
        fetchEvent();
      }
    }
  }, [accessToken, isEditMode]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (isEditMode && id) {
      const updateData: UpdateEventRequest = {
        event_type: eventType,
        organizer,
        name,
        start_date: startDate,
        year: parseInt(year),
        duration_days: parseInt(durationDays),
        state,
        city,
        venue,
        dealer,
      };

      const response = await apiService.put(
        API_ENDPOINTS.EVENTS.UPDATE(id),
        updateData,
        accessToken || undefined
      );

      if (response.error) {
        setError(response.error);
      } else {
        navigate('/events');
      }
    } else {
      const createData: CreateEventRequest = {
        brand_id: brandId,
        event_type: eventType,
        organizer,
        name,
        start_date: startDate,
        year: parseInt(year),
        duration_days: parseInt(durationDays),
        state,
        city,
        venue,
        dealer,
      };

      const response = await apiService.post(
        API_ENDPOINTS.EVENTS.CREATE,
        createData,
        accessToken || undefined
      );

      if (response.error) {
        setError(response.error);
      } else {
        navigate('/events');
      }
    }

    setLoading(false);
  };

  const handleCancel = () => {
    navigate('/events');
  };

  if (fetchingEvent) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box maxWidth="900px" mx="auto">
      <Box mb={4}>
        <Typography variant="h4" fontWeight={600} gutterBottom>
          {isEditMode ? 'Editar Evento' : 'Nuevo Evento'}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Complete la información del evento
        </Typography>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit}>
        <Stack spacing={3}>
          {/* Información General */}
          <Card elevation={0} sx={{ border: '1px solid', borderColor: 'divider' }}>
            <CardContent sx={{ p: 4 }}>
              <Typography variant="h6" fontWeight={600} gutterBottom>
                Información General
              </Typography>
              <Divider sx={{ mb: 3 }} />

              <Stack spacing={3}>
                <FormControl fullWidth required disabled={isEditMode}>
                  <InputLabel>Marca</InputLabel>
                  <Select value={brandId} label="Marca" onChange={(e) => setBrandId(e.target.value)}>
                    {brands.map((brand) => (
                      <MenuItem key={brand.id} value={brand.id}>
                        {brand.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>

                <TextField
                  fullWidth
                  required
                  label="Nombre del Evento"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Ej: Triatlón Chevrolet Guadalajara 2026"
                />

                <Box display="flex" gap={2}>
                  <TextField
                    fullWidth
                    required
                    label="Tipo de Evento"
                    value={eventType}
                    onChange={(e) => setEventType(e.target.value)}
                    placeholder="Ej: Triatlón, Maratón, Expo"
                  />
                  <TextField
                    fullWidth
                    required
                    label="Organizador"
                    value={organizer}
                    onChange={(e) => setOrganizer(e.target.value)}
                    placeholder="Ej: Chevrolet México"
                  />
                </Box>
              </Stack>
            </CardContent>
          </Card>

          {/* Fecha y Duración */}
          <Card elevation={0} sx={{ border: '1px solid', borderColor: 'divider' }}>
            <CardContent sx={{ p: 4 }}>
              <Typography variant="h6" fontWeight={600} gutterBottom>
                Fecha y Duración
              </Typography>
              <Divider sx={{ mb: 3 }} />

              <Box display="flex" gap={2}>
                <TextField
                  fullWidth
                  required
                  type="date"
                  label="Fecha de Inicio"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  InputLabelProps={{ shrink: true }}
                />
                <TextField
                  fullWidth
                  required
                  type="number"
                  label="Año"
                  value={year}
                  onChange={(e) => setYear(e.target.value)}
                  inputProps={{ min: 2000, max: 2100 }}
                />
                <TextField
                  fullWidth
                  required
                  type="number"
                  label="Duración (días)"
                  value={durationDays}
                  onChange={(e) => setDurationDays(e.target.value)}
                  inputProps={{ min: 1 }}
                />
              </Box>
            </CardContent>
          </Card>

          {/* Ubicación */}
          <Card elevation={0} sx={{ border: '1px solid', borderColor: 'divider' }}>
            <CardContent sx={{ p: 4 }}>
              <Typography variant="h6" fontWeight={600} gutterBottom>
                Ubicación
              </Typography>
              <Divider sx={{ mb: 3 }} />

              <Stack spacing={3}>
                <Box display="flex" gap={2}>
                  <TextField
                    fullWidth
                    required
                    label="Estado"
                    value={state}
                    onChange={(e) => setState(e.target.value)}
                    placeholder="Ej: Jalisco"
                  />
                  <TextField
                    fullWidth
                    required
                    label="Ciudad"
                    value={city}
                    onChange={(e) => setCity(e.target.value)}
                    placeholder="Ej: Guadalajara"
                  />
                </Box>

                <TextField
                  fullWidth
                  required
                  label="Sede"
                  value={venue}
                  onChange={(e) => setVenue(e.target.value)}
                  placeholder="Ej: Expo Guadalajara"
                />

                <TextField
                  fullWidth
                  required
                  label="Distribuidor"
                  value={dealer}
                  onChange={(e) => setDealer(e.target.value)}
                  placeholder="Ej: Chevrolet Andares"
                />
              </Stack>
            </CardContent>
          </Card>

          {/* Botones de acción */}
          <Box display="flex" gap={2} justifyContent="flex-end" pt={2}>
            <Button
              variant="outlined"
              size="large"
              onClick={handleCancel}
              disabled={loading}
              sx={{ minWidth: 120 }}
            >
              Cancelar
            </Button>
            <Button
              type="submit"
              variant="contained"
              size="large"
              disabled={loading}
              sx={{ minWidth: 120 }}
            >
              {loading ? 'Guardando...' : 'Guardar'}
            </Button>
          </Box>
        </Stack>
      </form>
    </Box>
  );
};
