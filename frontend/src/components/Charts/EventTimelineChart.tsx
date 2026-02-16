import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { MonthlyMetrics } from '../../types';

interface Props {
  data: MonthlyMetrics[];
  onMonthClick?: (year: number, month: number, monthName: string) => void;
}

export const EventTimelineChart: React.FC<Props> = ({ data, onMonthClick }) => {
  const handleClick = (data: any) => {
    if (onMonthClick && data && data.activePayload && data.activePayload[0]) {
      const monthData = data.activePayload[0].payload as MonthlyMetrics;
      onMonthClick(monthData.year, monthData.month, monthData.month_name);
    }
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={data} onClick={handleClick}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis
          dataKey="month_name"
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
          cursor={{ stroke: '#2563eb', strokeDasharray: '5 5' }}
        />
        <Legend />
        <Line
          type="monotone"
          dataKey="event_count"
          stroke="#2563eb"
          strokeWidth={2}
          name="Eventos"
          dot={{ fill: '#2563eb', r: 4 }}
          activeDot={{ r: 6 }}
          style={{ cursor: onMonthClick ? 'pointer' : 'default' }}
        />
        <Line
          type="monotone"
          dataKey="attendees"
          stroke="#10b981"
          strokeWidth={2}
          name="Asistentes"
          dot={{ fill: '#10b981', r: 4 }}
          activeDot={{ r: 6 }}
          style={{ cursor: onMonthClick ? 'pointer' : 'default' }}
        />
      </LineChart>
    </ResponsiveContainer>
  );
};
