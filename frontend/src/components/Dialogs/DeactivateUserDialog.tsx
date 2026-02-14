import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  DialogContentText,
  Alert,
  Box,
} from '@mui/material';

interface DeactivateUserDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  userName: string;
  userEmail: string;
  isActive: boolean;
}

export const DeactivateUserDialog: React.FC<DeactivateUserDialogProps> = ({
  open,
  onClose,
  onConfirm,
  userName,
  userEmail,
  isActive,
}) => {
  const action = isActive ? 'desactivar' : 'activar';
  const actionPast = isActive ? 'desactivado' : 'activado';

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {isActive ? 'Desactivar Usuario' : 'Activar Usuario'}
      </DialogTitle>
      <DialogContent>
        <Box sx={{ pt: 1 }}>
          <DialogContentText>
            ¿Estás seguro que deseas {action} al usuario <strong>{userName}</strong>{' '}
            ({userEmail})?
          </DialogContentText>

          <Alert severity={isActive ? 'warning' : 'info'} sx={{ mt: 2 }}>
            {isActive ? (
              <>
                El usuario <strong>no podrá iniciar sesión</strong> hasta que sea
                reactivado. Sus datos se mantendrán en el sistema.
              </>
            ) : (
              <>
                El usuario podrá <strong>iniciar sesión nuevamente</strong> con sus
                credenciales existentes.
              </>
            )}
          </Alert>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancelar</Button>
        <Button
          onClick={() => {
            onConfirm();
            onClose();
          }}
          variant="contained"
          color={isActive ? 'error' : 'primary'}
        >
          {isActive ? 'Desactivar' : 'Activar'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
