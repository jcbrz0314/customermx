import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Box,
  Stack,
} from '@mui/material';
import { Brand, UserRole } from '../../types';
import { useAppSelector } from '../../hooks/useRedux';
import { apiService, API_ENDPOINTS } from '../../services/api';

interface InviteUserDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (email: string, role: UserRole, brandId: string | null) => void;
}

export const InviteUserDialog: React.FC<InviteUserDialogProps> = ({
  open,
  onClose,
  onConfirm,
}) => {
  const [email, setEmail] = useState('');
  const [role, setRole] = useState<UserRole>('COORDINATOR');
  const [brandId, setBrandId] = useState<string>('');
  const [emailError, setEmailError] = useState('');
  const [brands, setBrands] = useState<Brand[]>([]);
  const [loadingBrands, setLoadingBrands] = useState(false);

  const { accessToken } = useAppSelector((state) => state.auth);

  useEffect(() => {
    if (open) {
      setEmail('');
      setRole('COORDINATOR');
      setBrandId('');
      setEmailError('');
      fetchBrands();
    }
  }, [open]);

  const fetchBrands = async () => {
    if (!accessToken) return;
    setLoadingBrands(true);
    const response = await apiService.get<Brand[]>(API_ENDPOINTS.BRANDS.LIST, accessToken);
    if (response.data) {
      setBrands(response.data);
    }
    setLoadingBrands(false);
  };

  const validateEmail = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const handleEmailChange = (value: string) => {
    setEmail(value);
    if (value && !validateEmail(value)) {
      setEmailError('Email inválido');
    } else {
      setEmailError('');
    }
  };

  const handleConfirm = () => {
    if (!validateEmail(email)) {
      setEmailError('Email inválido');
      return;
    }

    if (role === 'BRAND' && !brandId) {
      return;
    }

    onConfirm(email, role, role === 'BRAND' ? brandId : null);
    onClose();
  };

  const roles: { value: UserRole; label: string; description: string }[] = [
    {
      value: 'ADMIN',
      label: 'Administrador',
      description: 'Acceso completo al sistema',
    },
    {
      value: 'COORDINATOR',
      label: 'Coordinador',
      description: 'Gestión de eventos asignados',
    },
    {
      value: 'BRAND',
      label: 'Marca',
      description: 'Consulta de eventos de su marca',
    },
    {
      value: 'VISUALIZER',
      label: 'Visualizador',
      description: 'Acceso de lectura a todo el sistema, sin poder hacer cambios',
    },
  ];

  const isFormValid =
    email &&
    !emailError &&
    role &&
    (role !== 'BRAND' || brandId);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Invitar Nuevo Usuario</DialogTitle>
      <DialogContent>
        <Box sx={{ pt: 2 }}>
          <Alert severity="info" sx={{ mb: 3 }}>
            Se enviará un email de invitación al usuario para que complete su
            registro.
          </Alert>

          <Stack spacing={2}>
            <TextField
              fullWidth
              label="Email"
              type="email"
              value={email}
              onChange={(e) => handleEmailChange(e.target.value)}
              error={!!emailError}
              helperText={emailError}
              required
            />

            <FormControl fullWidth required>
              <InputLabel>Rol</InputLabel>
              <Select
                value={role}
                onChange={(e) => {
                  setRole(e.target.value as UserRole);
                  if (e.target.value !== 'BRAND') {
                    setBrandId('');
                  }
                }}
                label="Rol"
              >
                {roles.map((roleOption) => (
                  <MenuItem key={roleOption.value} value={roleOption.value}>
                    <Box>
                      <Box sx={{ fontWeight: 500 }}>{roleOption.label}</Box>
                      <Box sx={{ fontSize: '0.875rem', color: 'text.secondary' }}>
                        {roleOption.description}
                      </Box>
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            {role === 'BRAND' && (
              <FormControl fullWidth required disabled={loadingBrands}>
                <InputLabel>Marca</InputLabel>
                <Select
                  value={brandId}
                  onChange={(e) => setBrandId(e.target.value)}
                  label="Marca"
                >
                  {brands.map((brand) => (
                    <MenuItem key={brand.id} value={brand.id}>
                      {brand.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}
          </Stack>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancelar</Button>
        <Button
          onClick={handleConfirm}
          variant="contained"
          disabled={!isFormValid}
        >
          Enviar Invitación
        </Button>
      </DialogActions>
    </Dialog>
  );
};
