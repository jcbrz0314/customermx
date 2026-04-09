import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Grid,
  Alert,
  FormControlLabel,
  Checkbox,
  Rating,
  Typography,
  Box,
  InputAdornment,
  Divider,
  Stack,
} from '@mui/material';
import {
  Group as PeopleIcon,
  TrendingUp as TrendingIcon,
  Star as StarIcon,
} from '@mui/icons-material';
import { EventReport } from '../../types';
import { apiService, API_ENDPOINTS } from '../../services/api';
import { useAppSelector } from '../../hooks/useRedux';

interface EventReportDialogProps {
  open: boolean;
  onClose: () => void;
  eventId: string;
  existingReport?: EventReport | null;
  onSuccess: () => void;
}

export const EventReportDialog = ({
  open,
  onClose,
  eventId,
  existingReport,
  onSuccess,
}: EventReportDialogProps) => {
  const { accessToken } = useAppSelector((state) => state.auth);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Form fields
  const [hostessCount, setHostessCount] = useState('');
  const [setupVendor, setSetupVendor] = useState('');
  const [hasPromotional, setHasPromotional] = useState(false);
  const [attendees, setAttendees] = useState('');
  const [activitiesCount, setActivitiesCount] = useState('');
  const [leadsCollected, setLeadsCollected] = useState('');
  const [prospects, setProspects] = useState('');
  const [dealerRating, setDealerRating] = useState<number | null>(null);
  const [comments, setComments] = useState('');
  const [completed, setCompleted] = useState(false);

  // Pre-fill form if editing existing report
  useEffect(() => {
    if (open && existingReport) {
      setHostessCount(existingReport.hostess_count?.toString() || '');
      setSetupVendor(existingReport.setup_vendor || '');
      setHasPromotional(existingReport.has_promotional || false);
      setAttendees(existingReport.attendees?.toString() || '');
      setActivitiesCount(existingReport.activities_count?.toString() || '');
      setLeadsCollected(existingReport.leads_collected?.toString() || '');
      setProspects(existingReport.prospects?.toString() || '');
      setDealerRating(existingReport.dealer_rating || null);
      setComments(existingReport.comments || '');
      setCompleted(existingReport.completed || false);
    } else if (open) {
      // Reset form for new report
      resetForm();
    }
  }, [open, existingReport]);

  const resetForm = () => {
    setHostessCount('');
    setSetupVendor('');
    setHasPromotional(false);
    setAttendees('');
    setActivitiesCount('');
    setLeadsCollected('');
    setProspects('');
    setDealerRating(null);
    setComments('');
    setCompleted(false);
  };

  const handleSave = async () => {
    if (!accessToken) return;

    // Validate dealer_rating if provided
    if (dealerRating && (dealerRating < 1 || dealerRating > 5)) {
      setError('La calificación debe estar entre 1 y 5');
      return;
    }

    setLoading(true);
    setError('');

    // Build request with optional fields
    const requestData: any = {};

    if (hostessCount) requestData.hostess_count = parseInt(hostessCount);
    if (setupVendor) requestData.setup_vendor = setupVendor;
    requestData.has_promotional = hasPromotional;
    if (attendees) requestData.attendees = parseInt(attendees);
    if (activitiesCount) requestData.activities_count = parseInt(activitiesCount);
    if (leadsCollected) requestData.leads_collected = parseInt(leadsCollected);
    if (prospects) requestData.prospects = parseInt(prospects);
    if (dealerRating) requestData.dealer_rating = dealerRating;
    if (comments) requestData.comments = comments;
    requestData.completed = completed;

    const response = await apiService.post(
      API_ENDPOINTS.EVENTS.REPORT(eventId),
      requestData,
      accessToken
    );

    if (response.error) {
      setError(response.error);
      setLoading(false);
    } else {
      setLoading(false);
      resetForm();
      onSuccess();
      onClose();
    }
  };

  const handleClose = () => {
    if (!loading) {
      resetForm();
      setError('');
      onClose();
    }
  };

  const SectionLabel = ({ children }: { children: string }) => (
    <Box sx={{ mb: 2 }}>
      <Typography variant="overline" fontWeight={700} color="primary" letterSpacing={1}>
        {children}
      </Typography>
      <Divider sx={{ mt: 0.5 }} />
    </Box>
  );

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle sx={{ fontWeight: 700, fontSize: '1.2rem', pb: 1 }}>
        {existingReport ? 'Editar Reporte del Evento' : 'Reporte del Evento'}
      </DialogTitle>

      <DialogContent dividers sx={{ px: 3, py: 2.5 }}>
        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}

        {/* Personal y Logística */}
        <Box sx={{ mb: 3 }}>
          <SectionLabel>Personal y Logística</SectionLabel>
          <Grid container spacing={2}>
            <Grid xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Número de Edecanes"
                value={hostessCount}
                onChange={(e) => setHostessCount(e.target.value)}
                disabled={loading}
                inputProps={{ min: 0 }}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <PeopleIcon fontSize="small" />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>

            <Grid xs={12} sm={6}>
              <TextField
                fullWidth
                label="Proveedor de Montaje"
                value={setupVendor}
                onChange={(e) => setSetupVendor(e.target.value)}
                disabled={loading}
              />
            </Grid>

            <Grid xs={12}>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={hasPromotional}
                    onChange={(e) => setHasPromotional(e.target.checked)}
                    disabled={loading}
                  />
                }
                label="¿Se entregó material promocional?"
              />
            </Grid>
          </Grid>
        </Box>

        {/* Métricas del Evento */}
        <Box sx={{ mb: 3 }}>
          <SectionLabel>Métricas del Evento</SectionLabel>
          <Grid container spacing={2}>
            <Grid xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Asistentes"
                value={attendees}
                onChange={(e) => setAttendees(e.target.value)}
                disabled={loading}
                inputProps={{ min: 0 }}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <PeopleIcon fontSize="small" />
                    </InputAdornment>
                  ),
                }}
              />
            </Grid>

            <Grid xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Dinámicas"
                value={activitiesCount}
                onChange={(e) => setActivitiesCount(e.target.value)}
                disabled={loading}
                inputProps={{ min: 0 }}
              />
            </Grid>
          </Grid>
        </Box>

        {/* Conversiones */}
        <Box sx={{ mb: 3 }}>
          <SectionLabel>Conversiones y Oportunidades</SectionLabel>
          <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
            <TextField
              fullWidth
              type="number"
              label="Datos Levantados"
              value={leadsCollected}
              onChange={(e) => setLeadsCollected(e.target.value)}
              disabled={loading}
              inputProps={{ min: 0 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <TrendingIcon fontSize="small" />
                  </InputAdornment>
                ),
              }}
            />
            <TextField
              fullWidth
              type="number"
              label="Prospectos"
              value={prospects}
              onChange={(e) => setProspects(e.target.value)}
              disabled={loading}
              inputProps={{ min: 0 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <TrendingIcon fontSize="small" />
                  </InputAdornment>
                ),
              }}
            />
            <Box sx={{ flexShrink: 0, minWidth: { xs: '100%', sm: 180 } }}>
              <Typography variant="caption" color="text.secondary" display="block" sx={{ mb: 0.5 }}>
                Calificación del Distribuidor
              </Typography>
              <Box display="flex" alignItems="center" gap={1}>
                <Rating
                  value={dealerRating}
                  onChange={(_, newValue) => setDealerRating(newValue)}
                  disabled={loading}
                  size="large"
                  icon={<StarIcon fontSize="inherit" />}
                  emptyIcon={<StarIcon fontSize="inherit" />}
                />
                {dealerRating && (
                  <Typography variant="body2" color="text.secondary">
                    {dealerRating}/5
                  </Typography>
                )}
              </Box>
            </Box>
          </Stack>
        </Box>

        {/* Comentarios */}
        <Box sx={{ mb: 3 }}>
          <SectionLabel>Comentarios Adicionales</SectionLabel>
          <TextField
            fullWidth
            multiline
            rows={3}
            label="Comentarios"
            value={comments}
            onChange={(e) => setComments(e.target.value)}
            disabled={loading}
            placeholder="Detalles adicionales sobre el evento, incidentes, observaciones, etc."
          />
        </Box>

        {/* Marcar como completado */}
        <Box
          sx={{
            p: 2,
            bgcolor: completed ? 'success.50' : 'grey.50',
            borderRadius: 2,
            border: '1px solid',
            borderColor: completed ? 'success.main' : 'divider',
          }}
        >
          <FormControlLabel
            control={
              <Checkbox
                checked={completed}
                onChange={(e) => setCompleted(e.target.checked)}
                disabled={loading}
                color="success"
              />
            }
            label={
              <Box>
                <Typography variant="body2" fontWeight={600}>
                  Marcar reporte como completado
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  Indica que todos los datos del reporte han sido capturados y verificados
                </Typography>
              </Box>
            }
          />
        </Box>
      </DialogContent>

      <DialogActions sx={{ px: 3, py: 2, gap: 1 }}>
        <Button onClick={handleClose} variant="outlined" disabled={loading}>
          Cancelar
        </Button>
        <Button onClick={handleSave} variant="contained" disabled={loading}>
          {loading ? 'Guardando...' : 'Guardar Reporte'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
