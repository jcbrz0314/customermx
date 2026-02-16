import React from 'react';
import { Card, CardContent, Typography, Box } from '@mui/material';
import { TrendingUp, TrendingDown } from '@mui/icons-material';

interface Props {
  title: string;
  value: string | number;
  subtitle?: string;
  icon?: React.ReactElement;
  trend?: {
    direction: 'up' | 'down';
    value: string;
  };
  color?: string;
}

export const KPICard: React.FC<Props> = ({
  title,
  value,
  subtitle,
  icon,
  trend,
  color = '#2563eb',
}) => {
  return (
    <Card
      sx={{
        height: '100%',
        transition: 'transform 0.2s, box-shadow 0.2s',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: 3,
        },
      }}
    >
      <CardContent>
        <Box
          sx={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'flex-start',
          }}
        >
          <Box sx={{ flex: 1 }}>
            <Typography
              color="text.secondary"
              variant="body2"
              gutterBottom
              sx={{ fontWeight: 500 }}
            >
              {title}
            </Typography>
            <Typography
              variant="h4"
              sx={{
                fontWeight: 600,
                color,
                mb: 0.5,
              }}
            >
              {value}
            </Typography>
            {subtitle && (
              <Typography variant="caption" color="text.secondary">
                {subtitle}
              </Typography>
            )}
            {trend && (
              <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                {trend.direction === 'up' ? (
                  <TrendingUp
                    sx={{ fontSize: 16, color: 'success.main', mr: 0.5 }}
                  />
                ) : (
                  <TrendingDown
                    sx={{ fontSize: 16, color: 'error.main', mr: 0.5 }}
                  />
                )}
                <Typography
                  variant="caption"
                  color={trend.direction === 'up' ? 'success.main' : 'error.main'}
                  sx={{ fontWeight: 600 }}
                >
                  {trend.value}
                </Typography>
              </Box>
            )}
          </Box>
          {icon && (
            <Box
              sx={{
                bgcolor: `${color}15`,
                borderRadius: 2,
                p: 1.5,
                color,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              {React.cloneElement(icon, { sx: { fontSize: 28 } })}
            </Box>
          )}
        </Box>
      </CardContent>
    </Card>
  );
};
