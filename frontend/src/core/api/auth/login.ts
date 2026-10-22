import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { LoginRequest, LoginResponse } from "./types"

/** POST /auth/login */
export async function login(payload: LoginRequest): Promise<LoginResponse> {
  try {
    const { data } = await httpClient.post<ApiResult<LoginResponse>>(
      "/auth/login",
      payload,
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
