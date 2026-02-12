import { API_BASE_URL, HTTP_STATUS, ERROR_MESSAGES } from './apiConstants';

export interface ApiRequestConfig {
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  endpoint: string;
  body?: any;
  headers?: Record<string, string>;
  token?: string;
  useCredentials?: boolean;
}

export interface ApiResponse<T = any> {
  data?: T;
  error?: string;
  message?: string;
  status: number;
}

class ApiService {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  /**
   * Build full URL
   */
  private buildURL(endpoint: string): string {
    return `${this.baseURL}${endpoint}`;
  }

  /**
   * Build headers with authentication if token is provided
   */
  private buildHeaders(config: ApiRequestConfig): HeadersInit {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...config.headers,
    };

    // Add authorization header if token is provided
    if (config.token) {
      headers['Authorization'] = `Bearer ${config.token}`;
    }

    return headers;
  }

  /**
   * Handle API response
   * Backend returns: { data: {...}, message?: "..." } for success
   * Backend returns: { error: "...", message?: "..." } for errors
   */
  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    const status = response.status;

    // Handle no content
    if (status === HTTP_STATUS.NO_CONTENT) {
      return { status };
    }

    // Try to parse JSON response
    let data: T | undefined;
    let error: string | undefined;
    let message: string | undefined;

    try {
      const json = await response.json();

      if (response.ok) {
        // Success response from backend: { data: {...}, message?: "..." }
        data = json.data || json;
        message = json.message;
      } else {
        // Error response from backend: { error: "...", message?: "..." }
        error = json.error || json.message || this.getErrorMessage(status);
        message = json.message;
      }
    } catch (e) {
      // If response is not JSON
      if (!response.ok) {
        error = this.getErrorMessage(status);
      }
    }

    return { data, error, message, status };
  }

  /**
   * Get error message based on status code
   */
  private getErrorMessage(status: number): string {
    switch (status) {
      case HTTP_STATUS.UNAUTHORIZED:
        return ERROR_MESSAGES.UNAUTHORIZED;
      case HTTP_STATUS.FORBIDDEN:
        return ERROR_MESSAGES.UNAUTHORIZED;
      case HTTP_STATUS.NOT_FOUND:
        return ERROR_MESSAGES.NOT_FOUND;
      case HTTP_STATUS.CONFLICT:
        return ERROR_MESSAGES.CONFLICT;
      case HTTP_STATUS.BAD_REQUEST:
        return ERROR_MESSAGES.VALIDATION_ERROR;
      default:
        return ERROR_MESSAGES.SERVER_ERROR;
    }
  }

  /**
   * Make HTTP request
   */
  async request<T>(config: ApiRequestConfig): Promise<ApiResponse<T>> {
    const url = this.buildURL(config.endpoint);
    const headers = this.buildHeaders(config);

    const options: RequestInit = {
      method: config.method,
      headers,
      credentials: config.useCredentials ? 'include' : 'same-origin',
    };

    // Add body for POST, PUT, PATCH requests
    if (config.body && ['POST', 'PUT', 'PATCH'].includes(config.method)) {
      options.body = JSON.stringify(config.body);
    }

    try {
      const response = await fetch(url, options);
      return await this.handleResponse<T>(response);
    } catch (error) {
      console.error('API Request Error:', error);
      return {
        error: ERROR_MESSAGES.NETWORK_ERROR,
        status: 0,
      };
    }
  }

  /**
   * GET request
   */
  async get<T = any>(
    endpoint: string,
    token?: string,
    headers?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    return this.request<T>({
      method: 'GET',
      endpoint,
      token,
      headers,
      useCredentials: true,
    });
  }

  /**
   * POST request
   */
  async post<T = any>(
    endpoint: string,
    body?: any,
    token?: string,
    headers?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    return this.request<T>({
      method: 'POST',
      endpoint,
      body,
      token,
      headers,
      useCredentials: true,
    });
  }

  /**
   * PUT request
   */
  async put<T = any>(
    endpoint: string,
    body?: any,
    token?: string,
    headers?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    return this.request<T>({
      method: 'PUT',
      endpoint,
      body,
      token,
      headers,
      useCredentials: true,
    });
  }

  /**
   * PATCH request
   */
  async patch<T = any>(
    endpoint: string,
    body?: any,
    token?: string,
    headers?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    return this.request<T>({
      method: 'PATCH',
      endpoint,
      body,
      token,
      headers,
      useCredentials: true,
    });
  }

  /**
   * DELETE request
   */
  async delete<T = any>(
    endpoint: string,
    token?: string,
    headers?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    return this.request<T>({
      method: 'DELETE',
      endpoint,
      token,
      headers,
      useCredentials: true,
    });
  }
}

// Export singleton instance
export const apiService = new ApiService();

// Export class for testing or custom instances
export default ApiService;
