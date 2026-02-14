import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Tooltip,
  CircularProgress,
  Alert,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  InputAdornment,
} from '@mui/material';
import {
  PersonAdd as PersonAddIcon,
  Send as SendIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  CheckCircle as CheckCircleIcon,
  Schedule as ScheduleIcon,
  Cancel as CancelIcon,
  HourglassEmpty as HourglassEmptyIcon,
} from '@mui/icons-material';
import { useAppSelector } from '../../hooks/useRedux';
import { useSnackbar } from '../../hooks/useSnackbar';
import { apiService } from '../../services/api/apiService';
import { API_ENDPOINTS } from '../../services/api/apiConstants';
import { Invitation, InvitationStatus } from '../../types';
import { InviteUserDialog } from '../../components/Dialogs/InviteUserDialog';
import { ResendInvitationDialog } from '../../components/Dialogs/ResendInvitationDialog';
import { ConfirmDialog } from '../../components/Dialogs/ConfirmDialog';

export const Invitations: React.FC = () => {
  const { showSnackbar } = useSnackbar();
  const currentUser = useAppSelector((state) => state.auth.user);

  const [invitations, setInvitations] = useState<Invitation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const [statusFilter, setStatusFilter] = useState<InvitationStatus | ''>('');
  const [searchFilter, setSearchFilter] = useState('');

  // Dialogs
  const [inviteDialogOpen, setInviteDialogOpen] = useState(false);
  const [resendDialog, setResendDialog] = useState<{
    open: boolean;
    invitation: Invitation | null;
  }>({ open: false, invitation: null });
  const [deleteDialog, setDeleteDialog] = useState<{
    open: boolean;
    invitation: Invitation | null;
  }>({ open: false, invitation: null });

  useEffect(() => {
    fetchInvitations();
  }, []);

  const fetchInvitations = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiService.get(API_ENDPOINTS.INVITATIONS.LIST);
      setInvitations(response.data || []);
    } catch (err: any) {
      setError(err.message || 'Error al cargar invitaciones');
      showSnackbar('Error al cargar invitaciones', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleInviteUser = async (
    email: string,
    role: string,
    brandId: string | null
  ) => {
    try {
      await apiService.post(API_ENDPOINTS.INVITATIONS.CREATE, {
        email,
        role,
        brand_id: brandId,
      });
      showSnackbar('Invitación enviada exitosamente', 'success');
      fetchInvitations();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al enviar invitación', 'error');
    }
  };

  const handleResend = async (invitationId: string) => {
    try {
      await apiService.post(API_ENDPOINTS.INVITATIONS.RESEND(invitationId), {});
      showSnackbar('Invitación reenviada exitosamente', 'success');
      fetchInvitations();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al reenviar invitación', 'error');
    }
  };

  const handleDelete = async (invitationId: string) => {
    try {
      await apiService.delete(API_ENDPOINTS.INVITATIONS.DELETE(invitationId));
      showSnackbar('Invitación eliminada', 'success');
      fetchInvitations();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al eliminar invitación', 'error');
    }
  };

  const getStatusLabel = (status: InvitationStatus): string => {
    const labels: Record<InvitationStatus, string> = {
      PENDING: 'Pendiente',
      ACCEPTED: 'Aceptada',
      EXPIRED: 'Expirada',
    };
    return labels[status];
  };

  const getStatusColor = (
    status: InvitationStatus
  ): 'warning' | 'success' | 'error' => {
    const colors: Record<InvitationStatus, 'warning' | 'success' | 'error'> = {
      PENDING: 'warning',
      ACCEPTED: 'success',
      EXPIRED: 'error',
    };
    return colors[status];
  };

  const getStatusIcon = (status: InvitationStatus) => {
    const icons: Record<InvitationStatus, React.ReactElement> = {
      PENDING: <HourglassEmptyIcon />,
      ACCEPTED: <CheckCircleIcon />,
      EXPIRED: <CancelIcon />,
    };
    return icons[status];
  };

  const getRoleLabel = (role: string): string => {
    const labels: Record<string, string> = {
      ADMIN: 'Administrador',
      COORDINATOR: 'Coordinador',
      BRAND: 'Marca',
    };
    return labels[role] || role;
  };

  const filteredInvitations = invitations.filter((invitation) => {
    if (statusFilter && invitation.status !== statusFilter) return false;
    if (searchFilter) {
      const search = searchFilter.toLowerCase();
      return invitation.email.toLowerCase().includes(search);
    }
    return true;
  });

  if (!currentUser || currentUser.role !== 'ADMIN') {
    return (
      <Box sx={{ p: 3 }}>
        <Alert severity="error">
          No tienes permisos para acceder a esta sección.
        </Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" sx={{ fontWeight: 600 }}>
          Invitaciones
        </Typography>
        <Button
          variant="contained"
          startIcon={<PersonAddIcon />}
          onClick={() => setInviteDialogOpen(true)}
        >
          Nueva Invitación
        </Button>
      </Box>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                placeholder="Buscar por email..."
                value={searchFilter}
                onChange={(e) => setSearchFilter(e.target.value)}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Estado</InputLabel>
                <Select
                  value={statusFilter}
                  onChange={(e) =>
                    setStatusFilter(e.target.value as InvitationStatus | '')
                  }
                  label="Estado"
                >
                  <MenuItem value="">Todos</MenuItem>
                  <MenuItem value="PENDING">Pendientes</MenuItem>
                  <MenuItem value="ACCEPTED">Aceptadas</MenuItem>
                  <MenuItem value="EXPIRED">Expiradas</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent>
          {loading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : error ? (
            <Alert severity="error">{error}</Alert>
          ) : filteredInvitations.length === 0 ? (
            <Alert severity="info">No se encontraron invitaciones</Alert>
          ) : (
            <TableContainer component={Paper} elevation={0}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Email</TableCell>
                    <TableCell>Rol</TableCell>
                    <TableCell>Marca</TableCell>
                    <TableCell>Estado</TableCell>
                    <TableCell>Fecha de Envío</TableCell>
                    <TableCell>Expira</TableCell>
                    <TableCell align="right">Acciones</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {filteredInvitations.map((invitation) => (
                    <TableRow key={invitation.id} hover>
                      <TableCell>{invitation.email}</TableCell>
                      <TableCell>
                        <Chip
                          label={getRoleLabel(invitation.role)}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>{invitation.brand_name || '-'}</TableCell>
                      <TableCell>
                        <Chip
                          label={getStatusLabel(invitation.status)}
                          color={getStatusColor(invitation.status)}
                          size="small"
                          icon={getStatusIcon(invitation.status)}
                        />
                      </TableCell>
                      <TableCell>
                        {new Date(invitation.created_at).toLocaleDateString()}
                      </TableCell>
                      <TableCell>
                        {new Date(invitation.expires_at).toLocaleDateString()}
                      </TableCell>
                      <TableCell align="right">
                        {invitation.status === 'PENDING' && (
                          <>
                            <Tooltip title="Reenviar">
                              <IconButton
                                size="small"
                                onClick={() =>
                                  setResendDialog({ open: true, invitation })
                                }
                              >
                                <SendIcon />
                              </IconButton>
                            </Tooltip>
                            <Tooltip title="Eliminar">
                              <IconButton
                                size="small"
                                onClick={() =>
                                  setDeleteDialog({ open: true, invitation })
                                }
                              >
                                <DeleteIcon />
                              </IconButton>
                            </Tooltip>
                          </>
                        )}
                        {invitation.status === 'EXPIRED' && (
                          <Tooltip title="Eliminar">
                            <IconButton
                              size="small"
                              onClick={() =>
                                setDeleteDialog({ open: true, invitation })
                              }
                            >
                              <DeleteIcon />
                            </IconButton>
                          </Tooltip>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </CardContent>
      </Card>

      {/* Dialogs */}
      <InviteUserDialog
        open={inviteDialogOpen}
        onClose={() => setInviteDialogOpen(false)}
        onConfirm={handleInviteUser}
      />

      {resendDialog.invitation && (
        <ResendInvitationDialog
          open={resendDialog.open}
          onClose={() => setResendDialog({ open: false, invitation: null })}
          onConfirm={() => handleResend(resendDialog.invitation!.id)}
          invitationEmail={resendDialog.invitation.email}
        />
      )}

      {deleteDialog.invitation && (
        <ConfirmDialog
          open={deleteDialog.open}
          onClose={() => setDeleteDialog({ open: false, invitation: null })}
          onConfirm={() => handleDelete(deleteDialog.invitation!.id)}
          title="Eliminar Invitación"
          message={`¿Estás seguro que deseas eliminar la invitación enviada a ${deleteDialog.invitation.email}?`}
          severity="error"
        />
      )}
    </Box>
  );
};
