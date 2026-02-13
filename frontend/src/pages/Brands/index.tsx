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
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { format } from 'date-fns';
import { Brand } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

export const Brands = () => {
  const [brands, setBrands] = useState<Brand[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { accessToken } = useAppSelector((state) => state.auth);

  const fetchBrands = async () => {
    setLoading(true);
    setError('');

    const response = await apiService.get<Brand[]>(
      API_ENDPOINTS.BRANDS.LIST,
      accessToken || undefined
    );

    if (response.error) {
      setError(response.error);
    } else if (response.data) {
      setBrands(response.data);
    }

    setLoading(false);
  };

  useEffect(() => {
    fetchBrands();
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
            Marcas
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Gestiona las marcas automotrices
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          size="large"
          sx={{ boxShadow: 2 }}
        >
          Nueva Marca
        </Button>
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
                  <TableCell><strong>Nombre</strong></TableCell>
                  <TableCell><strong>Fecha de Creación</strong></TableCell>
                  <TableCell align="right"><strong>Acciones</strong></TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {brands.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={3} align="center">
                      No hay marcas registradas
                    </TableCell>
                  </TableRow>
                ) : (
                  brands.map((brand) => (
                    <TableRow key={brand.id} hover>
                      <TableCell>{brand.name}</TableCell>
                      <TableCell>
                        {format(new Date(brand.created_at), 'dd/MM/yyyy HH:mm')}
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
