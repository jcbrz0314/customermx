import { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Alert,
  Typography,
} from '@mui/material';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

interface EditVehicleQuantityDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  vehicleId: string;
  currentQuantity: number;
  vehicleName: string;
  onSuccess: () => void;
}

export const EditVehicleQuantityDialog = ({
  open,
  onClose,
  eventId,
  vehicleId,
  currentQuantity,
  vehicleName,
  onSuccess,
}: EditVehicleQuantityDialogProps) => {
  const { accessToken } = useAppSelector((state) => state.auth);
  const [quantity, setQuantity] = useState(currentQuantity.toString());
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSave = async () => {
    if (!accessToken) return;

    const quantityNum = parseInt(quantity);
    if (isNaN(quantityNum) || quantityNum < 1) {
      setError('La cantidad debe ser mayor a 0');
      return;
    }

    setLoading(true);
    setError('');

    const response = await apiService.patch(
      API_ENDPOINTS.EVENTS.UPDATE_VEHICLE_QUANTITY(eventId, vehicleId),
      { quantity: quantityNum },
      accessToken
    );

    if (response.error) {
      setError(response.error);
      setLoading(false);
    } else {
      setLoading(false);
      onSuccess();
      onClose();
    }
  };

  const handleClose = () => {
    if (!loading) {
      setQuantity(currentQuantity.toString());
      setError('');
      onClose();
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="xs" fullWidth>
      <DialogTitle sx={{ fontWeight: 600 }}>Editar Cantidad</DialogTitle>
      <DialogContent>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          {vehicleName}
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <TextField
          fullWidth
          type="number"
          label="Cantidad"
          value={quantity}
          onChange={(e) => setQuantity(e.target.value)}
          disabled={loading}
          inputProps={{ min: 1 }}
          autoFocus
        />
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button onClick={handleClose} variant="outlined" disabled={loading}>
          Cancelar
        </Button>
        <Button onClick={handleSave} variant="contained" disabled={loading}>
          {loading ? 'Guardando...' : 'Guardar'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
