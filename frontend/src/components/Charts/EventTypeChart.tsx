import React from 'react';
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { EventTypeMetrics } from '../../types';

interface Props {
  data: EventTypeMetrics[];
  onEventTypeClick?: (eventType: string) => void;
}

const COLORS = ['#2563eb', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899'];

export const EventTypeChart: React.FC<Props> = ({ data, onEventTypeClick }) => {
  const handleClick = (data: any) => {
    if (onEventTypeClick && data && data.event_type) {
      onEventTypeClick(data.event_type);
    }
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={data}
          dataKey="event_count"
          nameKey="event_type"
          cx="50%"
          cy="50%"
          outerRadius={80}
          label={({ event_type, event_count }) => `${event_type}: ${event_count}`}
          labelLine={true}
          onClick={handleClick}
          style={{ cursor: onEventTypeClick ? 'pointer' : 'default' }}
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
