import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  CardContent,
  Container,
  TextField,
  Typography,
  Alert,
  CircularProgress,
  Chip,
} from '@mui/material';
import {
  DirectionsCar as CarIcon,
  CheckCircle as CheckCircleIcon,
  ErrorOutline as ErrorOutlineIcon,
} from '@mui/icons-material';
import { apiService } from '../../services/api/apiService';
import { API_ENDPOINTS } from '../../services/api/apiConstants';
import { UserRole } from '../../types';

interface InvitationInfo {
  email: string;
  role: UserRole;
  brand_id?: string;
  expires_at: string;
}

const ROLE_LABELS: Record<UserRole, string> = {
  ADMIN: 'Administrador',
  COORDINATOR: 'Coordinador',
  BRAND: 'Marca',
};

export const AcceptInvitation = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const token = searchParams.get('token') || '';

  const [invitationInfo, setInvitationInfo] = useState<InvitationInfo | null>(null);
  const [loadingToken, setLoadingToken] = useState(true);
  const [tokenError, setTokenError] = useState('');

  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState('');
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    if (!token) {
      setTokenError('El enlace de invitación no es válido. No se encontró el token.');
      setLoadingToken(false);
      return;
    }
    validateToken();
  }, [token]);

  const validateToken = async () => {
    try {
      setLoadingToken(true);
      const response = await apiService.get(`${API_ENDPOINTS.INVITATIONS.VALIDATE}?token=${token}`);
      if (response.error) {
        if (response.status === 410) {
          setTokenError('Esta invitación ha expirado. Contacta al administrador para que te reenvíe una nueva.');
        } else if (response.status === 409) {
          setTokenError('Esta invitación ya fue aceptada. Si ya tienes cuenta, inicia sesión.');
        } else if (response.status === 404) {
          setTokenError('El enlace de invitación no es válido.');
        } else {
          setTokenError(response.error);
        }
      } else {
        setInvitationInfo(response.data);
      }
    } catch {
      setTokenError('El enlace de invitación no es válido o ha expirado.');
    } finally {
      setLoadingToken(false);
    }
  };

  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};
    if (!name.trim()) errors.name = 'El nombre es requerido';
    if (!password) errors.password = 'La contraseña es requerida';
    else if (password.length < 6) errors.password = 'La contraseña debe tener al menos 6 caracteres';
    if (!confirmPassword) errors.confirmPassword = 'Confirma tu contraseña';
    else if (password !== confirmPassword) errors.confirmPassword = 'Las contraseñas no coinciden';
    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    try {
      setSubmitting(true);
      setSubmitError('');
      const response = await apiService.post(API_ENDPOINTS.INVITATIONS.ACCEPT, { token, name, password });
      if (response.error) {
        setSubmitError(response.error);
      } else {
        setSuccess(true);
      }
    } catch {
      setSubmitError('Ocurrió un error al crear tu cuenta. Intenta nuevamente.');
    } finally {
      setSubmitting(false);
    }
  };

  const containerStyle = {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    background: 'linear-gradient(135deg, #1a237e 0%, #283593 50%, #3949ab 100%)',
    py: 4,
  };

  if (loadingToken) {
    return (
      <Box sx={containerStyle}>
        <Container maxWidth="sm">
          <Card elevation={8} sx={{ borderRadius: 3, textAlign: 'center', p: 4 }}>
            <CircularProgress size={48} />
            <Typography variant="h6" sx={{ mt: 2 }}>
              Validando invitación...
            </Typography>
          </Card>
        </Container>
      </Box>
    );
  }

  if (tokenError) {
    return (
      <Box sx={containerStyle}>
        <Container maxWidth="sm">
          <Card elevation={8} sx={{ borderRadius: 3 }}>
            <CardContent sx={{ p: 4, textAlign: 'center' }}>
              <ErrorOutlineIcon sx={{ fontSize: 64, color: 'error.main', mb: 2 }} />
              <Typography variant="h5" gutterBottom fontWeight={600}>
                Invitación no válida
              </Typography>
              <Alert severity="error" sx={{ mb: 3, textAlign: 'left' }}>
                {tokenError}
              </Alert>
              <Button variant="contained" onClick={() => navigate('/login')}>
                Ir al inicio de sesión
              </Button>
            </CardContent>
          </Card>
        </Container>
      </Box>
    );
  }

  if (success) {
    return (
      <Box sx={containerStyle}>
        <Container maxWidth="sm">
          <Card elevation={8} sx={{ borderRadius: 3 }}>
            <CardContent sx={{ p: 4, textAlign: 'center' }}>
              <CheckCircleIcon sx={{ fontSize: 64, color: 'success.main', mb: 2 }} />
              <Typography variant="h5" gutterBottom fontWeight={600}>
                ¡Cuenta creada exitosamente!
              </Typography>
              <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                Tu cuenta ha sido activada. Ya puedes iniciar sesión con tu correo y contraseña.
              </Typography>
              <Button variant="contained" size="large" onClick={() => navigate('/login')}>
                Ir al inicio de sesión
              </Button>
            </CardContent>
          </Card>
        </Container>
      </Box>
    );
  }

  return (
    <Box sx={containerStyle}>
      <Container maxWidth="sm">
        <Box sx={{ textAlign: 'center', mb: 3 }}>
          <CarIcon sx={{ fontSize: 48, color: 'white', mb: 1 }} />
          <Typography variant="h4" sx={{ color: 'white', fontWeight: 700 }}>
            CustomerMX
          </Typography>
        </Box>

        <Card elevation={8} sx={{ borderRadius: 3 }}>
          <CardContent sx={{ p: 4 }}>
            <Typography variant="h5" gutterBottom fontWeight={600} textAlign="center">
              Aceptar Invitación
            </Typography>
            <Typography variant="body2" color="text.secondary" textAlign="center" sx={{ mb: 3 }}>
              Completa tu perfil para activar tu cuenta
            </Typography>

            {/* Invitation info */}
            <Box sx={{ mb: 3, p: 2, bgcolor: 'grey.50', borderRadius: 2 }}>
              <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                Correo de la invitación
              </Typography>
              <Typography variant="body1" fontWeight={500}>
                {invitationInfo?.email}
              </Typography>
              <Box sx={{ mt: 1 }}>
                <Chip
                  label={ROLE_LABELS[invitationInfo?.role as UserRole] || invitationInfo?.role}
                  size="small"
                  color="primary"
                  variant="outlined"
                />
              </Box>
            </Box>

            {submitError && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {submitError}
              </Alert>
            )}

            <form onSubmit={handleSubmit}>
              <TextField
                fullWidth
                label="Nombre completo"
                value={name}
                onChange={(e) => setName(e.target.value)}
                error={!!formErrors.name}
                helperText={formErrors.name}
                sx={{ mb: 2 }}
                autoFocus
              />
              <TextField
                fullWidth
                label="Contraseña"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                error={!!formErrors.password}
                helperText={formErrors.password || 'Mínimo 6 caracteres'}
                sx={{ mb: 2 }}
              />
              <TextField
                fullWidth
                label="Confirmar contraseña"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                error={!!formErrors.confirmPassword}
                helperText={formErrors.confirmPassword}
                sx={{ mb: 3 }}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                size="large"
                disabled={submitting}
                sx={{ py: 1.5 }}
              >
                {submitting ? <CircularProgress size={24} color="inherit" /> : 'Activar cuenta'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
};
