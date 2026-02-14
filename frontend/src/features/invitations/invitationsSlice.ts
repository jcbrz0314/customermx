import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { Invitation, InvitationStatus } from '../../types';

export interface InvitationsState {
  invitations: Invitation[];
  selectedInvitation: Invitation | null;
  loading: boolean;
  error: string | null;
  filters: {
    status?: InvitationStatus;
    search?: string;
  };
}

const initialState: InvitationsState = {
  invitations: [],
  selectedInvitation: null,
  loading: false,
  error: null,
  filters: {},
};

const invitationsSlice = createSlice({
  name: 'invitations',
  initialState,
  reducers: {
    setInvitations: (state, action: PayloadAction<Invitation[]>) => {
      state.invitations = action.payload;
    },
    setSelectedInvitation: (state, action: PayloadAction<Invitation | null>) => {
      state.selectedInvitation = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    setFilters: (state, action: PayloadAction<InvitationsState['filters']>) => {
      state.filters = action.payload;
    },
    updateInvitationInList: (state, action: PayloadAction<Invitation>) => {
      const index = state.invitations.findIndex((i) => i.id === action.payload.id);
      if (index !== -1) {
        state.invitations[index] = action.payload;
      }
    },
    removeInvitationFromList: (state, action: PayloadAction<string>) => {
      state.invitations = state.invitations.filter((i) => i.id !== action.payload);
    },
    clearFilters: (state) => {
      state.filters = {};
    },
  },
});

export const {
  setInvitations,
  setSelectedInvitation,
  setLoading,
  setError,
  setFilters,
  updateInvitationInList,
  removeInvitationFromList,
  clearFilters,
} = invitationsSlice.actions;

export default invitationsSlice.reducer;
