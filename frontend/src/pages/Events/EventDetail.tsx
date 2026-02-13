import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  Typography,
  CircularProgress,
  Alert,
  Chip,
  List,
  ListItem,
  ListItemText,
  Divider,
} from '@mui/material';
import {
  Edit as EditIcon,
  ArrowBack as BackIcon,
  Add as AddIcon,
} from '@mui/icons-material';
import { format } from 'date-fns';
import {
  EventWithBrand,
  EventCoordinatorWithUser,
  EventVehicleWithDetails,
  EventReport,
  EventStatus,
} from '../../types';
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

export const EventDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { accessToken, user } = useAppSelector((state) => state.auth);

  const [event, setEvent] = useState<EventWithBrand | null>(null);
  const [coordinators, setCoordinators] = useState<EventCoordinatorWithUser[]>([]);
  const [vehicles, setVehicles] = useState<EventVehicleWithDetails[]>([]);
  const [report, setReport] = useState<EventReport | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchEventData = async () => {
    if (!id || !accessToken) return;

    setLoading(true);
    setError('');

    try {
      // Fetch all data in parallel
      const [eventRes, coordinatorsRes, vehiclesRes, reportRes] = await Promise.all([
        apiService.get<EventWithBrand>(API_ENDPOINTS.EVENTS.BY_ID(id), accessToken),
        apiService.get<EventCoordinatorWithUser[]>(
          API_ENDPOINTS.EVENTS.COORDINATORS(id),
          accessToken
        ),
        apiService.get<EventVehicleWithDetails[]>(
          API_ENDPOINTS.EVENTS.VEHICLES(id),
          accessToken
        ),
        apiService.get<EventReport>(API_ENDPOINTS.EVENTS.REPORT(id), accessToken),
      ]);

      if (eventRes.error) {
        setError(eventRes.error);
      } else if (eventRes.data) {
        setEvent(eventRes.data);
      }

      if (coordinatorsRes.data) {
        setCoordinators(coordinatorsRes.data);
      }

      if (vehiclesRes.data) {
        setVehicles(vehiclesRes.data);
      }

      if (reportRes.data) {
        setReport(reportRes.data);
      }
    } catch (err) {
      setError('Error al cargar los datos del evento');
    }

    setLoading(false);
  };

  useEffect(() => {
    if (accessToken && id) {
      fetchEventData();
    }
  }, [id, accessToken]);

  const handleEdit = () => {
    navigate(`/events/${id}/edit`);
  };

  const handleBack = () => {
    navigate('/events');
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  if (!event) {
    return (
      <Box>
        <Alert severity="error">Evento no encontrado</Alert>
        <Button onClick={handleBack} sx={{ mt: 2 }}>
          Volver
        </Button>
      </Box>
    );
  }

  const isAdmin = user?.role === 'ADMIN';
  const isBrand = user?.role === 'BRAND';

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Box display="flex" alignItems="center" gap={2}>
          <Button startIcon={<BackIcon />} onClick={handleBack}>
            Volver
          </Button>
          <Typography variant="h4">{event.name}</Typography>
        </Box>
        {isAdmin && (
          <Button variant="contained" startIcon={<EditIcon />} onClick={handleEdit}>
            Editar
          </Button>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Grid container spacing={3}>
        {/* Event Details Card */}
        <Grid xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Detalles del Evento
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <Grid container spacing={2}>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Marca
                  </Typography>
                  <Typography variant="body1">{event.brand_name}</Typography>
                </Grid>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Tipo
                  </Typography>
                  <Typography variant="body1">{event.event_type}</Typography>
                </Grid>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Organizador
                  </Typography>
                  <Typography variant="body1">{event.organizer}</Typography>
                </Grid>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Estado
                  </Typography>
                  <Chip
                    label={statusLabels[event.status]}
                    color={statusColors[event.status]}
                    size="small"
                  />
                </Grid>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Fecha de Inicio
                  </Typography>
                  <Typography variant="body1">
                    {format(new Date(event.start_date), 'dd/MM/yyyy')}
                  </Typography>
                </Grid>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Duración
                  </Typography>
                  <Typography variant="body1">{event.duration_days} días</Typography>
                </Grid>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Ubicación
                  </Typography>
                  <Typography variant="body1">{`${event.city}, ${event.state}`}</Typography>
                </Grid>
                <Grid xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Sede
                  </Typography>
                  <Typography variant="body1">{event.venue}</Typography>
                </Grid>
                <Grid xs={12}>
                  <Typography variant="body2" color="text.secondary">
                    Distribuidor
                  </Typography>
                  <Typography variant="body1">{event.dealer}</Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Coordinators Card */}
        <Grid xs={12} md={6}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6">Coordinadores</Typography>
                {isAdmin && (
                  <Button size="small" startIcon={<AddIcon />}>
                    Asignar
                  </Button>
                )}
              </Box>
              <Divider sx={{ mb: 2 }} />
              {coordinators.length === 0 ? (
                <Typography variant="body2" color="text.secondary">
                  No hay coordinadores asignados
                </Typography>
              ) : (
                <List dense>
                  {coordinators.map((coordinator) => (
                    <ListItem key={coordinator.id}>
                      <ListItemText
                        primary={coordinator.user_name}
                        secondary={coordinator.user_email}
                      />
                    </ListItem>
                  ))}
                </List>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Vehicles Card */}
        <Grid xs={12} md={6}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6">Vehículos Presentados</Typography>
                {isAdmin && (
                  <Button size="small" startIcon={<AddIcon />}>
                    Agregar
                  </Button>
                )}
              </Box>
              <Divider sx={{ mb: 2 }} />
              {vehicles.length === 0 ? (
                <Typography variant="body2" color="text.secondary">
                  No hay vehículos agregados
                </Typography>
              ) : (
                <List dense>
                  {vehicles.map((vehicle) => (
                    <ListItem key={vehicle.id}>
                      <ListItemText
                        primary={vehicle.model_name}
                        secondary={`Cantidad: ${vehicle.quantity}`}
                      />
                    </ListItem>
                  ))}
                </List>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Report Card - Only visible for ADMIN and BRAND */}
        {(isAdmin || isBrand) && (
          <Grid xs={12} md={6}>
            <Card>
              <CardContent>
                <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                  <Typography variant="h6">Reporte del Evento</Typography>
                  {isAdmin && (
                    <Button size="small" startIcon={<EditIcon />}>
                      Editar Reporte
                    </Button>
                  )}
                </Box>
                <Divider sx={{ mb: 2 }} />
                {!report ? (
                  <Typography variant="body2" color="text.secondary">
                    No hay reporte disponible
                  </Typography>
                ) : (
                  <Grid container spacing={2}>
                    {report.attendees && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Asistentes
                        </Typography>
                        <Typography variant="body1">{report.attendees}</Typography>
                      </Grid>
                    )}
                    {report.leads_collected && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Leads Capturados
                        </Typography>
                        <Typography variant="body1">{report.leads_collected}</Typography>
                      </Grid>
                    )}
                    {report.prospects && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Prospectos
                        </Typography>
                        <Typography variant="body1">{report.prospects}</Typography>
                      </Grid>
                    )}
                    {report.dealer_rating && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Calificación Distribuidor
                        </Typography>
                        <Typography variant="body1">{report.dealer_rating}/5</Typography>
                      </Grid>
                    )}
                    {report.comments && (
                      <Grid xs={12}>
                        <Typography variant="body2" color="text.secondary">
                          Comentarios
                        </Typography>
                        <Typography variant="body1">{report.comments}</Typography>
                      </Grid>
                    )}
                    <Grid xs={12}>
                      <Chip
                        label={report.completed ? 'Completado' : 'En Progreso'}
                        color={report.completed ? 'success' : 'default'}
                        size="small"
                      />
                    </Grid>
                  </Grid>
                )}
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    </Box>
  );
};
