import { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from './useRedux';
import { setCredentials, logout as logoutAction } from '../features/auth/authSlice';
import { apiService, API_ENDPOINTS } from '../services/api';
import { LoginRequest, LoginResponse } from '../types';

export const useAuth = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { user, accessToken, isAuthenticated } = useAppSelector((state) => state.auth);

  const login = useCallback(
    async (credentials: LoginRequest) => {
      const response = await apiService.post<LoginResponse>(
        API_ENDPOINTS.AUTH.LOGIN,
        credentials
      );

      if (response.error) {
        throw new Error(response.error);
      }

      if (response.data) {
        dispatch(setCredentials(response.data));
        return response.data;
      }

      throw new Error('Login failed');
    },
    [dispatch]
  );

  const logout = useCallback(() => {
    // Call logout endpoint (optional, since JWT is stateless)
    apiService.post(API_ENDPOINTS.AUTH.LOGOUT, {}, accessToken || undefined);

    // Clear Redux state and localStorage
    dispatch(logoutAction());

    // Redirect to login
    navigate('/login');
  }, [dispatch, navigate, accessToken]);

  return {
    user,
    accessToken,
    isAuthenticated,
    login,
    logout,
  };
};
