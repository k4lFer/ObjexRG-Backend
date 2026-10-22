import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { BucketListResponse } from "./types"

/** GET /bucket */
export async function listBuckets(): Promise<BucketListResponse> {
  try {
    const { data } = await httpClient.get<ApiResult<BucketListResponse>>("/bucket")
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
