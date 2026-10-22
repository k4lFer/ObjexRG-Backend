/** Raw payloads for the `/auth/*` endpoints (backend snake_case). */

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  id: string
  username: string
  first_name: string
  last_name: string
  email: string
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
  first_name: string
  last_name: string
  birth_date: string
}

export interface RefreshRequest {
  refresh_token: string
}

export interface RefreshResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface LogoutRequest {
  refresh_token: string
}
