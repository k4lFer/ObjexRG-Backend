import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { CreateShareRequest, ShareLinkResponse } from "./types"

/** POST /share/{bucket}/{filename} */
export async function createShare(
  bucket: string,
  fileName: string,
  expiresAt?: string,
): Promise<ShareLinkResponse> {
  try {
    const body: CreateShareRequest = expiresAt ? { expires_at: expiresAt } : {}
    const { data } = await httpClient.post<ApiResult<ShareLinkResponse>>(
      `/share/${encodeURIComponent(bucket)}/${encodeURIComponent(fileName)}`,
      body,
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
