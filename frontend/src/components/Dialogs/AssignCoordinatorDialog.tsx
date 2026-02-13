import { useState, useEffect } from 'react';
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
  CircularProgress,
  Box,
} from '@mui/material';
import { User } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

interface AssignCoordinatorDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  onSuccess: () => void;
}

export const AssignCoordinatorDialog = ({
  open,
  onClose,
  eventId,
  onSuccess,
}: AssignCoordinatorDialogProps) => {
  const { accessToken } = useAppSelector((state) => state.auth);
  const [coordinators, setCoordinators] = useState<User[]>([]);
  const [selectedUserId, setSelectedUserId] = useState('');
  const [loading, setLoading] = useState(false);
  const [fetchingCoordinators, setFetchingCoordinators] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (open && accessToken) {
      fetchCoordinators();
    }
  }, [open, accessToken]);

  const fetchCoordinators = async () => {
    if (!accessToken) return;

    setFetchingCoordinators(true);
    setError('');

    const response = await apiService.get<User[]>(
      API_ENDPOINTS.USERS.LIST,
      accessToken
    );

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      // Filter only COORDINATOR role users
      const coordinatorUsers = response.data.filter(
        (user) => user.role === 'COORDINATOR'
      );
      setCoordinators(coordinatorUsers);
    }

    setFetchingCoordinators(false);
  };

  const handleAssign = async () => {
    if (!selectedUserId || !accessToken) return;

    setLoading(true);
    setError('');

    const response = await apiService.post(
      API_ENDPOINTS.EVENTS.COORDINATORS(eventId),
      {
        event_id: eventId,
        user_id: selectedUserId,
      },
      accessToken
    );

    if (response.error) {
      setError(response.error);
      setLoading(false);
    } else {
      setLoading(false);
      setSelectedUserId('');
      onSuccess();
      onClose();
    }
  };

  const handleClose = () => {
    if (!loading) {
      setSelectedUserId('');
      setError('');
      onClose();
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle sx={{ fontWeight: 600 }}>Asignar Coordinador</DialogTitle>
      <DialogContent>
        {fetchingCoordinators ? (
          <Box display="flex" justifyContent="center" py={3}>
            <CircularProgress />
          </Box>
        ) : (
          <>
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}

            <FormControl fullWidth sx={{ mt: 2 }}>
              <InputLabel>Coordinador</InputLabel>
              <Select
                value={selectedUserId}
                label="Coordinador"
                onChange={(e) => setSelectedUserId(e.target.value)}
                disabled={loading}
              >
                {coordinators.length === 0 ? (
                  <MenuItem disabled>No hay coordinadores disponibles</MenuItem>
                ) : (
                  coordinators.map((coordinator) => (
                    <MenuItem key={coordinator.id} value={coordinator.id}>
                      {coordinator.name} ({coordinator.email})
                    </MenuItem>
                  ))
                )}
              </Select>
            </FormControl>
          </>
        )}
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button onClick={handleClose} variant="outlined" disabled={loading}>
          Cancelar
        </Button>
        <Button
          onClick={handleAssign}
          variant="contained"
          disabled={loading || !selectedUserId || fetchingCoordinators}
        >
          {loading ? 'Asignando...' : 'Asignar'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
