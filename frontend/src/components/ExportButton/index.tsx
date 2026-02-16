import React, { useState } from 'react';
import { Button, CircularProgress } from '@mui/material';
import { PictureAsPdf } from '@mui/icons-material';

interface ExportButtonProps {
  onExport: () => Promise<void>;
  label?: string;
  variant?: 'contained' | 'outlined' | 'text';
  color?: 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning';
  size?: 'small' | 'medium' | 'large';
  fullWidth?: boolean;
  disabled?: boolean;
}

export const ExportButton: React.FC<ExportButtonProps> = ({
  onExport,
  label = 'Exportar a PDF',
  variant = 'outlined',
  color = 'primary',
  size = 'medium',
  fullWidth = false,
  disabled = false,
}) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleExport = async () => {
    try {
      setLoading(true);
      setError(null);
      await onExport();
    } catch (err) {
      console.error('Error al exportar:', err);
      setError('Error al generar el PDF');
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Button
        variant={variant}
        color={error ? 'error' : color}
        size={size}
        fullWidth={fullWidth}
        disabled={disabled || loading}
        onClick={handleExport}
        startIcon={loading ? <CircularProgress size={16} /> : <PictureAsPdf />}
      >
        {loading ? 'Generando PDF...' : error || label}
      </Button>
    </>
  );
};
