/** Raw payloads for object-level endpoints. */

export interface UploadResponse {
  bucket: string
  filename: string
  path: string
  size: number
}

/** POST /upload/{bucket}/presign */
export interface PresignUploadRequest {
  fileName: string
  sizeBytes: number
  mimeType: string
}

/**
 * `mode: "s3-post"` carries a presigned POST `uploadUrl` (+ `fields` to send as
 * form data, with the file appended last); it enforces the max size server-side.
 * `mode: "s3"` carries a presigned PUT `uploadUrl` (+ `headers` to replay).
 * `mode: "proxy"` means the driver cannot presign, so fall back to the
 * multipart endpoint `POST /upload/{bucket}`.
 */
export interface PresignUploadResponse {
  mode: "s3" | "s3-post" | "proxy"
  method?: string
  uploadUrl?: string
  storageKey: string
  key: string
  headers?: Record<string, string>
  fields?: Record<string, string>
}

/** POST /upload/{bucket}/finalize */
export interface FinalizeUploadRequest {
  fileName: string
  storageKey: string
}

export interface FinalizeUploadResponse {
  id: string
  bucket: string
  fileName: string
  sizeBytes: number
  mimeType: string
  storageKey: string
  key: string
  createdAt: string
}

/** GET /bucket/{bucket}/objects/{filename}/link */
export interface ObjectLinkResponse {
  url: string
  fileName: string
  disposition: string
}

export interface SetVisibilityRequest {
  is_public: boolean
}

export interface VisibilityResponse {
  id: string
  bucket: string
  fileName: string
  isPublic: boolean
  publicUrl?: string
}
