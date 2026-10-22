import type { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from "axios"

import { refresh } from "../api/auth/refresh"
import { tokenStorage } from "./tokenStorage"

type RetriableConfig = InternalAxiosRequestConfig & { _retry?: boolean }

const AUTH_PATHS = ["/auth/login", "/auth/register", "/auth/refresh"]

/** In-flight refresh promise shared by every concurrent 401 (single-flight). */
let refreshPromise: Promise<string> | null = null
let onUnauthorized: (() => void) | null = null

/** Lets the app react (clear session, redirect) when refresh finally fails. */
export function setUnauthorizedHandler(handler: () => void): void {
  onUnauthorized = handler
}

function isAuthPath(url?: string): boolean {
  if (!url) {
    return false
  }
  return AUTH_PATHS.some((path) => url.includes(path))
}

async function requestNewAccessToken(): Promise<string> {
  const refreshToken = tokenStorage.getRefreshToken()
  if (!refreshToken) {
    throw new Error("No refresh token available")
  }

  const issued = await refresh({ refresh_token: refreshToken })
  tokenStorage.set(issued.access_token, issued.refresh_token)
  return issued.access_token
}

export function attachInterceptors(client: AxiosInstance): void {
  client.interceptors.request.use((config) => {
    const accessToken = tokenStorage.getAccessToken()
    if (accessToken) {
      config.headers.set("Authorization", `Bearer ${accessToken}`)
    }
    return config
  })

  client.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const original = error.config as RetriableConfig | undefined

      const shouldRefresh =
        error.response?.status === 401 &&
        original !== undefined &&
        !original._retry &&
        !isAuthPath(original.url)

      if (!shouldRefresh) {
        return Promise.reject(error)
      }

      original._retry = true

      try {
        refreshPromise ??= requestNewAccessToken().finally(() => {
          refreshPromise = null
        })
        const newToken = await refreshPromise

        original.headers.set("Authorization", `Bearer ${newToken}`)
        return client(original)
      } catch (refreshError) {
        tokenStorage.clear()
        onUnauthorized?.()
        return Promise.reject(refreshError)
      }
    },
  )
}
