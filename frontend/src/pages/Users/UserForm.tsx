import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Alert,
  FormControlLabel,
  Switch,
} from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { useAppSelector } from '../../hooks/useRedux';
import { useSnackbar } from '../../hooks/useSnackbar';
import { apiService } from '../../services/api/apiService';
import { API_ENDPOINTS } from '../../services/api/apiConstants';
import { User, UserRole } from '../../types';

export const UserForm: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { showSnackbar } = useSnackbar();
  const brands = useAppSelector((state) => state.auth.brands) || [];
  const currentUser = useAppSelector((state) => state.auth.user);

  const isEditMode = Boolean(id);

  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    role: 'COORDINATOR' as UserRole,
    brand_id: '',
    is_active: true,
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (isEditMode && id) {
      fetchUser(id);
    }
  }, [id, isEditMode]);

  const fetchUser = async (userId: string) => {
    try {
      setLoading(true);
      const response = await apiService.get(API_ENDPOINTS.USERS.BY_ID(userId));
      const user: User = response.data;
      setFormData({
        name: user.name,
        email: user.email,
        password: '',
        role: user.role,
        brand_id: user.brand_id || '',
        is_active: user.is_active,
      });
    } catch (err: any) {
      showSnackbar(err.message || 'Error al cargar usuario', 'error');
      navigate('/users');
    } finally {
      setLoading(false);
    }
  };

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'El nombre es requerido';
    }

    if (!formData.email.trim()) {
      newErrors.email = 'El email es requerido';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Email inválido';
    }

    if (!isEditMode && !formData.password) {
      newErrors.password = 'La contraseña es requerida';
    }

    if (!isEditMode && formData.password && formData.password.length < 6) {
      newErrors.password = 'La contraseña debe tener al menos 6 caracteres';
    }

    if (formData.role === 'BRAND' && !formData.brand_id) {
      newErrors.brand_id = 'La marca es requerida para usuarios de tipo Marca';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      setLoading(true);

      const payload: any = {
        name: formData.name,
        email: formData.email,
        role: formData.role,
        brand_id: formData.role === 'BRAND' ? formData.brand_id : null,
        is_active: formData.is_active,
      };

      if (formData.password) {
        payload.password = formData.password;
      }

      if (isEditMode && id) {
        await apiService.put(API_ENDPOINTS.USERS.UPDATE(id), payload);
        showSnackbar('Usuario actualizado exitosamente', 'success');
      } else {
        await apiService.post(API_ENDPOINTS.USERS.CREATE, payload);
        showSnackbar('Usuario creado exitosamente', 'success');
      }

      navigate('/users');
    } catch (err: any) {
      showSnackbar(
        err.message || `Error al ${isEditMode ? 'actualizar' : 'crear'} usuario`,
        'error'
      );
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    navigate('/users');
  };

  if (!currentUser || currentUser.role !== 'ADMIN') {
    return (
      <Box sx={{ p: 3 }}>
        <Alert severity="error">
          No tienes permisos para acceder a esta sección.
        </Alert>
      </Box>
    );
  }

  if (loading && isEditMode) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" sx={{ fontWeight: 600 }}>
          {isEditMode ? 'Editar Usuario' : 'Nuevo Usuario'}
        </Typography>
      </Box>

      <Card>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Nombre Completo"
                  value={formData.name}
                  onChange={(e) =>
                    setFormData({ ...formData, name: e.target.value })
                  }
                  error={!!errors.name}
                  helperText={errors.name}
                  required
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Email"
                  type="email"
                  value={formData.email}
                  onChange={(e) =>
                    setFormData({ ...formData, email: e.target.value })
                  }
                  error={!!errors.email}
                  helperText={errors.email}
                  required
                  disabled={isEditMode}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label={isEditMode ? 'Nueva Contraseña (opcional)' : 'Contraseña'}
                  type="password"
                  value={formData.password}
                  onChange={(e) =>
                    setFormData({ ...formData, password: e.target.value })
                  }
                  error={!!errors.password}
                  helperText={
                    errors.password ||
                    (isEditMode ? 'Dejar en blanco para mantener la actual' : '')
                  }
                  required={!isEditMode}
                />
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth required error={!!errors.role}>
                  <InputLabel>Rol</InputLabel>
                  <Select
                    value={formData.role}
                    onChange={(e) => {
                      setFormData({
                        ...formData,
                        role: e.target.value as UserRole,
                        brand_id: e.target.value !== 'BRAND' ? '' : formData.brand_id,
                      });
                    }}
                    label="Rol"
                  >
                    <MenuItem value="ADMIN">
                      <Box>
                        <Typography variant="body1">Administrador</Typography>
                        <Typography variant="caption" color="text.secondary">
                          Acceso completo al sistema
                        </Typography>
                      </Box>
                    </MenuItem>
                    <MenuItem value="COORDINATOR">
                      <Box>
                        <Typography variant="body1">Coordinador</Typography>
                        <Typography variant="caption" color="text.secondary">
                          Gestión de eventos asignados
                        </Typography>
                      </Box>
                    </MenuItem>
                    <MenuItem value="BRAND">
                      <Box>
                        <Typography variant="body1">Marca</Typography>
                        <Typography variant="caption" color="text.secondary">
                          Consulta de eventos de su marca
                        </Typography>
                      </Box>
                    </MenuItem>
                  </Select>
                </FormControl>
              </Grid>

              {formData.role === 'BRAND' && (
                <Grid item xs={12} md={6}>
                  <FormControl fullWidth required error={!!errors.brand_id}>
                    <InputLabel>Marca</InputLabel>
                    <Select
                      value={formData.brand_id}
                      onChange={(e) =>
                        setFormData({ ...formData, brand_id: e.target.value })
                      }
                      label="Marca"
                    >
                      {brands.map((brand) => (
                        <MenuItem key={brand.id} value={brand.id}>
                          {brand.name}
                        </MenuItem>
                      ))}
                    </Select>
                    {errors.brand_id && (
                      <Typography variant="caption" color="error" sx={{ mt: 0.5 }}>
                        {errors.brand_id}
                      </Typography>
                    )}
                  </FormControl>
                </Grid>
              )}

              <Grid item xs={12}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={formData.is_active}
                      onChange={(e) =>
                        setFormData({ ...formData, is_active: e.target.checked })
                      }
                    />
                  }
                  label="Usuario Activo"
                />
              </Grid>

              <Grid item xs={12}>
                <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
                  <Button onClick={handleCancel} disabled={loading}>
                    Cancelar
                  </Button>
                  <Button
                    type="submit"
                    variant="contained"
                    disabled={loading}
                  >
                    {loading ? (
                      <CircularProgress size={24} />
                    ) : isEditMode ? (
                      'Actualizar'
                    ) : (
                      'Crear Usuario'
                    )}
                  </Button>
                </Box>
              </Grid>
            </Grid>
          </form>
        </CardContent>
      </Card>
    </Box>
  );
};
