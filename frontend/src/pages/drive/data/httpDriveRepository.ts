import { createBucket as createBucketRequest } from "../../../core/api/bucket/createBucket"
import { deleteBucket as deleteBucketRequest } from "../../../core/api/bucket/deleteBucket"
import { listBuckets as listBucketsRequest } from "../../../core/api/bucket/listBuckets"
import { listObjects as listObjectsRequest } from "../../../core/api/bucket/listObjects"
import type { StorageObjectDto } from "../../../core/api/bucket/types"
import { deleteObject as deleteObjectRequest } from "../../../core/api/object/deleteObject"
import { downloadObject as downloadObjectRequest } from "../../../core/api/object/downloadObject"
import { getObjectLink } from "../../../core/api/object/getObjectLink"
import { setObjectVisibility } from "../../../core/api/object/setObjectVisibility"
import { uploadObject as uploadObjectRequest } from "../../../core/api/object/uploadObject"
import { createShare as createShareRequest } from "../../../core/api/share/createShare"
import { env } from "../../../core/config/env"
import type { DriveRepository } from "../domain/drive.repository"
import type { Bucket, DriveObject, ShareLink } from "../domain/models"

function toObject(dto: StorageObjectDto): DriveObject {
  return {
    id: dto.id,
    fileName: dto.fileName,
    sizeBytes: dto.sizeBytes,
    extension: dto.extension,
    mimeType: dto.mimeType,
    isPublic: dto.isPublic,
    createdAt: dto.createdAt,
  }
}

function absoluteUrl(path: string): string {
  return path.startsWith("http") ? path : `${env.apiBaseUrl}${path}`
}

export class HttpDriveRepository implements DriveRepository {
  async listBuckets(): Promise<Bucket[]> {
    const { buckets } = await listBucketsRequest()
    return (buckets ?? []).map((bucket) => ({
      id: bucket.id,
      name: bucket.name,
      status: bucket.status,
      createdAt: bucket.createdAt,
    }))
  }

  async createBucket(name: string): Promise<void> {
    await createBucketRequest({ name })
  }

  async deleteBucket(name: string): Promise<void> {
    await deleteBucketRequest(name)
  }

  async listObjects(bucket: string): Promise<DriveObject[]> {
    const { objects } = await listObjectsRequest(bucket)
    return (objects ?? []).map(toObject)
  }

  async uploadObject(
    bucket: string,
    file: File,
    onProgress?: (loadedBytes: number) => void,
  ): Promise<void> {
    await uploadObjectRequest(bucket, file, onProgress)
  }

  async deleteObject(bucket: string, fileName: string): Promise<void> {
    await deleteObjectRequest(bucket, fileName)
  }

  async setVisibility(
    bucket: string,
    fileName: string,
    isPublic: boolean,
  ): Promise<void> {
    await setObjectVisibility(bucket, fileName, isPublic)
  }

  async getObjectUrl(
    bucket: string,
    fileName: string,
    disposition: "inline" | "attachment",
  ): Promise<string | null> {
    const { url } = await getObjectLink(bucket, fileName, disposition)
    return url.startsWith("http") ? url : null
  }

  async downloadObject(bucket: string, fileName: string): Promise<Blob> {
    return downloadObjectRequest(bucket, fileName)
  }

  async createShare(bucket: string, fileName: string): Promise<ShareLink> {
    const link = await createShareRequest(bucket, fileName)
    return { token: link.token, url: absoluteUrl(link.url) }
  }
}
