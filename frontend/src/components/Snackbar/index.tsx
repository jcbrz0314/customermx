import { Snackbar as MuiSnackbar, Alert, IconButton } from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';
import { useSnackbar } from '../../hooks/useSnackbar';

export const SnackbarContainer = () => {
  const { messages, hideSnackbar } = useSnackbar();

  // Show only the most recent message
  const currentMessage = messages[messages.length - 1];

  if (!currentMessage) {
    return null;
  }

  return (
    <MuiSnackbar
      open={true}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      onClose={() => hideSnackbar(currentMessage.id)}
    >
      <Alert
        severity={currentMessage.severity}
        variant="filled"
        action={
          <IconButton
            size="small"
            color="inherit"
            onClick={() => hideSnackbar(currentMessage.id)}
          >
            <CloseIcon fontSize="small" />
          </IconButton>
        }
        sx={{
          minWidth: 300,
          boxShadow: 4,
        }}
      >
        {currentMessage.message}
      </Alert>
    </MuiSnackbar>
  );
};
