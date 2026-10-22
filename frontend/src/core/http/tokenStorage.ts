/**
 * Where tokens live between requests.
 *
 * The access token stays in memory (fast, XSS-resistant enough for a SPA) while
 * the refresh token is persisted so the session survives a page reload.
 * Swap the persistence here if you move the refresh token to an httpOnly cookie.
 */
let accessToken: string | null = null

const REFRESH_KEY = "objexrg.refresh_token"

export const tokenStorage = {
  getAccessToken(): string | null {
    return accessToken
  },

  getRefreshToken(): string | null {
    return localStorage.getItem(REFRESH_KEY)
  },

  set(access: string, refresh: string): void {
    accessToken = access
    localStorage.setItem(REFRESH_KEY, refresh)
  },

  setAccess(access: string): void {
    accessToken = access
  },

  clear(): void {
    accessToken = null
    localStorage.removeItem(REFRESH_KEY)
  },
}
