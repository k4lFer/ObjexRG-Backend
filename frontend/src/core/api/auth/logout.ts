import { toApiError } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { LogoutRequest } from "./types"

/** POST /auth/logout */
export async function logout(payload: LogoutRequest): Promise<void> {
  try {
    await httpClient.post<ApiResult<null>>("/auth/logout", payload)
  } catch (error) {
    throw toApiError(error)
  }
}
