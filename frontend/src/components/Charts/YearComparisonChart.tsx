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
import { YearComparison } from '../../types';

interface Props {
  data: YearComparison[];
  onYearClick?: (year: number) => void;
}

export const YearComparisonChart: React.FC<Props> = ({ data, onYearClick }) => {
  const handleClick = (data: any) => {
    if (onYearClick && data && data.activePayload && data.activePayload[0]) {
      const yearData = data.activePayload[0].payload as YearComparison;
      onYearClick(yearData.year);
    }
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data} onClick={handleClick}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis
          dataKey="year"
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
          style={{ cursor: onYearClick ? 'pointer' : 'default' }}
        />
        <Bar
          dataKey="total_attendees"
          fill="#10b981"
          name="Asistentes"
          radius={[8, 8, 0, 0]}
          style={{ cursor: onYearClick ? 'pointer' : 'default' }}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};
