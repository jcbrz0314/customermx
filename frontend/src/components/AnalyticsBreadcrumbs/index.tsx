import React from 'react';
import { Breadcrumbs, Link, Typography, Chip, Box } from '@mui/material';
import { NavigateNext, Home, Close } from '@mui/icons-material';

export interface BreadcrumbFilter {
  type: 'brand' | 'year' | 'month' | 'state' | 'vehicle' | 'eventType' | 'dealer';
  label: string;
  value: string;
}

interface AnalyticsBreadcrumbsProps {
  filters: BreadcrumbFilter[];
  onRemoveFilter: (index: number) => void;
  onClearAll: () => void;
}

export const AnalyticsBreadcrumbs: React.FC<AnalyticsBreadcrumbsProps> = ({
  filters,
  onRemoveFilter,
  onClearAll,
}) => {
  if (filters.length === 0) {
    return null;
  }

  const getFilterColor = (type: BreadcrumbFilter['type']) => {
    switch (type) {
      case 'brand':
        return 'primary';
      case 'year':
        return 'secondary';
      case 'month':
        return 'info';
      case 'state':
        return 'success';
      case 'vehicle':
        return 'warning';
      case 'eventType':
        return 'error';
      case 'dealer':
        return 'default';
      default:
        return 'default';
    }
  };

  const getFilterTypeLabel = (type: BreadcrumbFilter['type']) => {
    switch (type) {
      case 'brand':
        return 'Marca';
      case 'year':
        return 'Año';
      case 'month':
        return 'Mes';
      case 'state':
        return 'Estado';
      case 'vehicle':
        return 'Vehículo';
      case 'eventType':
        return 'Tipo';
      case 'dealer':
        return 'Distribuidor';
      default:
        return type;
    }
  };

  return (
    <Box
      sx={{
        bgcolor: '#f9fafb',
        p: 2,
        borderRadius: 2,
        mb: 3,
        border: '1px solid #e5e7eb',
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
        <Typography variant="body2" sx={{ color: '#6b7280', fontWeight: 500 }}>
          Filtros activos:
        </Typography>

        <Breadcrumbs
          separator={<NavigateNext fontSize="small" sx={{ color: '#9ca3af' }} />}
          sx={{ flex: 1 }}
        >
          <Link
            component="button"
            onClick={onClearAll}
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
              color: '#6b7280',
              textDecoration: 'none',
              '&:hover': {
                color: '#2563eb',
                textDecoration: 'underline',
              },
            }}
          >
            <Home fontSize="small" />
            <Typography variant="body2">Inicio</Typography>
          </Link>

          {filters.map((filter, index) => (
            <Chip
              key={index}
              label={
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                  <Typography variant="caption" sx={{ fontWeight: 600 }}>
                    {getFilterTypeLabel(filter.type)}:
                  </Typography>
                  <Typography variant="caption">{filter.label}</Typography>
                </Box>
              }
              onDelete={() => onRemoveFilter(index)}
              deleteIcon={<Close sx={{ fontSize: 16 }} />}
              color={getFilterColor(filter.type)}
              size="small"
              sx={{
                height: 24,
                '& .MuiChip-deleteIcon': {
                  fontSize: 16,
                },
              }}
            />
          ))}
        </Breadcrumbs>

        {filters.length > 1 && (
          <Link
            component="button"
            onClick={onClearAll}
            sx={{
              fontSize: 12,
              color: '#ef4444',
              textDecoration: 'none',
              fontWeight: 500,
              '&:hover': {
                textDecoration: 'underline',
              },
            }}
          >
            Limpiar todo
          </Link>
        )}
      </Box>
    </Box>
  );
};
