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
  IconButton,
} from '@mui/material';
import {
  Edit as EditIcon,
  ArrowBack as BackIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
  ChangeCircle as ChangeStatusIcon,
} from '@mui/icons-material';
import { format } from 'date-fns';
import {
  EventWithBrand,
  EventCoordinatorWithUser,
  EventVehicleWithDetails,
  EventReport,
  EventPhoto,
  EventStatus,
} from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';
import { useSnackbar } from '../../hooks/useSnackbar';
import { EventPhotosSection } from '../../components/EventPhotos/EventPhotosSection';
import { AssignCoordinatorDialog } from '../../components/Dialogs/AssignCoordinatorDialog';
import { AddVehicleDialog } from '../../components/Dialogs/AddVehicleDialog';
import { EditVehicleQuantityDialog } from '../../components/Dialogs/EditVehicleQuantityDialog';
import { ChangeStatusDialog } from '../../components/Dialogs/ChangeStatusDialog';
import { EventReportDialog } from '../../components/Dialogs/EventReportDialog';
import { ConfirmDialog } from '../../components/Dialogs/ConfirmDialog';
import { ExportButton } from '../../components/ExportButton';
import { exportEventToPDF } from '../../utils/pdfExport';

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
  const { showSnackbar } = useSnackbar();

  const [event, setEvent] = useState<EventWithBrand | null>(null);
  const [coordinators, setCoordinators] = useState<EventCoordinatorWithUser[]>([]);
  const [vehicles, setVehicles] = useState<EventVehicleWithDetails[]>([]);
  const [report, setReport] = useState<EventReport | null>(null);
  const [photos, setPhotos] = useState<EventPhoto[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // Dialog states
  const [assignCoordinatorOpen, setAssignCoordinatorOpen] = useState(false);
  const [addVehicleOpen, setAddVehicleOpen] = useState(false);
  const [editQuantityOpen, setEditQuantityOpen] = useState(false);
  const [selectedVehicle, setSelectedVehicle] = useState<EventVehicleWithDetails | null>(null);
  const [changeStatusOpen, setChangeStatusOpen] = useState(false);
  const [reportDialogOpen, setReportDialogOpen] = useState(false);
  const [confirmDialog, setConfirmDialog] = useState<{
    open: boolean;
    title: string;
    message: string;
    onConfirm: () => void;
    severity?: 'warning' | 'error' | 'info';
  } | null>(null);

  const fetchEventData = async () => {
    if (!id || !accessToken) return;

    setLoading(true);
    setError('');

    try {
      // Fetch all data in parallel
      const [eventRes, coordinatorsRes, vehiclesRes, reportRes, photosRes] = await Promise.all([
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
        apiService.get<EventPhoto[]>(API_ENDPOINTS.EVENTS.PHOTOS(id), accessToken),
      ]);

      if (eventRes.error) {
        setError(eventRes.error);
      } else if (eventRes.data) {
        setEvent(eventRes.data);
      }

      setCoordinators(Array.isArray(coordinatorsRes.data) ? coordinatorsRes.data : []);
      setVehicles(Array.isArray(vehiclesRes.data) ? vehiclesRes.data : []);

      if (reportRes.data) {
        setReport(reportRes.data);
      }

      setPhotos(Array.isArray(photosRes.data) ? photosRes.data : []);
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

  // Coordinator handlers
  const handleAssignCoordinator = () => {
    setAssignCoordinatorOpen(true);
  };

  const handleRemoveCoordinator = async (userId: string, userName: string) => {
    if (!accessToken || !id) return;

    setConfirmDialog({
      open: true,
      title: 'Remover Coordinador',
      message: `¿Está seguro de remover a ${userName} de este evento?`,
      severity: 'warning',
      onConfirm: async () => {
        const response = await apiService.delete(
          API_ENDPOINTS.EVENTS.REMOVE_COORDINATOR(id, userId),
          accessToken
        );

        if (response.error) {
          showSnackbar(response.error, 'error');
        } else {
          showSnackbar('Coordinador removido exitosamente', 'success');
          fetchEventData();
        }
        setConfirmDialog(null);
      },
    });
  };

  // Vehicle handlers
  const handleAddVehicle = () => {
    setAddVehicleOpen(true);
  };

  const handleEditVehicleQuantity = (vehicle: EventVehicleWithDetails) => {
    setSelectedVehicle(vehicle);
    setEditQuantityOpen(true);
  };

  const handleRemoveVehicle = async (vehicleId: string, vehicleName: string) => {
    if (!accessToken || !id) return;

    setConfirmDialog({
      open: true,
      title: 'Remover Vehículo',
      message: `¿Está seguro de remover ${vehicleName} de este evento?`,
      severity: 'warning',
      onConfirm: async () => {
        const response = await apiService.delete(
          API_ENDPOINTS.EVENTS.REMOVE_VEHICLE(id, vehicleId),
          accessToken
        );

        if (response.error) {
          showSnackbar(response.error, 'error');
        } else {
          showSnackbar('Vehículo removido exitosamente', 'success');
          fetchEventData();
        }
        setConfirmDialog(null);
      },
    });
  };

  // Status handler
  const handleChangeStatus = () => {
    setChangeStatusOpen(true);
  };

  // Report handler
  const handleEditReport = () => {
    setReportDialogOpen(true);
  };

  // Export handler
  const handleExport = async () => {
    if (!event) return;
    await exportEventToPDF(event, 'event-detail-content');
    showSnackbar('Evento exportado exitosamente a PDF', 'success');
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
  const isCoordinator = user?.role === 'COORDINATOR';
  const isBrand = user?.role === 'BRAND';

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={4}>
        <Box>
          <Button
            startIcon={<BackIcon />}
            onClick={handleBack}
            sx={{ mb: 1.5, color: 'text.secondary' }}
          >
            Volver a Eventos
          </Button>
          <Typography variant="h4" sx={{ fontWeight: 700 }}>
            {event.name}
          </Typography>
        </Box>
        <Box display="flex" gap={2}>
          <ExportButton
            onExport={handleExport}
            label="Exportar PDF"
            variant="outlined"
            size="large"
          />
          {isAdmin && (
            <>
              <Button
                variant="outlined"
                startIcon={<ChangeStatusIcon />}
                onClick={handleChangeStatus}
                size="large"
              >
                Cambiar Status
              </Button>
              <Button
                variant="contained"
                startIcon={<EditIcon />}
                onClick={handleEdit}
                size="large"
                sx={{ boxShadow: 2 }}
              >
                Editar
              </Button>
            </>
          )}
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <div id="event-detail-content">
        <Grid container spacing={3}>
        {/* Event Details Card */}
        <Grid xs={12} md={6}>
          <Card sx={{ height: '100%' }}>
            <CardContent sx={{ p: 3 }}>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                Detalles del Evento
              </Typography>
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
          <Card sx={{ height: '100%' }}>
            <CardContent sx={{ p: 3 }}>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  Coordinadores
                </Typography>
                {isAdmin && (
                  <Button
                    size="small"
                    startIcon={<AddIcon />}
                    variant="outlined"
                    sx={{ borderRadius: 2 }}
                    onClick={handleAssignCoordinator}
                  >
                    Asignar
                  </Button>
                )}
              </Box>
              {coordinators.length === 0 ? (
                <Typography variant="body2" color="text.secondary">
                  No hay coordinadores asignados
                </Typography>
              ) : (
                <List dense>
                  {coordinators.map((coordinator) => (
                    <ListItem
                      key={coordinator.id}
                      secondaryAction={
                        isAdmin && (
                          <IconButton
                            edge="end"
                            size="small"
                            onClick={() =>
                              handleRemoveCoordinator(
                                coordinator.user_id,
                                coordinator.user_name
                              )
                            }
                          >
                            <DeleteIcon fontSize="small" />
                          </IconButton>
                        )
                      }
                    >
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
          <Card sx={{ height: '100%' }}>
            <CardContent sx={{ p: 3 }}>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  Vehículos Presentados
                </Typography>
                {isAdmin && (
                  <Button
                    size="small"
                    startIcon={<AddIcon />}
                    variant="outlined"
                    sx={{ borderRadius: 2 }}
                    onClick={handleAddVehicle}
                  >
                    Agregar
                  </Button>
                )}
              </Box>
              {vehicles.length === 0 ? (
                <Typography variant="body2" color="text.secondary">
                  No hay vehículos agregados
                </Typography>
              ) : (
                <List dense>
                  {vehicles.map((vehicle) => (
                    <ListItem
                      key={vehicle.id}
                      secondaryAction={
                        isAdmin && (
                          <Box>
                            <IconButton
                              edge="end"
                              size="small"
                              onClick={() => handleEditVehicleQuantity(vehicle)}
                              sx={{ mr: 0.5 }}
                            >
                              <EditIcon fontSize="small" />
                            </IconButton>
                            <IconButton
                              edge="end"
                              size="small"
                              onClick={() =>
                                handleRemoveVehicle(vehicle.vehicle_id, vehicle.model_name)
                              }
                            >
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          </Box>
                        )
                      }
                    >
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
            <Card sx={{ height: '100%' }}>
              <CardContent sx={{ p: 3 }}>
                <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>
                    Reporte del Evento
                  </Typography>
                  {isAdmin && (
                    <Button
                      size="small"
                      startIcon={<EditIcon />}
                      variant="outlined"
                      sx={{ borderRadius: 2 }}
                      onClick={handleEditReport}
                    >
                      Editar Reporte
                    </Button>
                  )}
                </Box>
                {!report ? (
                  <Typography variant="body2" color="text.secondary">
                    No hay reporte disponible
                  </Typography>
                ) : (
                  <Grid container spacing={2}>
                    {report.hostess_count != null && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Edecanes
                        </Typography>
                        <Typography variant="body1">{report.hostess_count}</Typography>
                      </Grid>
                    )}
                    {report.setup_vendor && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Montaje
                        </Typography>
                        <Typography variant="body1">{report.setup_vendor}</Typography>
                      </Grid>
                    )}
                    {report.has_promotional != null && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Promocionales
                        </Typography>
                        <Typography variant="body1">{report.has_promotional ? 'Sí' : 'No'}</Typography>
                      </Grid>
                    )}
                    {report.attendees != null && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Asistencia
                        </Typography>
                        <Typography variant="body1">{report.attendees?.toLocaleString()}</Typography>
                      </Grid>
                    )}
                    {report.activities_count != null && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Dinámicas
                        </Typography>
                        <Typography variant="body1">{report.activities_count}</Typography>
                      </Grid>
                    )}
                    {report.leads_collected != null && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Datos Levantados
                        </Typography>
                        <Typography variant="body1">{report.leads_collected}</Typography>
                      </Grid>
                    )}
                    {report.prospects != null && (
                      <Grid xs={6}>
                        <Typography variant="body2" color="text.secondary">
                          Prospectos
                        </Typography>
                        <Typography variant="body1">{report.prospects}</Typography>
                      </Grid>
                    )}
                    {report.dealer_rating != null && (
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
        {/* Photos Section */}
        <Grid xs={12}>
          <EventPhotosSection
            eventId={id!}
            photos={photos}
            canWrite={isAdmin || isCoordinator}
            accessToken={accessToken!}
            onPhotosChange={fetchEventData}
          />
        </Grid>
        </Grid>
      </div>

      {/* Dialogs */}
      <AssignCoordinatorDialog
        open={assignCoordinatorOpen}
        onClose={() => setAssignCoordinatorOpen(false)}
        eventId={id!}
        onSuccess={() => {
          fetchEventData();
          showSnackbar('Coordinador asignado exitosamente', 'success');
        }}
      />

      <AddVehicleDialog
        open={addVehicleOpen}
        onClose={() => setAddVehicleOpen(false)}
        eventId={id!}
        eventBrandId={event.brand_id}
        onSuccess={() => {
          fetchEventData();
          showSnackbar('Vehículo agregado exitosamente', 'success');
        }}
      />

      {selectedVehicle && (
        <EditVehicleQuantityDialog
          open={editQuantityOpen}
          onClose={() => {
            setEditQuantityOpen(false);
            setSelectedVehicle(null);
          }}
          eventId={id!}
          vehicleId={selectedVehicle.vehicle_id}
          currentQuantity={selectedVehicle.quantity}
          vehicleName={selectedVehicle.model_name}
          onSuccess={() => {
            fetchEventData();
            showSnackbar('Cantidad actualizada exitosamente', 'success');
            setSelectedVehicle(null);
          }}
        />
      )}

      <ChangeStatusDialog
        open={changeStatusOpen}
        onClose={() => setChangeStatusOpen(false)}
        event={event}
        onSuccess={() => {
          fetchEventData();
          showSnackbar('Status del evento actualizado exitosamente', 'success');
        }}
      />

      <EventReportDialog
        open={reportDialogOpen}
        onClose={() => setReportDialogOpen(false)}
        eventId={id!}
        existingReport={report}
        onSuccess={() => {
          fetchEventData();
          showSnackbar('Reporte guardado exitosamente', 'success');
        }}
      />

      {confirmDialog && (
        <ConfirmDialog
          open={confirmDialog.open}
          onClose={() => setConfirmDialog(null)}
          onConfirm={confirmDialog.onConfirm}
          title={confirmDialog.title}
          message={confirmDialog.message}
          severity={confirmDialog.severity}
        />
      )}
    </Box>
  );
};
