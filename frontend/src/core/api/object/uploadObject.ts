import axios from "axios"

import { toApiError, unwrapResult } from "../../http/apiError"
import type { ApiResult } from "../../http/apiResult"
import { httpClient } from "../../http/axiosClient"
import { finalizeUpload } from "./finalizeUpload"
import { presignUpload } from "./presignUpload"
import type { FinalizeUploadResponse, UploadResponse } from "./types"

/** Called with the number of bytes uploaded so far (best effort). */
export type UploadProgressHandler = (loadedBytes: number) => void

/**
 * Bare axios instance for talking to the object store: the presigned URL is
 * absolute and must not inherit the app baseURL nor the auth interceptor.
 */
const directClient = axios.create()

function toUploadResponse(result: FinalizeUploadResponse): UploadResponse {
  return {
    bucket: result.bucket,
    filename: result.fileName,
    path: result.key,
    size: result.sizeBytes,
  }
}

/**
 * Uploads a file to a bucket. When the backend driver supports presigning
 * (S3/RustFS) the bytes go straight to the object store via a presigned PUT
 * followed by a finalize call; otherwise it falls back to the multipart proxy
 * endpoint `POST /upload/{bucket}`.
 */
export async function uploadObject(
  bucket: string,
  file: File,
  onProgress?: UploadProgressHandler,
): Promise<UploadResponse> {
  const mimeType = file.type || "application/octet-stream"

  try {
    const presign = await presignUpload(bucket, {
      fileName: file.name,
      sizeBytes: file.size,
      mimeType,
    })

    if (presign.mode === "s3-post" && presign.uploadUrl) {
      // Presigned POST: send every policy field first and the file last. Let the
      // browser set the multipart Content-Type (with boundary); do not override it.
      const form = new FormData()
      for (const [name, value] of Object.entries(presign.fields ?? {})) {
        form.append(name, value)
      }
      form.append("file", file)

      await directClient.post(presign.uploadUrl, form, {
        onUploadProgress: (event) => onProgress?.(event.loaded),
      })

      const finalized = await finalizeUpload(bucket, {
        fileName: file.name,
        storageKey: presign.storageKey,
      })
      onProgress?.(file.size)
      return toUploadResponse(finalized)
    }

    if (presign.mode === "s3" && presign.uploadUrl) {
      await directClient.put(presign.uploadUrl, file, {
        headers: { ...presign.headers, "Content-Type": mimeType },
        onUploadProgress: (event) => onProgress?.(event.loaded),
      })

      const finalized = await finalizeUpload(bucket, {
        fileName: file.name,
        storageKey: presign.storageKey,
      })
      onProgress?.(file.size)
      return toUploadResponse(finalized)
    }

    return await uploadViaProxy(bucket, file, onProgress)
  } catch (error) {
    throw toApiError(error)
  }
}

async function uploadViaProxy(
  bucket: string,
  file: File,
  onProgress?: UploadProgressHandler,
): Promise<UploadResponse> {
  const form = new FormData()
  form.append("file", file)

  const { data } = await httpClient.post<ApiResult<UploadResponse>>(
    `/upload/${encodeURIComponent(bucket)}`,
    form,
    { onUploadProgress: (event) => onProgress?.(event.loaded) },
  )
  return unwrapResult(data)
}
