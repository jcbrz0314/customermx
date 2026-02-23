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
  Stack,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { Brand, CreateVehicleRequest, UpdateVehicleRequest, VehicleWithBrand } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

export const VehicleForm = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { accessToken } = useAppSelector((state) => state.auth);
  const isEditMode = Boolean(id);

  const [brands, setBrands] = useState<Brand[]>([]);
  const [loading, setLoading] = useState(false);
  const [fetchingVehicle, setFetchingVehicle] = useState(isEditMode);
  const [error, setError] = useState('');
  const [brandId, setBrandId] = useState('');
  const [modelName, setModelName] = useState('');

  const fetchBrands = async () => {
    if (!accessToken) return;
    const response = await apiService.get<Brand[]>(API_ENDPOINTS.BRANDS.LIST, accessToken);
    if (response.data) setBrands(response.data);
  };

  const fetchVehicle = async () => {
    if (!id || !accessToken) return;
    setFetchingVehicle(true);
    const response = await apiService.get<VehicleWithBrand>(
      API_ENDPOINTS.VEHICLES.BY_ID(id),
      accessToken
    );
    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      setBrandId(response.data.brand_id);
      setModelName(response.data.model_name);
    }
    setFetchingVehicle(false);
  };

  useEffect(() => {
    if (accessToken) {
      fetchBrands();
      if (isEditMode) fetchVehicle();
    }
  }, [accessToken, isEditMode]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (isEditMode && id) {
      const updateData: UpdateVehicleRequest = { model_name: modelName };
      const response = await apiService.put(
        API_ENDPOINTS.VEHICLES.UPDATE(id),
        updateData,
        accessToken || undefined
      );
      if (response.error) {
        setError(response.error);
      } else {
        navigate('/vehicles');
      }
    } else {
      const createData: CreateVehicleRequest = { brand_id: brandId, model_name: modelName };
      const response = await apiService.post(
        API_ENDPOINTS.VEHICLES.CREATE,
        createData,
        accessToken || undefined
      );
      if (response.error) {
        setError(response.error);
      } else {
        navigate('/vehicles');
      }
    }

    setLoading(false);
  };

  if (fetchingVehicle) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box maxWidth="600px" mx="auto">
      <Box mb={4}>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
          {isEditMode ? 'Editar Vehículo' : 'Nuevo Vehículo'}
        </Typography>
        <Typography variant="body1" color="text.secondary">
          {isEditMode ? 'Modifica los datos del vehículo' : 'Registra un nuevo modelo de vehículo'}
        </Typography>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <form onSubmit={handleSubmit}>
        <Stack spacing={3}>
          <Card>
            <CardContent sx={{ p: 4 }}>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
                Datos del Vehículo
              </Typography>
              <Stack spacing={3}>
                <FormControl fullWidth required disabled={isEditMode}>
                  <InputLabel>Marca</InputLabel>
                  <Select
                    value={brandId}
                    label="Marca"
                    onChange={(e) => setBrandId(e.target.value)}
                  >
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
                  label="Nombre del Modelo"
                  value={modelName}
                  onChange={(e) => setModelName(e.target.value)}
                  placeholder="Ej: Silverado 2500HD, Equinox, Trax"
                />
              </Stack>
            </CardContent>
          </Card>

          <Box display="flex" gap={2} justifyContent="flex-end" pt={2}>
            <Button
              variant="outlined"
              size="large"
              onClick={() => navigate('/vehicles')}
              disabled={loading}
              sx={{ minWidth: 140, py: 1.5 }}
            >
              Cancelar
            </Button>
            <Button
              type="submit"
              variant="contained"
              size="large"
              disabled={loading}
              sx={{ minWidth: 140, py: 1.5, boxShadow: 2 }}
            >
              {loading ? 'Guardando...' : 'Guardar'}
            </Button>
          </Box>
        </Stack>
      </form>
    </Box>
  );
};
