import type { Bucket, DriveObject, ShareLink } from "./models"

/** Contract the drive (filesystem) application layer depends on. */
export interface DriveRepository {
  listBuckets(): Promise<Bucket[]>
  createBucket(name: string): Promise<void>
  deleteBucket(name: string): Promise<void>

  listObjects(bucket: string): Promise<DriveObject[]>
  uploadObject(
    bucket: string,
    file: File,
    onProgress?: (loadedBytes: number) => void,
  ): Promise<void>
  deleteObject(bucket: string, fileName: string): Promise<void>
  setVisibility(bucket: string, fileName: string, isPublic: boolean): Promise<void>
  /** Absolute presigned URL when the backend can mint one, otherwise null. */
  getObjectUrl(
    bucket: string,
    fileName: string,
    disposition: "inline" | "attachment",
  ): Promise<string | null>
  downloadObject(bucket: string, fileName: string): Promise<Blob>
  createShare(bucket: string, fileName: string): Promise<ShareLink>
}
