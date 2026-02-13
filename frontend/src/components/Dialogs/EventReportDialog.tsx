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

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
      <DialogTitle sx={{ fontWeight: 600 }}>
        {existingReport ? 'Editar Reporte del Evento' : 'Crear Reporte del Evento'}
      </DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Grid container spacing={3} sx={{ mt: 0.5 }}>
          {/* Personal y Logística */}
          <Grid xs={12}>
            <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1 }}>
              Personal y Logística
            </Typography>
          </Grid>

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
                    <PeopleIcon />
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

          {/* Métricas del Evento */}
          <Grid xs={12}>
            <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1, mt: 2 }}>
              Métricas del Evento
            </Typography>
          </Grid>

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
                    <PeopleIcon />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>

          <Grid xs={12} sm={6}>
            <TextField
              fullWidth
              type="number"
              label="Actividades Realizadas"
              value={activitiesCount}
              onChange={(e) => setActivitiesCount(e.target.value)}
              disabled={loading}
              inputProps={{ min: 0 }}
            />
          </Grid>

          {/* Conversiones */}
          <Grid xs={12}>
            <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1, mt: 2 }}>
              Conversiones y Oportunidades
            </Typography>
          </Grid>

          <Grid xs={12} sm={4}>
            <TextField
              fullWidth
              type="number"
              label="Leads Capturados"
              value={leadsCollected}
              onChange={(e) => setLeadsCollected(e.target.value)}
              disabled={loading}
              inputProps={{ min: 0 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <TrendingIcon />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>

          <Grid xs={12} sm={4}>
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
                    <TrendingIcon />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>

          {/* Calificación del Distribuidor */}
          <Grid xs={12} sm={4}>
            <Box>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
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
                    ({dealerRating}/5)
                  </Typography>
                )}
              </Box>
            </Box>
          </Grid>

          {/* Comentarios */}
          <Grid xs={12}>
            <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1, mt: 2 }}>
              Comentarios Adicionales
            </Typography>
          </Grid>

          <Grid xs={12}>
            <TextField
              fullWidth
              multiline
              rows={4}
              label="Comentarios"
              value={comments}
              onChange={(e) => setComments(e.target.value)}
              disabled={loading}
              placeholder="Detalles adicionales sobre el evento, incidentes, observaciones, etc."
            />
          </Grid>

          {/* Marcar como completado */}
          <Grid xs={12}>
            <Box
              sx={{
                p: 2,
                bgcolor: completed ? 'success.50' : 'background.default',
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
                    <Typography variant="body1" fontWeight={600}>
                      Marcar reporte como completado
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Indica que todos los datos del reporte han sido capturados y verificados
                    </Typography>
                  </Box>
                }
              />
            </Box>
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2 }}>
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
