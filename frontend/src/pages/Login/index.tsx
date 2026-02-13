import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Card,
  CardContent,
  Container,
  TextField,
  Typography,
  Alert,
} from '@mui/material';
import { DirectionsCar as CarIcon } from '@mui/icons-material';
import { useAuth } from '../../hooks/useAuth';

export const Login = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login({ email, password });
      navigate('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al iniciar sesión');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        py: 4,
      }}
    >
      <Container component="main" maxWidth="xs">
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
          }}
        >
          <Box
            sx={{
              mb: 3,
              p: 2.5,
              borderRadius: '20px',
              background: 'linear-gradient(135deg, #2563eb 0%, #3b82f6 100%)',
              boxShadow: '0 10px 40px rgba(37, 99, 235, 0.3)',
            }}
          >
            <CarIcon sx={{ fontSize: 48, color: 'white' }} />
          </Box>
          <Typography
            component="h1"
            variant="h3"
            sx={{ color: 'white', fontWeight: 700, mb: 1 }}
          >
            CustomerMX
          </Typography>
          <Typography variant="body1" sx={{ color: 'rgba(255,255,255,0.9)', mb: 4 }}>
            Sistema de Gestión de Eventos Automotrices
          </Typography>

          <Card
            sx={{
              width: '100%',
              boxShadow: '0 20px 60px rgba(0, 0, 0, 0.3)',
            }}
          >
            <CardContent sx={{ p: 4 }}>
              <Typography
                component="h2"
                variant="h5"
                sx={{ mb: 3, fontWeight: 600, textAlign: 'center' }}
              >
                Iniciar Sesión
              </Typography>

              {error && (
                <Alert severity="error" sx={{ mb: 3 }}>
                  {error}
                </Alert>
              )}

              <Box component="form" onSubmit={handleSubmit}>
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  id="email"
                  label="Correo Electrónico"
                  name="email"
                  autoComplete="email"
                  autoFocus
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  disabled={loading}
                  sx={{ mb: 2 }}
                />
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  name="password"
                  label="Contraseña"
                  type="password"
                  id="password"
                  autoComplete="current-password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={loading}
                  sx={{ mb: 3 }}
                />
                <Button
                  type="submit"
                  fullWidth
                  variant="contained"
                  size="large"
                  disabled={loading}
                  sx={{
                    py: 1.5,
                    fontSize: '1rem',
                    fontWeight: 600,
                    background: 'linear-gradient(135deg, #2563eb 0%, #3b82f6 100%)',
                    '&:hover': {
                      background: 'linear-gradient(135deg, #1e40af 0%, #2563eb 100%)',
                    },
                  }}
                >
                  {loading ? 'Iniciando sesión...' : 'Iniciar Sesión'}
                </Button>
              </Box>
            </CardContent>
          </Card>

          <Box
            sx={{
              mt: 4,
              p: 3,
              backgroundColor: 'rgba(255, 255, 255, 0.1)',
              backdropFilter: 'blur(10px)',
              borderRadius: 3,
              border: '1px solid rgba(255, 255, 255, 0.2)',
            }}
          >
            <Typography
              variant="caption"
              sx={{ color: 'white', display: 'block', textAlign: 'center', mb: 1 }}
            >
              <strong>Credenciales de prueba:</strong>
            </Typography>
            <Typography variant="body2" sx={{ color: 'rgba(255, 255, 255, 0.9)' }}>
              Email: admin@customermx.com
              <br />
              Password: admin123
            </Typography>
          </Box>
        </Box>
      </Container>
    </Box>
  );
};
