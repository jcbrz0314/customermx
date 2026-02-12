import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { User } from '../../types';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
}

// Load initial state from localStorage
const loadAuthFromStorage = (): AuthState => {
  try {
    const user = localStorage.getItem('user');
    const accessToken = localStorage.getItem('accessToken');
    const refreshToken = localStorage.getItem('refreshToken');

    if (user && accessToken) {
      return {
        user: JSON.parse(user),
        accessToken,
        refreshToken,
        isAuthenticated: true,
      };
    }
  } catch (error) {
    console.error('Error loading auth from storage:', error);
  }

  return {
    user: null,
    accessToken: null,
    refreshToken: null,
    isAuthenticated: false,
  };
};

const initialState: AuthState = loadAuthFromStorage();

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials: (
      state,
      action: PayloadAction<{
        user: User;
        access_token: string;
        refresh_token: string;
      }>
    ) => {
      state.user = action.payload.user;
      state.accessToken = action.payload.access_token;
      state.refreshToken = action.payload.refresh_token;
      state.isAuthenticated = true;

      // Persist to localStorage
      localStorage.setItem('user', JSON.stringify(action.payload.user));
      localStorage.setItem('accessToken', action.payload.access_token);
      localStorage.setItem('refreshToken', action.payload.refresh_token);
    },

    updateTokens: (
      state,
      action: PayloadAction<{
        access_token: string;
        refresh_token: string;
      }>
    ) => {
      state.accessToken = action.payload.access_token;
      state.refreshToken = action.payload.refresh_token;

      // Persist to localStorage
      localStorage.setItem('accessToken', action.payload.access_token);
      localStorage.setItem('refreshToken', action.payload.refresh_token);
    },

    updateUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload;

      // Persist to localStorage
      localStorage.setItem('user', JSON.stringify(action.payload));
    },

    logout: (state) => {
      state.user = null;
      state.accessToken = null;
      state.refreshToken = null;
      state.isAuthenticated = false;

      // Clear localStorage
      localStorage.removeItem('user');
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
    },
  },
});

export const { setCredentials, updateTokens, updateUser, logout } = authSlice.actions;
export default authSlice.reducer;
