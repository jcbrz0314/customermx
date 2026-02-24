import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { store } from './store';
import { theme } from './theme';
import { SnackbarProvider } from './contexts/SnackbarContext';
import { SnackbarContainer } from './components/Snackbar';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Layout } from './components/Layout';
import { Login } from './pages/Login';
import { Dashboard } from './pages/Dashboard';
import { Brands } from './pages/Brands';
import { BrandForm } from './pages/Brands/BrandForm';
import { Vehicles } from './pages/Vehicles';
import { VehicleForm } from './pages/Vehicles/VehicleForm';
import { Events } from './pages/Events';
import { EventForm } from './pages/Events/EventForm';
import { EventDetail } from './pages/Events/EventDetail';
import { Users } from './pages/Users';
import { UserForm } from './pages/Users/UserForm';
import { Invitations } from './pages/Invitations';
import { AcceptInvitation } from './pages/AcceptInvitation';
import { Analytics } from './pages/Analytics';

function App() {
  return (
    <Provider store={store}>
      <ThemeProvider theme={theme}>
        <SnackbarProvider>
          <CssBaseline />
          <SnackbarContainer />
          <BrowserRouter>
          <Routes>
            {/* Public routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/invite/accept" element={<AcceptInvitation />} />

            {/* Protected routes */}
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Navigate to="/dashboard" replace />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Dashboard />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/brands"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Brands />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/brands/new"
              element={
                <ProtectedRoute>
                  <Layout>
                    <BrandForm />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/brands/:id/edit"
              element={
                <ProtectedRoute>
                  <Layout>
                    <BrandForm />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Vehicles />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles/new"
              element={
                <ProtectedRoute>
                  <Layout>
                    <VehicleForm />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles/:id/edit"
              element={
                <ProtectedRoute>
                  <Layout>
                    <VehicleForm />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/events"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Events />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/events/new"
              element={
                <ProtectedRoute>
                  <Layout>
                    <EventForm />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/events/:id"
              element={
                <ProtectedRoute>
                  <Layout>
                    <EventDetail />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/events/:id/edit"
              element={
                <ProtectedRoute>
                  <Layout>
                    <EventForm />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/users"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Users />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/users/:id/edit"
              element={
                <ProtectedRoute>
                  <Layout>
                    <UserForm />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/invitations"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Invitations />
                  </Layout>
                </ProtectedRoute>
              }
            />
            <Route
              path="/analytics"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Analytics />
                  </Layout>
                </ProtectedRoute>
              }
            />

            {/* Fallback route */}
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
        </BrowserRouter>
        </SnackbarProvider>
      </ThemeProvider>
    </Provider>
  );
}

export default App;
