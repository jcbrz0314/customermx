// Base URL configuration
export const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// API Endpoints
export const API_ENDPOINTS = {
  // Auth endpoints
  AUTH: {
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    ME: '/auth/me',
  },

  // User endpoints
  USERS: {
    LIST: '/users',
    CREATE: '/users',
    BY_ID: (id: string) => `/users/${id}`,
    BY_ROLE: (role: string) => `/users/role/${role}`,
  },

  // Brand endpoints
  BRANDS: {
    LIST: '/brands',
    CREATE: '/brands',
    BY_ID: (id: string) => `/brands/${id}`,
    VEHICLES: (brandId: string) => `/brands/${brandId}/vehicles`,
  },

  // Vehicle endpoints
  VEHICLES: {
    LIST: '/vehicles',
    CREATE: '/vehicles',
    BY_ID: (id: string) => `/vehicles/${id}`,
  },

  // Invitation endpoints
  INVITATIONS: {
    LIST: '/invitations',
    CREATE: '/invitations',
    BY_ID: (id: string) => `/invitations/${id}`,
    ACCEPT: '/invitations/accept',
    RESEND: (id: string) => `/invitations/${id}/resend`,
  },

  // Event endpoints (Fase 2)
  EVENTS: {
    LIST: '/events',
    CREATE: '/events',
    BY_ID: (id: string) => `/events/${id}`,
    COORDINATORS: (eventId: string) => `/events/${eventId}/coordinators`,
    REPORT: (eventId: string) => `/events/${eventId}/report`,
    VEHICLES: (eventId: string) => `/events/${eventId}/vehicles`,
  },

  // Analytics endpoints (Fase 3)
  ANALYTICS: {
    EVENTS: '/analytics/events',
    DASHBOARD: '/analytics/dashboard',
  },

  // Notification endpoints (Fase 4)
  NOTIFICATIONS: {
    LIST: '/notifications',
    UNREAD_COUNT: '/notifications/unread-count',
    MARK_READ: (id: string) => `/notifications/${id}/read`,
    MARK_ALL_READ: '/notifications/read-all',
  },
} as const;

// HTTP Status Codes
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  INTERNAL_SERVER_ERROR: 500,
} as const;

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Error de conexión. Por favor verifica tu conexión a internet.',
  UNAUTHORIZED: 'No tienes autorización para realizar esta acción.',
  SESSION_EXPIRED: 'Tu sesión ha expirado. Por favor inicia sesión nuevamente.',
  SERVER_ERROR: 'Error del servidor. Por favor intenta nuevamente.',
  VALIDATION_ERROR: 'Por favor verifica los datos ingresados.',
  NOT_FOUND: 'El recurso solicitado no fue encontrado.',
  CONFLICT: 'Ya existe un registro con estos datos.',
} as const;

// User Roles
export const USER_ROLES = {
  ADMIN: 'ADMIN',
  COORDINATOR: 'COORDINATOR',
  BRAND: 'BRAND',
} as const;

export type UserRole = typeof USER_ROLES[keyof typeof USER_ROLES];
