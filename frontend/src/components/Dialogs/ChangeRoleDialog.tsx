import React, { useState, useEffect } from 'react';
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
  Box,
} from '@mui/material';
import { UserRole } from '../../types';

interface ChangeRoleDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (newRole: UserRole) => void;
  currentRole: UserRole;
  userName: string;
}

export const ChangeRoleDialog: React.FC<ChangeRoleDialogProps> = ({
  open,
  onClose,
  onConfirm,
  currentRole,
  userName,
}) => {
  const [selectedRole, setSelectedRole] = useState<UserRole>(currentRole);

  useEffect(() => {
    if (open) {
      setSelectedRole(currentRole);
    }
  }, [open, currentRole]);

  const handleConfirm = () => {
    if (selectedRole !== currentRole) {
      onConfirm(selectedRole);
    }
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

  const isRoleChanged = selectedRole !== currentRole;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Cambiar Rol de Usuario</DialogTitle>
      <DialogContent>
        <Box sx={{ pt: 2 }}>
          <Alert severity="info" sx={{ mb: 3 }}>
            Estás cambiando el rol de <strong>{userName}</strong>
          </Alert>

          <FormControl fullWidth>
            <InputLabel>Rol</InputLabel>
            <Select
              value={selectedRole}
              onChange={(e) => setSelectedRole(e.target.value as UserRole)}
              label="Rol"
            >
              {roles.map((role) => (
                <MenuItem key={role.value} value={role.value}>
                  <Box>
                    <Box sx={{ fontWeight: 500 }}>{role.label}</Box>
                    <Box sx={{ fontSize: '0.875rem', color: 'text.secondary' }}>
                      {role.description}
                    </Box>
                  </Box>
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {isRoleChanged && (
            <Alert severity="warning" sx={{ mt: 2 }}>
              El cambio de rol afectará los permisos y accesos del usuario inmediatamente.
            </Alert>
          )}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancelar</Button>
        <Button
          onClick={handleConfirm}
          variant="contained"
          disabled={!isRoleChanged}
        >
          Cambiar Rol
        </Button>
      </DialogActions>
    </Dialog>
  );
};
