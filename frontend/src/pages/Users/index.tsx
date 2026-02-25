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
  Edit as EditIcon,
  Block as BlockIcon,
  CheckCircle as CheckCircleIcon,
  Search as SearchIcon,
  PersonAdd as PersonAddIcon,
  Send as SendIcon,
  Delete as DeleteIcon,
  HourglassEmpty as HourglassEmptyIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useAppSelector } from '../../hooks/useRedux';
import { useSnackbar } from '../../hooks/useSnackbar';
import { apiService } from '../../services/api/apiService';
import { API_ENDPOINTS } from '../../services/api/apiConstants';
import { User, UserRole, Invitation } from '../../types';
import { ChangeRoleDialog } from '../../components/Dialogs/ChangeRoleDialog';
import { DeactivateUserDialog } from '../../components/Dialogs/DeactivateUserDialog';
import { InviteUserDialog } from '../../components/Dialogs/InviteUserDialog';
import { ResendInvitationDialog } from '../../components/Dialogs/ResendInvitationDialog';
import { ConfirmDialog } from '../../components/Dialogs/ConfirmDialog';

type StatusFilter = 'active' | 'inactive' | 'pending' | '';

export const Users: React.FC = () => {
  const navigate = useNavigate();
  const { showSnackbar } = useSnackbar();
  const { user: currentUser, accessToken } = useAppSelector((state) => state.auth);

  const [users, setUsers] = useState<User[]>([]);
  const [pendingInvitations, setPendingInvitations] = useState<Invitation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const [roleFilter, setRoleFilter] = useState<UserRole | ''>('');
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('');
  const [searchFilter, setSearchFilter] = useState('');

  // Dialogs
  const [changeRoleDialog, setChangeRoleDialog] = useState<{
    open: boolean;
    user: User | null;
  }>({ open: false, user: null });
  const [deactivateDialog, setDeactivateDialog] = useState<{
    open: boolean;
    user: User | null;
  }>({ open: false, user: null });
  const [inviteDialogOpen, setInviteDialogOpen] = useState(false);
  const [resendDialog, setResendDialog] = useState<{
    open: boolean;
    invitation: Invitation | null;
  }>({ open: false, invitation: null });
  const [deleteInvitationDialog, setDeleteInvitationDialog] = useState<{
    open: boolean;
    invitation: Invitation | null;
  }>({ open: false, invitation: null });

  useEffect(() => {
    if (accessToken) fetchData();
  }, [accessToken]);

  const fetchData = async () => {
    try {
      setLoading(true);
      setError(null);
      const [usersRes, invitationsRes] = await Promise.all([
        apiService.get(API_ENDPOINTS.USERS.LIST, accessToken || undefined),
        apiService.get(`${API_ENDPOINTS.INVITATIONS.LIST}?pending=true`, accessToken || undefined),
      ]);
      if (usersRes.error) {
        setError(usersRes.error);
        showSnackbar('Error al cargar usuarios', 'error');
      } else {
        setUsers(Array.isArray(usersRes.data) ? usersRes.data : []);
        setPendingInvitations(Array.isArray(invitationsRes.data) ? invitationsRes.data : []);
      }
    } catch (err: any) {
      setError(err.message || 'Error al cargar usuarios');
      showSnackbar('Error al cargar usuarios', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleChangeRole = async (userId: string, newRole: UserRole) => {
    try {
      await apiService.patch(API_ENDPOINTS.USERS.CHANGE_ROLE(userId), { role: newRole }, accessToken || undefined);
      showSnackbar('Rol actualizado exitosamente', 'success');
      fetchData();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al cambiar rol', 'error');
    }
  };

  const handleDeactivate = async (userId: string, currentStatus: boolean) => {
    try {
      await apiService.patch(API_ENDPOINTS.USERS.DEACTIVATE(userId), { is_active: !currentStatus }, accessToken || undefined);
      showSnackbar(
        currentStatus ? 'Usuario desactivado' : 'Usuario activado',
        'success'
      );
      fetchData();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al cambiar estado', 'error');
    }
  };

  const handleInviteUser = async (
    email: string,
    role: UserRole,
    brandId: string | null
  ) => {
    try {
      await apiService.post(API_ENDPOINTS.INVITATIONS.CREATE, { email, role, brand_id: brandId }, accessToken || undefined);
      showSnackbar('Invitación enviada exitosamente', 'success');
      fetchData();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al enviar invitación', 'error');
    }
  };

  const handleResendInvitation = async (invitationId: string) => {
    try {
      await apiService.post(API_ENDPOINTS.INVITATIONS.RESEND(invitationId), {}, accessToken || undefined);
      showSnackbar('Invitación reenviada exitosamente', 'success');
      fetchData();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al reenviar invitación', 'error');
    }
  };

  const handleDeleteInvitation = async (invitationId: string) => {
    try {
      await apiService.delete(API_ENDPOINTS.INVITATIONS.DELETE(invitationId), accessToken || undefined);
      showSnackbar('Invitación cancelada', 'success');
      fetchData();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al cancelar invitación', 'error');
    }
  };

  const getRoleLabel = (role: UserRole): string => {
    const labels: Record<UserRole, string> = {
      ADMIN: 'Administrador',
      COORDINATOR: 'Coordinador',
      BRAND: 'Marca',
    };
    return labels[role];
  };

  const getRoleColor = (
    role: UserRole
  ): 'error' | 'primary' | 'secondary' | 'default' => {
    const colors: Record<UserRole, 'error' | 'primary' | 'secondary'> = {
      ADMIN: 'error',
      COORDINATOR: 'primary',
      BRAND: 'secondary',
    };
    return colors[role];
  };

  // Filter users
  const filteredUsers = statusFilter === 'pending' ? [] : users.filter((user) => {
    if (roleFilter && user.role !== roleFilter) return false;
    if (statusFilter === 'active' && !user.is_active) return false;
    if (statusFilter === 'inactive' && user.is_active) return false;
    if (searchFilter) {
      const search = searchFilter.toLowerCase();
      return (
        user.name.toLowerCase().includes(search) ||
        user.email.toLowerCase().includes(search)
      );
    }
    return true;
  });

  // Filter pending invitations
  const filteredInvitations = (statusFilter === 'active' || statusFilter === 'inactive') ? [] : pendingInvitations.filter((inv) => {
    if (roleFilter && inv.role !== roleFilter) return false;
    if (searchFilter) {
      return inv.email.toLowerCase().includes(searchFilter.toLowerCase());
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

  const totalRows = filteredUsers.length + filteredInvitations.length;

  return (
    <Box sx={{ p: 3 }}>
      {/* Header */}
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" sx={{ fontWeight: 600 }}>
          Gestión de Usuarios
        </Typography>
        <Button
          variant="contained"
          startIcon={<PersonAddIcon />}
          onClick={() => setInviteDialogOpen(true)}
        >
          Invitar Usuario
        </Button>
      </Box>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2}>
            <Grid size={{ xs: 12, sm: 4 }}>
              <TextField
                fullWidth
                size="small"
                placeholder="Buscar por nombre o email..."
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
            <Grid size={{ xs: 12, sm: 4 }}>
              <FormControl fullWidth size="small">
                <InputLabel>Rol</InputLabel>
                <Select
                  value={roleFilter}
                  onChange={(e) => setRoleFilter(e.target.value as UserRole | '')}
                  label="Rol"
                >
                  <MenuItem value="">Todos</MenuItem>
                  <MenuItem value="ADMIN">Administrador</MenuItem>
                  <MenuItem value="COORDINATOR">Coordinador</MenuItem>
                  <MenuItem value="BRAND">Marca</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid size={{ xs: 12, sm: 4 }}>
              <FormControl fullWidth size="small">
                <InputLabel>Estado</InputLabel>
                <Select
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value as StatusFilter)}
                  label="Estado"
                >
                  <MenuItem value="">Todos</MenuItem>
                  <MenuItem value="active">Activos</MenuItem>
                  <MenuItem value="inactive">Inactivos</MenuItem>
                  <MenuItem value="pending">Invitación pendiente</MenuItem>
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
          ) : totalRows === 0 ? (
            <Alert severity="info">No se encontraron usuarios</Alert>
          ) : (
            <TableContainer component={Paper} elevation={0}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Nombre</TableCell>
                    <TableCell>Email</TableCell>
                    <TableCell>Rol</TableCell>
                    <TableCell>Marca</TableCell>
                    <TableCell>Estado</TableCell>
                    <TableCell align="right">Acciones</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {/* Active/Inactive users */}
                  {filteredUsers.map((user) => (
                    <TableRow key={`user-${user.id}`} hover>
                      <TableCell>{user.name}</TableCell>
                      <TableCell>{user.email}</TableCell>
                      <TableCell>
                        <Chip
                          label={getRoleLabel(user.role)}
                          color={getRoleColor(user.role)}
                          size="small"
                        />
                      </TableCell>
                      <TableCell>{user.brand_name || '-'}</TableCell>
                      <TableCell>
                        <Chip
                          label={user.is_active ? 'Activo' : 'Inactivo'}
                          color={user.is_active ? 'success' : 'default'}
                          size="small"
                          icon={
                            user.is_active ? <CheckCircleIcon /> : <BlockIcon />
                          }
                        />
                      </TableCell>
                      <TableCell align="right">
                        <Tooltip title="Editar">
                          <IconButton
                            size="small"
                            onClick={() => navigate(`/users/${user.id}/edit`)}
                          >
                            <EditIcon />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Cambiar Rol">
                          <IconButton
                            size="small"
                            onClick={() =>
                              setChangeRoleDialog({ open: true, user })
                            }
                          >
                            <PersonAddIcon />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title={user.is_active ? 'Desactivar' : 'Activar'}>
                          <IconButton
                            size="small"
                            onClick={() =>
                              setDeactivateDialog({ open: true, user })
                            }
                          >
                            {user.is_active ? <BlockIcon /> : <CheckCircleIcon />}
                          </IconButton>
                        </Tooltip>
                      </TableCell>
                    </TableRow>
                  ))}

                  {/* Pending invitations */}
                  {filteredInvitations.map((invitation) => (
                    <TableRow
                      key={`inv-${invitation.id}`}
                      hover
                      sx={{ bgcolor: 'warning.50', opacity: 0.9 }}
                    >
                      <TableCell sx={{ color: 'text.secondary', fontStyle: 'italic' }}>
                        —
                      </TableCell>
                      <TableCell>{invitation.email}</TableCell>
                      <TableCell>
                        <Chip
                          label={getRoleLabel(invitation.role)}
                          color={getRoleColor(invitation.role)}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>{invitation.brand_name || '-'}</TableCell>
                      <TableCell>
                        <Chip
                          label="Invitación enviada"
                          color="warning"
                          size="small"
                          icon={<HourglassEmptyIcon />}
                        />
                      </TableCell>
                      <TableCell align="right">
                        <Tooltip title="Reenviar invitación">
                          <IconButton
                            size="small"
                            onClick={() => setResendDialog({ open: true, invitation })}
                          >
                            <SendIcon />
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Cancelar invitación">
                          <IconButton
                            size="small"
                            onClick={() => setDeleteInvitationDialog({ open: true, invitation })}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Tooltip>
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
      {changeRoleDialog.user && (
        <ChangeRoleDialog
          open={changeRoleDialog.open}
          onClose={() => setChangeRoleDialog({ open: false, user: null })}
          onConfirm={(newRole) =>
            handleChangeRole(changeRoleDialog.user!.id, newRole)
          }
          currentRole={changeRoleDialog.user.role}
          userName={changeRoleDialog.user.name}
        />
      )}

      {deactivateDialog.user && (
        <DeactivateUserDialog
          open={deactivateDialog.open}
          onClose={() => setDeactivateDialog({ open: false, user: null })}
          onConfirm={() =>
            handleDeactivate(
              deactivateDialog.user!.id,
              deactivateDialog.user!.is_active
            )
          }
          userName={deactivateDialog.user.name}
          userEmail={deactivateDialog.user.email}
          isActive={deactivateDialog.user.is_active}
        />
      )}

      <InviteUserDialog
        open={inviteDialogOpen}
        onClose={() => setInviteDialogOpen(false)}
        onConfirm={handleInviteUser}
      />

      {resendDialog.invitation && (
        <ResendInvitationDialog
          open={resendDialog.open}
          onClose={() => setResendDialog({ open: false, invitation: null })}
          onConfirm={() => handleResendInvitation(resendDialog.invitation!.id)}
          invitationEmail={resendDialog.invitation.email}
        />
      )}

      {deleteInvitationDialog.invitation && (
        <ConfirmDialog
          open={deleteInvitationDialog.open}
          onClose={() => setDeleteInvitationDialog({ open: false, invitation: null })}
          onConfirm={() => handleDeleteInvitation(deleteInvitationDialog.invitation!.id)}
          title="Cancelar Invitación"
          message={`¿Estás seguro que deseas cancelar la invitación enviada a ${deleteInvitationDialog.invitation.email}?`}
          severity="error"
        />
      )}
    </Box>
  );
};
