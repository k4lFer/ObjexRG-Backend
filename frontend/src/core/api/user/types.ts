/** Raw payload for the `/user/me` endpoint (backend snake_case). */

export interface MeResponse {
  id: string
  email: string
  username: string
  first_name: string
  last_name: string
  birth_date: string
  avatar_url?: string | null
  created_at: string
}
