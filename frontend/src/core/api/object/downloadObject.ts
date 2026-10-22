import { toApiError } from "../../http/apiError"
import { httpClient } from "../../http/axiosClient"

/** GET /download/{bucket}/{filename} (returns the raw file as a Blob) */
export async function downloadObject(
  bucket: string,
  fileName: string,
): Promise<Blob> {
  try {
    const { data } = await httpClient.get<Blob>(
      `/download/${encodeURIComponent(bucket)}/${encodeURIComponent(fileName)}`,
      { responseType: "blob" },
    )
    return data
  } catch (error) {
    throw toApiError(error)
  }
}
