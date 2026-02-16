import React from 'react';
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { StateMetrics } from '../../types';

interface Props {
  data: StateMetrics[];
  onStateClick?: (state: string) => void;
}

const COLORS = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899'];

export const StateDistributionChart: React.FC<Props> = ({ data, onStateClick }) => {
  const handleClick = (data: any) => {
    if (onStateClick && data && data.state) {
      onStateClick(data.state);
    }
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={data}
          dataKey="event_count"
          nameKey="state"
          cx="50%"
          cy="50%"
          outerRadius={80}
          label={({ state, event_count }) => `${state}: ${event_count}`}
          labelLine={true}
          onClick={handleClick}
          style={{ cursor: onStateClick ? 'pointer' : 'default' }}
        >
          {data.map((entry, index) => (
            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip
          contentStyle={{
            backgroundColor: '#fff',
            border: '1px solid #e5e7eb',
            borderRadius: '8px',
          }}
        />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
};
