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
  Grid,
  Alert,
  Box,
} from '@mui/material';
import { UserRole } from '../../types';
import { useAppSelector } from '../../hooks/useRedux';

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

  const brands = useAppSelector((state) => state.auth.brands) || [];

  useEffect(() => {
    if (open) {
      setEmail('');
      setRole('COORDINATOR');
      setBrandId('');
      setEmailError('');
    }
  }, [open]);

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

          <Grid container spacing={2}>
            <Grid item xs={12}>
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
            </Grid>

            <Grid item xs={12}>
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
            </Grid>

            {role === 'BRAND' && (
              <Grid item xs={12}>
                <FormControl fullWidth required>
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
              </Grid>
            )}
          </Grid>
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
