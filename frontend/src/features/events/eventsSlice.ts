import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { EventWithBrand, EventStatus } from '../../types';

interface EventsState {
  events: EventWithBrand[];
  selectedEvent: EventWithBrand | null;
  loading: boolean;
  error: string | null;
  filters: {
    brandId?: string;
    year?: number;
    status?: EventStatus;
    state?: string;
  };
}

const initialState: EventsState = {
  events: [],
  selectedEvent: null,
  loading: false,
  error: null,
  filters: {},
};

const eventsSlice = createSlice({
  name: 'events',
  initialState,
  reducers: {
    setEvents: (state, action: PayloadAction<EventWithBrand[]>) => {
      state.events = action.payload;
      state.loading = false;
      state.error = null;
    },
    setSelectedEvent: (state, action: PayloadAction<EventWithBrand | null>) => {
      state.selectedEvent = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
      if (action.payload) {
        state.error = null;
      }
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
      state.loading = false;
    },
    setFilters: (state, action: PayloadAction<EventsState['filters']>) => {
      state.filters = action.payload;
    },
    clearFilters: (state) => {
      state.filters = {};
    },
    clearError: (state) => {
      state.error = null;
    },
  },
});

export const {
  setEvents,
  setSelectedEvent,
  setLoading,
  setError,
  setFilters,
  clearFilters,
  clearError,
} = eventsSlice.actions;

export default eventsSlice.reducer;
