import { Box, Card, CardContent, Typography, Grid } from '@mui/material';
import {
  DirectionsCar as CarIcon,
  Business as BusinessIcon,
  Event as EventIcon,
  People as PeopleIcon,
} from '@mui/icons-material';
import { useAppSelector } from '../../hooks/useRedux';

export const Dashboard = () => {
  const { user } = useAppSelector((state) => state.auth);

  const stats = [
    { title: 'Marcas', value: '4', icon: <BusinessIcon />, color: '#1976d2' },
    { title: 'Vehículos', value: '39', icon: <CarIcon />, color: '#2e7d32' },
    { title: 'Eventos', value: '0', icon: <EventIcon />, color: '#ed6c02' },
    { title: 'Usuarios', value: '1', icon: <PeopleIcon />, color: '#9c27b0' },
  ];

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Bienvenido, {user?.name}
      </Typography>
      <Typography variant="body1" color="text.secondary" paragraph>
        Panel de control de CustomerMX
      </Typography>

      <Grid container spacing={3} sx={{ mt: 2 }}>
        {stats.map((stat) => (
          <Grid xs={12} sm={6} md={3} key={stat.title}>
            <Card>
              <CardContent>
                <Box display="flex" alignItems="center" justifyContent="space-between">
                  <Box>
                    <Typography color="text.secondary" gutterBottom>
                      {stat.title}
                    </Typography>
                    <Typography variant="h4">{stat.value}</Typography>
                  </Box>
                  <Box
                    sx={{
                      p: 2,
                      borderRadius: 2,
                      backgroundColor: `${stat.color}20`,
                      color: stat.color,
                    }}
                  >
                    {stat.icon}
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Información del Usuario
          </Typography>
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2">
              <strong>Email:</strong> {user?.email}
            </Typography>
            <Typography variant="body2" sx={{ mt: 1 }}>
              <strong>Rol:</strong> {user?.role}
            </Typography>
            <Typography variant="body2" sx={{ mt: 1 }}>
              <strong>Estado:</strong> {user?.is_active ? 'Activo' : 'Inactivo'}
            </Typography>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
};
