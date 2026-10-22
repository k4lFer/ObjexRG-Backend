import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { MeResponse } from "./types"

/** GET /user/me */
export async function getProfile(): Promise<MeResponse> {
  try {
    const { data } = await httpClient.get<ApiResult<MeResponse>>("/user/me")
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
