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
import { VehicleMetrics } from '../../types';

interface Props {
  data: VehicleMetrics[];
  onVehicleClick?: (vehicleId: string, modelName: string) => void;
}

export const TopVehiclesChart: React.FC<Props> = ({ data, onVehicleClick }) => {
  const handleClick = (data: any) => {
    if (onVehicleClick && data && data.activePayload && data.activePayload[0]) {
      const vehicle = data.activePayload[0].payload as VehicleMetrics;
      onVehicleClick(vehicle.vehicle_id, vehicle.model_name);
    }
  };

  return (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart data={data} layout="vertical" margin={{ left: 20 }} onClick={handleClick}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis type="number" tick={{ fontSize: 12 }} stroke="#6b7280" />
        <YAxis
          dataKey="model_name"
          type="category"
          width={120}
          tick={{ fontSize: 11 }}
          stroke="#6b7280"
        />
        <Tooltip
          contentStyle={{
            backgroundColor: '#fff',
            border: '1px solid #e5e7eb',
            borderRadius: '8px',
          }}
          cursor={{ fill: 'rgba(139, 92, 246, 0.1)' }}
        />
        <Bar
          dataKey="total_quantity"
          fill="#8b5cf6"
          name="Cantidad Total"
          radius={[0, 8, 8, 0]}
          style={{ cursor: onVehicleClick ? 'pointer' : 'default' }}
        />
      </BarChart>
    </ResponsiveContainer>
  );
};
