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
  TextField,
  Alert,
  CircularProgress,
  Box,
} from '@mui/material';
import { Vehicle } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

interface AddVehicleDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  eventBrandId: string;
  onSuccess: () => void;
}

export const AddVehicleDialog = ({
  open,
  onClose,
  eventId,
  eventBrandId,
  onSuccess,
}: AddVehicleDialogProps) => {
  const { accessToken } = useAppSelector((state) => state.auth);
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [selectedVehicleId, setSelectedVehicleId] = useState('');
  const [quantity, setQuantity] = useState('1');
  const [loading, setLoading] = useState(false);
  const [fetchingVehicles, setFetchingVehicles] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (open && accessToken) {
      fetchVehicles();
    }
  }, [open, accessToken, eventBrandId]);

  const fetchVehicles = async () => {
    if (!accessToken) return;

    setFetchingVehicles(true);
    setError('');

    // Fetch vehicles for the event's brand
    const response = await apiService.get<Vehicle[]>(
      API_ENDPOINTS.BRANDS.VEHICLES(eventBrandId),
      accessToken
    );

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      setVehicles(response.data);
    }

    setFetchingVehicles(false);
  };

  const handleAdd = async () => {
    if (!selectedVehicleId || !accessToken || !quantity) return;

    const quantityNum = parseInt(quantity);
    if (isNaN(quantityNum) || quantityNum < 1) {
      setError('La cantidad debe ser mayor a 0');
      return;
    }

    setLoading(true);
    setError('');

    const response = await apiService.post(
      API_ENDPOINTS.EVENTS.VEHICLES(eventId),
      {
        event_id: eventId,
        vehicle_id: selectedVehicleId,
        quantity: quantityNum,
      },
      accessToken
    );

    if (response.error) {
      setError(response.error);
      setLoading(false);
    } else {
      setLoading(false);
      setSelectedVehicleId('');
      setQuantity('1');
      onSuccess();
      onClose();
    }
  };

  const handleClose = () => {
    if (!loading) {
      setSelectedVehicleId('');
      setQuantity('1');
      setError('');
      onClose();
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle sx={{ fontWeight: 600 }}>Agregar Vehículo</DialogTitle>
      <DialogContent>
        {fetchingVehicles ? (
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

            <FormControl fullWidth sx={{ mt: 2, mb: 2 }}>
              <InputLabel>Vehículo</InputLabel>
              <Select
                value={selectedVehicleId}
                label="Vehículo"
                onChange={(e) => setSelectedVehicleId(e.target.value)}
                disabled={loading}
              >
                {vehicles.length === 0 ? (
                  <MenuItem disabled>No hay vehículos disponibles para esta marca</MenuItem>
                ) : (
                  vehicles.map((vehicle) => (
                    <MenuItem key={vehicle.id} value={vehicle.id}>
                      {vehicle.model_name}
                    </MenuItem>
                  ))
                )}
              </Select>
            </FormControl>

            <TextField
              fullWidth
              type="number"
              label="Cantidad"
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
              disabled={loading}
              inputProps={{ min: 1 }}
              helperText="Número de unidades que se presentarán"
            />
          </>
        )}
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button onClick={handleClose} variant="outlined" disabled={loading}>
          Cancelar
        </Button>
        <Button
          onClick={handleAdd}
          variant="contained"
          disabled={loading || !selectedVehicleId || fetchingVehicles}
        >
          {loading ? 'Agregando...' : 'Agregar'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
