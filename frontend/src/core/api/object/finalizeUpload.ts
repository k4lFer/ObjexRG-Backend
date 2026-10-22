import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { FinalizeUploadRequest, FinalizeUploadResponse } from "./types"

/** POST /upload/{bucket}/finalize (registers a directly-uploaded blob) */
export async function finalizeUpload(
  bucket: string,
  payload: FinalizeUploadRequest,
): Promise<FinalizeUploadResponse> {
  try {
    const { data } = await httpClient.post<ApiResult<FinalizeUploadResponse>>(
      `/upload/${encodeURIComponent(bucket)}/finalize`,
      payload,
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
