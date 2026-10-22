export interface Bucket {
  id: string
  name: string
  status: string
  createdAt: string
}

export interface DriveObject {
  id: string
  fileName: string
  sizeBytes: number
  extension: string
  mimeType: string
  isPublic: boolean
  createdAt: string
}

export interface ShareLink {
  token: string
  url: string
}
