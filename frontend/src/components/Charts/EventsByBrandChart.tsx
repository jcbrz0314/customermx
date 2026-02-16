import React from 'react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { BrandMetrics } from '../../types';

interface Props {
  data: BrandMetrics[];
  onBrandClick?: (brandId: string, brandName: string) => void;
}

export const EventsByBrandChart: React.FC<Props> = ({ data, onBrandClick }) => {
  const handleClick = (data: any) => {
    if (onBrandClick && data && data.activePayload && data.activePayload[0]) {
      const brand = data.activePayload[0].payload as BrandMetrics;
      onBrandClick(brand.brand_id, brand.brand_name);
    }
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data} onClick={handleClick}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis
          dataKey="brand_name"
          tick={{ fontSize: 12 }}
          stroke="#6b7280"
        />
        <YAxis tick={{ fontSize: 12 }} stroke="#6b7280" />
        <Tooltip
          contentStyle={{
            backgroundColor: '#fff',
            border: '1px solid #e5e7eb',
            borderRadius: '8px',
          }}
          cursor={{ fill: 'rgba(37, 99, 235, 0.1)' }}
        />
        <Legend />
        <Bar
          dataKey="event_count"
          fill="#2563eb"
          name="Eventos"
          radius={[8, 8, 0, 0]}
          style={{ cursor: onBrandClick ? 'pointer' : 'default' }}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};
