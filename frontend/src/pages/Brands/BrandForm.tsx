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
} from '@mui/material';
import { Save as SaveIcon, Cancel as CancelIcon } from '@mui/icons-material';
import { Brand, CreateBrandRequest, UpdateBrandRequest } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

export const BrandForm = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { accessToken } = useAppSelector((state) => state.auth);
  const isEditMode = Boolean(id);

  const [loading, setLoading] = useState(false);
  const [fetchingBrand, setFetchingBrand] = useState(isEditMode);
  const [error, setError] = useState('');
  const [name, setName] = useState('');

  const fetchBrand = async () => {
    if (!id || !accessToken) return;

    setFetchingBrand(true);
    const response = await apiService.get<Brand>(
      API_ENDPOINTS.BRANDS.BY_ID(id),
      accessToken
    );

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      setName(response.data.name);
    }

    setFetchingBrand(false);
  };

  useEffect(() => {
    if (accessToken && isEditMode) {
      fetchBrand();
    }
  }, [accessToken, isEditMode]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    if (isEditMode && id) {
      const updateData: UpdateBrandRequest = { name };
      const response = await apiService.put(
        API_ENDPOINTS.BRANDS.UPDATE(id),
        updateData,
        accessToken || undefined
      );

      if (response.error) {
        setError(response.error);
      } else {
        navigate('/brands');
      }
    } else {
      const createData: CreateBrandRequest = { name };
      const response = await apiService.post(
        API_ENDPOINTS.BRANDS.CREATE,
        createData,
        accessToken || undefined
      );

      if (response.error) {
        setError(response.error);
      } else {
        navigate('/brands');
      }
    }

    setLoading(false);
  };

  if (fetchingBrand) {
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
          {isEditMode ? 'Editar Marca' : 'Nueva Marca'}
        </Typography>
        <Typography variant="body1" color="text.secondary">
          {isEditMode ? 'Modifica el nombre de la marca' : 'Registra una nueva marca automotriz'}
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
                Información de la Marca
              </Typography>
              <TextField
                fullWidth
                required
                label="Nombre de la Marca"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Ej: Chevrolet, Ford, Toyota"
              />
            </CardContent>
          </Card>

          <Box display="flex" gap={2} justifyContent="flex-end" pt={2}>
            <Button
              variant="outlined"
              size="large"
              onClick={() => navigate('/brands')}
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
