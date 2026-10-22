import { ApiError, toApiError } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { RegisterRequest } from "./types"

/** POST /auth/register */
export async function register(payload: RegisterRequest): Promise<void> {
  try {
    const { data } = await httpClient.post<ApiResult<null>>(
      "/auth/register",
      payload,
    )
    if (!data.success) {
      throw new ApiError(data.messages, data.code)
    }
  } catch (error) {
    throw toApiError(error)
  }
}
