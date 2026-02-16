import React from 'react';
import { Card, CardContent, Typography, Box } from '@mui/material';
import { ConversionMetrics } from '../../types';

interface Props {
  data: ConversionMetrics;
}

export const ConversionFunnel: React.FC<Props> = ({ data }) => {
  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom sx={{ fontWeight: 600 }}>
          Funnel de Conversión
        </Typography>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
          {/* Level 1: Attendees */}
          <Box
            sx={{
              bgcolor: '#2563eb',
              p: 2,
              borderRadius: 2,
              color: 'white',
              width: '100%',
              transition: 'transform 0.2s',
              '&:hover': {
                transform: 'scale(1.02)',
              },
            }}
          >
            <Typography variant="h5" sx={{ fontWeight: 600 }}>
              {data.total_attendees.toLocaleString()}
            </Typography>
            <Typography variant="body2">Asistentes Totales</Typography>
          </Box>

          {/* Level 2: Leads */}
          <Box
            sx={{
              bgcolor: '#10b981',
              p: 2,
              borderRadius: 2,
              color: 'white',
              width: '75%',
              ml: 'auto',
              mr: 'auto',
              transition: 'transform 0.2s',
              '&:hover': {
                transform: 'scale(1.02)',
              },
            }}
          >
            <Typography variant="h5" sx={{ fontWeight: 600 }}>
              {data.total_leads.toLocaleString()}
            </Typography>
            <Typography variant="body2">
              Leads Capturados ({data.lead_rate.toFixed(1)}%)
            </Typography>
          </Box>

          {/* Level 3: Prospects */}
          <Box
            sx={{
              bgcolor: '#f59e0b',
              p: 2,
              borderRadius: 2,
              color: 'white',
              width: '50%',
              ml: 'auto',
              mr: 'auto',
              transition: 'transform 0.2s',
              '&:hover': {
                transform: 'scale(1.02)',
              },
            }}
          >
            <Typography variant="h5" sx={{ fontWeight: 600 }}>
              {data.total_prospects.toLocaleString()}
            </Typography>
            <Typography variant="body2">
              Prospectos Calificados ({data.prospect_rate.toFixed(1)}%)
            </Typography>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};
