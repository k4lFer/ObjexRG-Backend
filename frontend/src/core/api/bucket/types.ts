/** Raw payloads for the `/bucket` endpoints (backend camelCase). */

export interface BucketDto {
  id: string
  name: string
  status: string
  createdAt: string
}

export interface BucketListResponse {
  buckets: BucketDto[] | null
}

export interface CreateBucketRequest {
  name: string
}

export interface CreateBucketResponse {
  bucket: string
  user: string
}

export interface StorageObjectDto {
  id: string
  fileName: string
  sizeBytes: number
  extension: string
  mimeType: string
  isPublic: boolean
  createdAt: string
}

export interface ObjectsListResponse {
  bucket: string
  objects: StorageObjectDto[] | null
}
