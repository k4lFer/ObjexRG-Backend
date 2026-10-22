import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { PresignUploadRequest, PresignUploadResponse } from "./types"

/** POST /upload/{bucket}/presign (reserves a key and returns a direct URL) */
export async function presignUpload(
  bucket: string,
  payload: PresignUploadRequest,
): Promise<PresignUploadResponse> {
  try {
    const { data } = await httpClient.post<ApiResult<PresignUploadResponse>>(
      `/upload/${encodeURIComponent(bucket)}/presign`,
      payload,
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
