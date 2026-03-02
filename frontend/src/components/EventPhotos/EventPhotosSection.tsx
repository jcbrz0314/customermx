import { useCallback, useEffect, useRef, useState } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  Tooltip,
  Typography,
} from '@mui/material';
import {
  AddPhotoAlternate as UploadIcon,
  Delete as DeleteIcon,
  Autorenew as ReplaceIcon,
  NavigateBefore as PrevIcon,
  NavigateNext as NextIcon,
  Close as CloseIcon,
} from '@mui/icons-material';
import { EventPhoto } from '../../types';
import { API_BASE_URL, API_ENDPOINTS } from '../../services/api/apiConstants';
import { useSnackbar } from '../../hooks/useSnackbar';

const MAX_PHOTOS = 10;
const MAX_FILE_SIZE_MB = 10;
const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];

interface EventPhotosSectionProps {
  eventId: string;
  photos: EventPhoto[];
  canWrite: boolean; // ADMIN or COORDINATOR
  accessToken: string;
  onPhotosChange: () => void;
}

export const EventPhotosSection = ({
  eventId,
  photos,
  canWrite,
  accessToken,
  onPhotosChange,
}: EventPhotosSectionProps) => {
  const { showSnackbar } = useSnackbar();
  const uploadInputRef = useRef<HTMLInputElement>(null);
  const replaceInputRef = useRef<HTMLInputElement>(null);

  const [uploading, setUploading] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [replacingPhoto, setReplacingPhoto] = useState<EventPhoto | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<EventPhoto | null>(null);
  const [lightboxIndex, setLightboxIndex] = useState<number | null>(null);

  const viewablePhotos = photos.filter((p) => deletingId !== p.id);

  const openLightbox = (photo: EventPhoto) => {
    const idx = viewablePhotos.findIndex((p) => p.id === photo.id);
    if (idx !== -1) setLightboxIndex(idx);
  };

  const closeLightbox = useCallback(() => setLightboxIndex(null), []);

  const goPrev = useCallback(() => {
    setLightboxIndex((i) => (i !== null ? (i - 1 + viewablePhotos.length) % viewablePhotos.length : null));
  }, [viewablePhotos.length]);

  const goNext = useCallback(() => {
    setLightboxIndex((i) => (i !== null ? (i + 1) % viewablePhotos.length : null));
  }, [viewablePhotos.length]);

  useEffect(() => {
    if (lightboxIndex === null) return;
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'ArrowLeft') goPrev();
      else if (e.key === 'ArrowRight') goNext();
      else if (e.key === 'Escape') closeLightbox();
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [lightboxIndex, goPrev, goNext, closeLightbox]);

  const authHeader = { Authorization: `Bearer ${accessToken}` };

  const validateFiles = (files: File[]): string | null => {
    for (const file of files) {
      if (!ACCEPTED_TYPES.includes(file.type)) {
        return `"${file.name}" no es una imagen válida (jpeg, png, gif, webp)`;
      }
      if (file.size > MAX_FILE_SIZE_MB * 1024 * 1024) {
        return `"${file.name}" excede el límite de ${MAX_FILE_SIZE_MB} MB`;
      }
    }
    return null;
  };

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || e.target.files.length === 0) return;
    const files = Array.from(e.target.files);

    const validationError = validateFiles(files);
    if (validationError) {
      showSnackbar(validationError, 'error');
      e.target.value = '';
      return;
    }

    if (photos.length + files.length > MAX_PHOTOS) {
      showSnackbar(
        `Solo puedes tener ${MAX_PHOTOS} fotos. Ya tienes ${photos.length}.`,
        'error'
      );
      e.target.value = '';
      return;
    }

    setUploading(true);
    try {
      const formData = new FormData();
      files.forEach((file) => formData.append('files', file));

      const response = await fetch(
        `${API_BASE_URL}${API_ENDPOINTS.EVENTS.PHOTOS(eventId)}`,
        { method: 'POST', headers: authHeader, body: formData }
      );

      if (!response.ok) {
        const json = await response.json().catch(() => ({}));
        throw new Error(json.error || 'Error al subir las fotos');
      }

      showSnackbar(`${files.length} foto(s) subida(s) exitosamente`, 'success');
      onPhotosChange();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al subir las fotos', 'error');
    } finally {
      setUploading(false);
      e.target.value = '';
    }
  };

  const handleDeleteConfirm = async () => {
    if (!confirmDelete) return;
    setDeletingId(confirmDelete.id);
    setConfirmDelete(null);

    try {
      const response = await fetch(
        `${API_BASE_URL}${API_ENDPOINTS.EVENTS.DELETE_PHOTO(eventId, confirmDelete.id)}`,
        { method: 'DELETE', headers: authHeader }
      );

      if (!response.ok) {
        const json = await response.json().catch(() => ({}));
        throw new Error(json.error || 'Error al eliminar la foto');
      }

      showSnackbar('Foto eliminada exitosamente', 'success');
      onPhotosChange();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al eliminar la foto', 'error');
    } finally {
      setDeletingId(null);
    }
  };

  const handleReplaceSelect = (photo: EventPhoto) => {
    setReplacingPhoto(photo);
    replaceInputRef.current?.click();
  };

  const handleReplaceUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!e.target.files || e.target.files.length === 0 || !replacingPhoto) return;
    const file = e.target.files[0];

    const validationError = validateFiles([file]);
    if (validationError) {
      showSnackbar(validationError, 'error');
      e.target.value = '';
      setReplacingPhoto(null);
      return;
    }

    setDeletingId(replacingPhoto.id); // reuse loading indicator
    try {
      const formData = new FormData();
      formData.append('file', file);

      const response = await fetch(
        `${API_BASE_URL}${API_ENDPOINTS.EVENTS.REPLACE_PHOTO(eventId, replacingPhoto.id)}`,
        { method: 'PUT', headers: authHeader, body: formData }
      );

      if (!response.ok) {
        const json = await response.json().catch(() => ({}));
        throw new Error(json.error || 'Error al reemplazar la foto');
      }

      showSnackbar('Foto reemplazada exitosamente', 'success');
      onPhotosChange();
    } catch (err: any) {
      showSnackbar(err.message || 'Error al reemplazar la foto', 'error');
    } finally {
      setDeletingId(null);
      setReplacingPhoto(null);
      e.target.value = '';
    }
  };

  return (
    <Card>
      <CardContent sx={{ p: 3 }}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Box>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>
              Fotos del Evento
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {photos.length}/{MAX_PHOTOS} fotos
            </Typography>
          </Box>
          {canWrite && photos.length < MAX_PHOTOS && (
            <Button
              size="small"
              variant="outlined"
              startIcon={uploading ? <CircularProgress size={14} /> : <UploadIcon />}
              disabled={uploading}
              onClick={() => uploadInputRef.current?.click()}
              sx={{ borderRadius: 2 }}
            >
              Subir Fotos
            </Button>
          )}
        </Box>

        {/* Hidden file inputs */}
        <input
          ref={uploadInputRef}
          type="file"
          accept="image/*"
          multiple
          style={{ display: 'none' }}
          onChange={handleUpload}
        />
        <input
          ref={replaceInputRef}
          type="file"
          accept="image/*"
          style={{ display: 'none' }}
          onChange={handleReplaceUpload}
        />

        {photos.length === 0 ? (
          <Typography variant="body2" color="text.secondary">
            No hay fotos adjuntas
          </Typography>
        ) : (
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: 'repeat(3, 1fr)',
              gap: 1,
            }}
          >
            {photos.map((photo) => {
              const isProcessing = deletingId === photo.id;
              return (
                <Box
                  key={photo.id}
                  sx={{
                    position: 'relative',
                    paddingTop: '75%',
                    borderRadius: 1,
                    overflow: 'hidden',
                    bgcolor: 'grey.100',
                    cursor: isProcessing ? 'default' : 'zoom-in',
                  }}
                  onClick={() => !isProcessing && openLightbox(photo)}
                >
                  {isProcessing ? (
                    <Box
                      sx={{
                        position: 'absolute',
                        inset: 0,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                      }}
                    >
                      <CircularProgress size={32} />
                    </Box>
                  ) : (
                    <>
                      <img
                        src={`${API_BASE_URL}${photo.url}`}
                        alt={photo.filename}
                        loading="lazy"
                        style={{
                          position: 'absolute',
                          inset: 0,
                          width: '100%',
                          height: '100%',
                          objectFit: 'cover',
                          pointerEvents: 'none', // let the parent Box handle clicks
                        }}
                      />
                      {canWrite && (
                        <Box
                          sx={{
                            position: 'absolute',
                            bottom: 0,
                            left: 0,
                            right: 0,
                            px: 1,
                            py: 0.5,
                            bgcolor: 'rgba(0,0,0,0.55)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'flex-end',
                            gap: 0.5,
                          }}
                          onClick={(e) => e.stopPropagation()} // don't open lightbox when clicking actions
                        >
                          <Tooltip title="Reemplazar">
                            <IconButton
                              size="small"
                              sx={{ color: 'white', p: 0.5 }}
                              onClick={() => handleReplaceSelect(photo)}
                            >
                              <ReplaceIcon sx={{ fontSize: 18 }} />
                            </IconButton>
                          </Tooltip>
                          <Tooltip title="Eliminar">
                            <IconButton
                              size="small"
                              sx={{ color: 'white', p: 0.5 }}
                              onClick={() => setConfirmDelete(photo)}
                            >
                              <DeleteIcon sx={{ fontSize: 18 }} />
                            </IconButton>
                          </Tooltip>
                        </Box>
                      )}
                    </>
                  )}
                </Box>
              );
            })}
          </Box>
        )}

        {/* Lightbox */}
        <Dialog
          open={lightboxIndex !== null}
          onClose={closeLightbox}
          fullScreen
          PaperProps={{
            sx: {
              bgcolor: 'rgba(0,0,0,0.93)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            },
          }}
        >
          {lightboxIndex !== null && (
            <>
              {/* Close */}
              <IconButton
                onClick={closeLightbox}
                sx={{
                  position: 'fixed',
                  top: 16,
                  right: 16,
                  color: 'white',
                  bgcolor: 'rgba(255,255,255,0.1)',
                  '&:hover': { bgcolor: 'rgba(255,255,255,0.25)' },
                  zIndex: 1,
                }}
              >
                <CloseIcon />
              </IconButton>

              {/* Counter */}
              <Typography
                variant="body2"
                sx={{
                  position: 'fixed',
                  top: 20,
                  left: '50%',
                  transform: 'translateX(-50%)',
                  color: 'rgba(255,255,255,0.7)',
                  zIndex: 1,
                }}
              >
                {lightboxIndex + 1} / {viewablePhotos.length}
              </Typography>

              {/* Prev */}
              {viewablePhotos.length > 1 && (
                <IconButton
                  onClick={goPrev}
                  sx={{
                    position: 'fixed',
                    left: 16,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    color: 'white',
                    bgcolor: 'rgba(255,255,255,0.1)',
                    '&:hover': { bgcolor: 'rgba(255,255,255,0.25)' },
                    zIndex: 1,
                  }}
                >
                  <PrevIcon sx={{ fontSize: 36 }} />
                </IconButton>
              )}

              {/* Image */}
              <img
                key={viewablePhotos[lightboxIndex].id}
                src={`${API_BASE_URL}${viewablePhotos[lightboxIndex].url}`}
                alt={viewablePhotos[lightboxIndex].filename}
                style={{
                  maxWidth: 'calc(100vw - 140px)',
                  maxHeight: 'calc(100vh - 80px)',
                  objectFit: 'contain',
                  borderRadius: 8,
                  display: 'block',
                }}
              />

              {/* Next */}
              {viewablePhotos.length > 1 && (
                <IconButton
                  onClick={goNext}
                  sx={{
                    position: 'fixed',
                    right: 16,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    color: 'white',
                    bgcolor: 'rgba(255,255,255,0.1)',
                    '&:hover': { bgcolor: 'rgba(255,255,255,0.25)' },
                    zIndex: 1,
                  }}
                >
                  <NextIcon sx={{ fontSize: 36 }} />
                </IconButton>
              )}
            </>
          )}
        </Dialog>

        {/* Delete confirmation dialog */}
        <Dialog open={Boolean(confirmDelete)} onClose={() => setConfirmDelete(null)}>
          <DialogTitle>Eliminar Foto</DialogTitle>
          <DialogContent>
            <DialogContentText>
              ¿Estás seguro de eliminar "{confirmDelete?.filename}"? Esta acción no se puede
              deshacer.
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setConfirmDelete(null)}>Cancelar</Button>
            <Button onClick={handleDeleteConfirm} color="error" variant="contained">
              Eliminar
            </Button>
          </DialogActions>
        </Dialog>
      </CardContent>
    </Card>
  );
};
