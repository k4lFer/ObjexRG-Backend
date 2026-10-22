import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { CreateBucketRequest, CreateBucketResponse } from "./types"

/** POST /bucket */
export async function createBucket(
  payload: CreateBucketRequest,
): Promise<CreateBucketResponse> {
  try {
    const { data } = await httpClient.post<ApiResult<CreateBucketResponse>>(
      "/bucket",
      payload,
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
