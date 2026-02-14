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

interface ResendInvitationDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  invitationEmail: string;
}

export const ResendInvitationDialog: React.FC<ResendInvitationDialogProps> = ({
  open,
  onClose,
  onConfirm,
  invitationEmail,
}) => {
  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Reenviar Invitación</DialogTitle>
      <DialogContent>
        <Box sx={{ pt: 1 }}>
          <DialogContentText>
            ¿Deseas reenviar la invitación a <strong>{invitationEmail}</strong>?
          </DialogContentText>

          <Alert severity="info" sx={{ mt: 2 }}>
            Se generará un nuevo código de invitación y se enviará un email al
            destinatario. El código anterior quedará invalidado.
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
        >
          Reenviar
        </Button>
      </DialogActions>
    </Dialog>
  );
};
