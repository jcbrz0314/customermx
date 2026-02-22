import { useEffect, useState } from 'react';
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
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { useAppSelector } from '../../hooks/useRedux';
import { format } from 'date-fns';
import { VehicleWithBrand } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';

export const Vehicles = () => {
  const [vehicles, setVehicles] = useState<VehicleWithBrand[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { accessToken, user } = useAppSelector((state) => state.auth);

  const fetchVehicles = async () => {
    setLoading(true);
    setError('');

    const response = await apiService.get<VehicleWithBrand[]>(
      API_ENDPOINTS.VEHICLES.LIST,
      accessToken || undefined
    );

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      setVehicles(response.data);
    }

    setLoading(false);
  };

  useEffect(() => {
    fetchVehicles();
  }, []);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={4}>
        <Box>
          <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
            Vehículos
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Catálogo de modelos de vehículos
          </Typography>
        </Box>
        {user?.role === 'ADMIN' && (
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            size="large"
            sx={{ boxShadow: 2 }}
          >
            Nuevo Vehículo
          </Button>
        )}
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <Card>
        <CardContent sx={{ p: 0 }}>
          <TableContainer component={Paper} elevation={0}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell><strong>Modelo</strong></TableCell>
                  <TableCell><strong>Marca</strong></TableCell>
                  <TableCell><strong>Fecha de Creación</strong></TableCell>
                  <TableCell align="right"><strong>Acciones</strong></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {vehicles.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} align="center">
                      No hay vehículos registrados
                    </TableCell>
                  </TableRow>
                ) : (
                  vehicles.map((vehicle) => (
                    <TableRow key={vehicle.id} hover>
                      <TableCell>{vehicle.model_name}</TableCell>
                      <TableCell>
                        <Chip label={vehicle.brand_name} color="primary" size="small" />
                      </TableCell>
                      <TableCell>
                        {format(new Date(vehicle.created_at), 'dd/MM/yyyy HH:mm')}
                      </TableCell>
                      <TableCell align="right">
                        <Button size="small">Ver Detalles</Button>
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
