/** Raw payloads for the `/share` endpoints. */

export interface CreateShareRequest {
  expires_at?: string
}

export interface ShareLinkResponse {
  token: string
  url: string
}

export interface ShareItemDto {
  token: string
  url: string
  bucket_name: string
  file_name: string
  mime_type: string
  size_bytes: number
  expires_at: string | null
  revoked: boolean
  download_count: number
  created_at: string
}

export interface ShareListResponse {
  details: ShareItemDto[]
}
