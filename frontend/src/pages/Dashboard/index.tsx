import { useEffect, useState } from 'react';
import { Box, Card, CardContent, Typography, Grid, CircularProgress } from '@mui/material';
import {
  DirectionsCar as CarIcon,
  Business as BusinessIcon,
  Event as EventIcon,
  People as PeopleIcon,
} from '@mui/icons-material';
import { useAppSelector } from '../../hooks/useRedux';
import { apiService } from '../../services/api/apiService';
import { API_ENDPOINTS } from '../../services/api/apiConstants';

interface DashboardCounts {
  brands: number;
  vehicles: number;
  events: number;
  users: number;
}

export const Dashboard = () => {
  const { user, accessToken } = useAppSelector((state) => state.auth);
  const [counts, setCounts] = useState<DashboardCounts>({ brands: 0, vehicles: 0, events: 0, users: 0 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!accessToken) return;

    const fetchCounts = async () => {
      setLoading(true);

      const [brandsRes, vehiclesRes, eventsRes, usersRes] = await Promise.all([
        apiService.get<any[]>(API_ENDPOINTS.BRANDS.LIST, accessToken),
        apiService.get<any[]>(API_ENDPOINTS.VEHICLES.LIST, accessToken),
        apiService.get<any[]>(API_ENDPOINTS.EVENTS.LIST, accessToken),
        apiService.get<any[]>(API_ENDPOINTS.USERS.LIST, accessToken),
      ]);

      setCounts({
        brands: brandsRes.data?.length ?? 0,
        vehicles: vehiclesRes.data?.length ?? 0,
        events: eventsRes.data?.length ?? 0,
        users: usersRes.data?.length ?? 0,
      });

      setLoading(false);
    };

    fetchCounts();
  }, [accessToken]);

  const stats = [
    {
      title: 'Marcas',
      value: counts.brands,
      icon: <BusinessIcon sx={{ fontSize: 32 }} />,
      gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    },
    {
      title: 'Vehículos',
      value: counts.vehicles,
      icon: <CarIcon sx={{ fontSize: 32 }} />,
      gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    },
    {
      title: 'Eventos',
      value: counts.events,
      icon: <EventIcon sx={{ fontSize: 32 }} />,
      gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    },
    {
      title: 'Usuarios',
      value: counts.users,
      icon: <PeopleIcon sx={{ fontSize: 32 }} />,
      gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    },
  ];

  return (
    <Box>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 1 }}>
          Bienvenido, {user?.name}
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Panel de control de CustomerMX
        </Typography>
      </Box>

      <Grid container spacing={3}>
        {stats.map((stat) => (
          <Grid xs={12} sm={6} md={3} key={stat.title}>
            <Card
              sx={{
                position: 'relative',
                overflow: 'hidden',
                '&::before': {
                  content: '""',
                  position: 'absolute',
                  top: 0,
                  left: 0,
                  right: 0,
                  height: '4px',
                  background: stat.gradient,
                },
              }}
            >
              <CardContent sx={{ p: 3 }}>
                <Box display="flex" alignItems="center" justifyContent="space-between">
                  <Box>
                    <Typography
                      variant="body2"
                      color="text.secondary"
                      sx={{ mb: 1, fontWeight: 500 }}
                    >
                      {stat.title}
                    </Typography>
                    <Typography variant="h3" sx={{ fontWeight: 700 }}>
                      {loading ? <CircularProgress size={24} /> : stat.value}
                    </Typography>
                  </Box>
                  <Box
                    sx={{
                      width: 64,
                      height: 64,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      borderRadius: '16px',
                      background: stat.gradient,
                      color: 'white',
                      boxShadow: '0 8px 16px rgba(0, 0, 0, 0.1)',
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

      <Card sx={{ mt: 4 }}>
        <CardContent sx={{ p: 4 }}>
          <Typography variant="h6" sx={{ fontWeight: 600, mb: 3 }}>
            Información del Usuario
          </Typography>
          <Grid container spacing={3}>
            <Grid xs={12} md={4}>
              <Box>
                <Typography variant="caption" color="text.secondary" sx={{ mb: 0.5 }}>
                  Correo Electrónico
                </Typography>
                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                  {user?.email}
                </Typography>
              </Box>
            </Grid>
            <Grid xs={12} md={4}>
              <Box>
                <Typography variant="caption" color="text.secondary" sx={{ mb: 0.5 }}>
                  Rol
                </Typography>
                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                  {user?.role}
                </Typography>
              </Box>
            </Grid>
            <Grid xs={12} md={4}>
              <Box>
                <Typography variant="caption" color="text.secondary" sx={{ mb: 0.5 }}>
                  Estado
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
                  <Box
                    sx={{
                      width: 8,
                      height: 8,
                      borderRadius: '50%',
                      backgroundColor: user?.is_active ? 'success.main' : 'error.main',
                    }}
                  />
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {user?.is_active ? 'Activo' : 'Inactivo'}
                  </Typography>
                </Box>
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>
    </Box>
  );
};
