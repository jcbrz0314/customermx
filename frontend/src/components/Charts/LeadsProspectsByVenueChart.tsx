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
import { VenueMetrics } from '../../types';

interface Props {
  data: VenueMetrics[];
}

export const LeadsProspectsByVenueChart: React.FC<Props> = ({ data }) => {
  return (
    <ResponsiveContainer width="100%" height={Math.max(300, data.length * 40)}>
      <BarChart data={data} layout="vertical" margin={{ left: 20 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis type="number" tick={{ fontSize: 12 }} stroke="#6b7280" />
        <YAxis
          dataKey="venue"
          type="category"
          width={180}
          tick={{ fontSize: 11 }}
          stroke="#6b7280"
        />
        <Tooltip
          contentStyle={{
            backgroundColor: '#fff',
            border: '1px solid #e5e7eb',
            borderRadius: '8px',
          }}
        />
        <Legend />
        <Bar
          dataKey="total_leads"
          fill="#f59e0b"
          name="Datos Levantados"
          radius={[0, 8, 8, 0]}
        />
        <Bar
          dataKey="total_prospects"
          fill="#10b981"
          name="Prospectos"
          radius={[0, 8, 8, 0]}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};
