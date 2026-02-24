// User Types
export type UserRole = 'ADMIN' | 'COORDINATOR' | 'BRAND';

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  brand_id?: string;
  brand_name?: string;
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
export type InvitationStatus = 'PENDING' | 'ACCEPTED' | 'EXPIRED';

export interface Invitation {
  id: string;
  email: string;
  role: UserRole;
  brand_id?: string;
  brand_name?: string;
  token: string;
  expires_at: string;
  accepted: boolean;
  created_at: string;
  status: InvitationStatus;
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

// Analytics Types
export interface TotalMetrics {
  total_events: number;
  total_attendees: number;
  total_leads: number;
  total_prospects: number;
  average_attendees: number;
  average_rating: number;
}

export interface BrandMetrics {
  brand_id: string;
  brand_name: string;
  event_count: number;
  total_attendees: number;
  total_leads: number;
  average_rating: number;
}

export interface MonthlyMetrics {
  year: number;
  month: number;
  month_name: string;
  event_count: number;
  attendees: number;
}

export interface StateMetrics {
  state: string;
  event_count: number;
  attendees: number;
}

export interface VehicleMetrics {
  vehicle_id: string;
  model_name: string;
  brand_name: string;
  times_presented: number;
  total_quantity: number;
}

export interface YearComparison {
  year: number;
  event_count: number;
  total_attendees: number;
  average_rating: number;
}

export interface EventTypeMetrics {
  event_type: string;
  event_count: number;
  attendees: number;
}

export interface DealerMetrics {
  dealer: string;
  event_count: number;
  average_rating: number;
  total_attendees: number;
}

export interface ConversionMetrics {
  total_attendees: number;
  total_leads: number;
  total_prospects: number;
  lead_rate: number;
  prospect_rate: number;
}

export interface CityMetrics {
  state: string;
  city: string;
  event_count: number;
  attendees: number;
}

export interface VenueMetrics {
  venue: string;
  event_count: number;
  total_attendees: number;
  total_leads: number;
  total_prospects: number;
}

export interface DashboardAnalytics {
  totals: TotalMetrics;
  by_brand: BrandMetrics[];
  by_month: MonthlyMetrics[];
  by_state: StateMetrics[];
  top_vehicles: VehicleMetrics[];
  year_comparison: YearComparison[];
  by_event_type: EventTypeMetrics[];
  top_dealers: DealerMetrics[];
  conversion: ConversionMetrics;
  top_cities: CityMetrics[];
  by_venue: VenueMetrics[];
}

export interface AnalyticsFilters {
  brand_id?: string;
  year?: number;
  start_date?: string;
  end_date?: string;
}
