// User Types
export type UserRole = 'ADMIN' | 'COORDINATOR' | 'BRAND';

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  brand_id?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  password: string;
  role: UserRole;
  brand_id?: string;
}

export interface UpdateUserRequest {
  name?: string;
  email?: string;
  is_active?: boolean;
}

// Brand Types
export interface Brand {
  id: string;
  name: string;
  created_at: string;
}

export interface CreateBrandRequest {
  name: string;
}

export interface UpdateBrandRequest {
  name: string;
}

// Vehicle Types
export interface Vehicle {
  id: string;
  brand_id: string;
  model_name: string;
  created_at: string;
}

export interface VehicleWithBrand {
  id: string;
  brand_id: string;
  brand_name: string;
  model_name: string;
  created_at: string;
}

export interface CreateVehicleRequest {
  brand_id: string;
  model_name: string;
}

export interface UpdateVehicleRequest {
  model_name: string;
}

// Invitation Types
export interface Invitation {
  id: string;
  email: string;
  role: UserRole;
  brand_id?: string;
  token: string;
  expires_at: string;
  accepted: boolean;
  created_at: string;
}

export interface CreateInvitationRequest {
  email: string;
  role: UserRole;
  brand_id?: string;
}

export interface AcceptInvitationRequest {
  token: string;
  name: string;
  password: string;
}

// Event Types (Fase 2)
export type EventStatus = 'PLANNED' | 'ACTIVE' | 'COMPLETED' | 'CLOSED';

export interface Event {
  id: string;
  brand_id: string;
  event_type: string;
  organizer: string;
  name: string;
  start_date: string;
  year: number;
  duration_days: number;
  state: string;
  city: string;
  venue: string;
  dealer: string;
  status: EventStatus;
  created_by?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateEventRequest {
  brand_id: string;
  event_type: string;
  organizer: string;
  name: string;
  start_date: string;
  year: number;
  duration_days: number;
  state: string;
  city: string;
  venue: string;
  dealer: string;
}

export interface EventWithBrand extends Event {
  brand_name: string;
}

export interface UpdateEventRequest {
  event_type: string;
  organizer: string;
  name: string;
  start_date: string;
  year: number;
  duration_days: number;
  state: string;
  city: string;
  venue: string;
  dealer: string;
}

export interface ChangeStatusRequest {
  status: EventStatus;
}

// Event Coordinator Types
export interface EventCoordinator {
  id: string;
  event_id: string;
  user_id: string;
  assigned_at: string;
}

export interface EventCoordinatorWithUser extends EventCoordinator {
  user_name: string;
  user_email: string;
}

export interface AssignCoordinatorRequest {
  event_id: string;
  user_id: string;
}

// Event Vehicle Types
export interface EventVehicle {
  id: string;
  event_id: string;
  vehicle_id: string;
  quantity: number;
  created_at: string;
}

export interface EventVehicleWithDetails extends EventVehicle {
  model_name: string;
  brand_id: string;
  brand_name: string;
}

export interface AddVehicleRequest {
  event_id: string;
  vehicle_id: string;
  quantity: number;
}

export interface UpdateVehicleQuantityRequest {
  quantity: number;
}

// Event Report Types
export interface EventReport {
  id: string;
  event_id: string;
  hostess_count?: number;
  setup_vendor?: string;
  has_promotional?: boolean;
  attendees?: number;
  activities_count?: number;
  leads_collected?: number;
  prospects?: number;
  dealer_rating?: number;
  comments?: string;
  completed: boolean;
  updated_at: string;
}

// Notification Types
export interface Notification {
  id: string;
  user_id?: string;
  type: string;
  payload?: any;
  read: boolean;
  sent_at: string;
  created_at: string;
}
