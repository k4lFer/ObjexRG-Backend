import axios from "axios"

import { env } from "../../config/env"
import { ApiError } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import type { RefreshRequest, RefreshResponse } from "./types"

/**
 * Plain axios instance (not httpClient) so the refresh request never re-enters
 * the 401 interceptor and loops. Using `baseURL` also avoids malformed URLs
 * like `//auth/refresh` when the app is served from the same origin.
 */
const refreshClient = axios.create({
  baseURL: env.apiBaseUrl,
})

/**
 * POST /auth/refresh
 */
export async function refresh(payload: RefreshRequest): Promise<RefreshResponse> {
  const { data } = await refreshClient.post<ApiResult<RefreshResponse>>(
    "/auth/refresh",
    payload,
  )

  if (!data.success || data.data === null) {
    throw new ApiError(data.messages, data.code)
  }

  return data.data
}
