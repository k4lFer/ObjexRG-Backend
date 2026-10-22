import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"

/** DELETE /bucket/{bucket}/objects/{filename} */
export async function deleteObject(
  bucket: string,
  fileName: string,
): Promise<void> {
  try {
    const { data } = await httpClient.delete<ApiResult<null>>(
      `/bucket/${encodeURIComponent(bucket)}/objects/${encodeURIComponent(fileName)}`,
    )
    unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
