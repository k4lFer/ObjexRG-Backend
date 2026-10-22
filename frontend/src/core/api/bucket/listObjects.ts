import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { ObjectsListResponse } from "./types"

/** GET /bucket/{bucket} */
export async function listObjects(bucket: string): Promise<ObjectsListResponse> {
  try {
    const { data } = await httpClient.get<ApiResult<ObjectsListResponse>>(
      `/bucket/${encodeURIComponent(bucket)}`,
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
