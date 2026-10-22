import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { ShareListResponse } from "./types"

/** GET /share/list */
export async function listShares(): Promise<ShareListResponse> {
  try {
    const { data } = await httpClient.get<ApiResult<ShareListResponse>>(
      "/share/list",
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
