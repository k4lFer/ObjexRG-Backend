import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import type { ObjectLinkResponse } from "./types"

export type LinkDisposition = "inline" | "attachment"

/**
 * GET /bucket/{bucket}/objects/{filename}/link
 *
 * Returns a URL for the object: an absolute presigned URL under the S3 driver,
 * or a relative backend path under the local driver.
 */
export async function getObjectLink(
  bucket: string,
  fileName: string,
  disposition: LinkDisposition = "inline",
): Promise<ObjectLinkResponse> {
  try {
    const { data } = await httpClient.get<ApiResult<ObjectLinkResponse>>(
      `/bucket/${encodeURIComponent(bucket)}/objects/${encodeURIComponent(fileName)}/link`,
      { params: { disposition } },
    )
    return unwrapResult(data)
  } catch (error) {
    throw toApiError(error)
  }
}
