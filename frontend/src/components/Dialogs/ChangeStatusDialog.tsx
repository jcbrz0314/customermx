import { useState, useMemo } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Typography,
  Box,
} from '@mui/material';
import { EventStatus, EventWithBrand } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

interface ChangeStatusDialogProps {
  open: boolean;
  onClose: () => void;
  event: EventWithBrand;
  onSuccess: () => void;
}

const statusLabels: Record<EventStatus, string> = {
  PLANNED: 'Planeado',
  ACTIVE: 'Activo',
  COMPLETED: 'Completado',
  CLOSED: 'Cerrado',
};

const validTransitions: Record<EventStatus, EventStatus[]> = {
  PLANNED: ['ACTIVE'],
  ACTIVE: ['COMPLETED'],
  COMPLETED: ['CLOSED'],
  CLOSED: [],
};

export const ChangeStatusDialog = ({
  open,
  onClose,
  event,
  onSuccess,
}: ChangeStatusDialogProps) => {
  const { accessToken } = useAppSelector((state) => state.auth);
  const [newStatus, setNewStatus] = useState<EventStatus | ''>('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const availableStatuses = useMemo(() => {
    return validTransitions[event.status] || [];
  }, [event.status]);

  const handleChange = async () => {
    if (!newStatus || !accessToken) return;

    setLoading(true);
    setError('');

    const response = await apiService.patch(
      API_ENDPOINTS.EVENTS.CHANGE_STATUS(event.id),
      { status: newStatus },
      accessToken
    );

    if (response.error) {
      setError(response.error);
      setLoading(false);
    } else {
      setLoading(false);
      setNewStatus('');
      onSuccess();
      onClose();
    }
  };

  const handleClose = () => {
    if (!loading) {
      setNewStatus('');
      setError('');
      onClose();
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle sx={{ fontWeight: 600 }}>Cambiar Status del Evento</DialogTitle>
      <DialogContent>
        <Alert severity="info" sx={{ mb: 3 }}>
          <strong>Status actual:</strong> {statusLabels[event.status]}
        </Alert>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        {availableStatuses.length === 0 ? (
          <Alert severity="warning">
            El evento está en status final. No se pueden hacer más cambios de status.
          </Alert>
        ) : (
          <>
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Nuevo Status</InputLabel>
              <Select
                value={newStatus}
                label="Nuevo Status"
                onChange={(e) => setNewStatus(e.target.value as EventStatus)}
                disabled={loading}
              >
                {availableStatuses.map((status) => (
                  <MenuItem key={status} value={status}>
                    {statusLabels[status]}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <Box sx={{ p: 2, bgcolor: 'background.default', borderRadius: 2 }}>
              <Typography variant="body2" color="text.secondary">
                <strong>Nota:</strong> Los cambios de status son unidireccionales y no se
                pueden revertir. Asegúrese de que el evento esté listo para la transición.
              </Typography>
            </Box>
          </>
        )}
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button onClick={handleClose} variant="outlined" disabled={loading}>
          Cancelar
        </Button>
        <Button
          onClick={handleChange}
          variant="contained"
          disabled={loading || !newStatus || availableStatuses.length === 0}
          color="primary"
        >
          {loading ? 'Cambiando...' : 'Confirmar Cambio'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
