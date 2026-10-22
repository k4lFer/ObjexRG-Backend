import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { SetVisibilityRequest, VisibilityResponse } from "./types"

/** PATCH /bucket/{bucket}/objects/{filename} */
export async function setObjectVisibility(
  bucket: string,
  fileName: string,
  isPublic: boolean,
): Promise<VisibilityResponse> {
  try {
    const body: SetVisibilityRequest = { is_public: isPublic }
    const { data } = await httpClient.patch<ApiResult<VisibilityResponse>>(
      `/bucket/${encodeURIComponent(bucket)}/objects/${encodeURIComponent(fileName)}`,
      body,
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
