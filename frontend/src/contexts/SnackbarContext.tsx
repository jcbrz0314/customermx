import { createContext, useState, useCallback, ReactNode } from 'react';

export type SnackbarSeverity = 'success' | 'error' | 'warning' | 'info';

export interface SnackbarMessage {
  id: string;
  message: string;
  severity: SnackbarSeverity;
}

interface SnackbarContextType {
  messages: SnackbarMessage[];
  showSnackbar: (message: string, severity: SnackbarSeverity) => void;
  hideSnackbar: (id: string) => void;
}

export const SnackbarContext = createContext<SnackbarContextType>({
  messages: [],
  showSnackbar: () => {},
  hideSnackbar: () => {},
});

interface SnackbarProviderProps {
  children: ReactNode;
}

export const SnackbarProvider = ({ children }: SnackbarProviderProps) => {
  const [messages, setMessages] = useState<SnackbarMessage[]>([]);

  const showSnackbar = useCallback((message: string, severity: SnackbarSeverity) => {
    const id = Date.now().toString();
    setMessages((prev) => [...prev, { id, message, severity }]);

    // Auto-hide after 5 seconds
    setTimeout(() => {
      hideSnackbar(id);
    }, 5000);
  }, []);

  const hideSnackbar = useCallback((id: string) => {
    setMessages((prev) => prev.filter((msg) => msg.id !== id));
  }, []);

  return (
    <SnackbarContext.Provider value={{ messages, showSnackbar, hideSnackbar }}>
      {children}
    </SnackbarContext.Provider>
  );
};
