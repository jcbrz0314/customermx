import React from 'react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { DealerMetrics } from '../../types';

interface Props {
  data: DealerMetrics[];
  onDealerClick?: (dealer: string) => void;
}

export const DealerRankingChart: React.FC<Props> = ({ data, onDealerClick }) => {
  const handleClick = (data: any) => {
    if (onDealerClick && data && data.activePayload && data.activePayload[0]) {
      const dealerData = data.activePayload[0].payload as DealerMetrics;
      onDealerClick(dealerData.dealer);
    }
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={data} layout="vertical" margin={{ left: 20 }} onClick={handleClick}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis type="number" tick={{ fontSize: 12 }} stroke="#6b7280" />
        <YAxis
          dataKey="dealer"
          type="category"
          width={150}
          tick={{ fontSize: 11 }}
          stroke="#6b7280"
        />
        <Tooltip
          contentStyle={{
            backgroundColor: '#fff',
            border: '1px solid #e5e7eb',
            borderRadius: '8px',
          }}
          cursor={{ fill: 'rgba(245, 158, 11, 0.1)' }}
        />
        <Bar
          dataKey="average_rating"
          fill="#f59e0b"
          name="Rating Promedio"
          radius={[0, 8, 8, 0]}
          style={{ cursor: onDealerClick ? 'pointer' : 'default' }}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};
