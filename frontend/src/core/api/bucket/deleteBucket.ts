import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"

/** DELETE /bucket/{bucket} */
export async function deleteBucket(bucket: string): Promise<void> {
  try {
    const { data } = await httpClient.delete<ApiResult<null>>(
      `/bucket/${encodeURIComponent(bucket)}`,
    )
    unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
