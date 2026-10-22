import { useCallback, useEffect, useState } from "react"
import { toast } from "sonner"

import { notifyError } from "../../../shared/utils/notify"
import { HttpDriveRepository } from "../data/httpDriveRepository"
import type { DriveRepository } from "../domain/drive.repository"
import type { Bucket, DriveObject, ShareLink } from "../domain/models"

const driveRepository: DriveRepository = new HttpDriveRepository()

const MIN_UPLOAD_FEEDBACK_MS = 1200

export interface UploadProgress {
  done: number
  total: number
  loadedBytes: number
  totalBytes: number
}

export interface DriveState {
  buckets: Bucket[]
  objects: DriveObject[]
  selectedBucket: string | null
  loadingBuckets: boolean
  loadingObjects: boolean
  uploading: boolean
  uploadProgress: UploadProgress | null
  selectBucket: (name: string) => void
  refreshObjects: () => void
  createBucket: (name: string) => Promise<void>
  removeBucket: (name: string) => Promise<void>
  uploadFiles: (files: File[]) => Promise<void>
  removeObject: (fileName: string) => Promise<void>
  toggleVisibility: (object: DriveObject) => Promise<void>
  createShare: (fileName: string) => Promise<ShareLink | null>
  download: (fileName: string) => Promise<void>
  getObjectBlob: (fileName: string) => Promise<Blob>
}

export function useDrive(): DriveState {
  const [buckets, setBuckets] = useState<Bucket[]>([])
  const [objects, setObjects] = useState<DriveObject[]>([])
  const [selectedBucket, setSelectedBucket] = useState<string | null>(null)
  const [loadingBuckets, setLoadingBuckets] = useState(true)
  const [loadingObjects, setLoadingObjects] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [uploadProgress, setUploadProgress] = useState<UploadProgress | null>(
    null,
  )

  const loadBuckets = useCallback(async () => {
    try {
      const list = await driveRepository.listBuckets()
      setBuckets(list)
    } catch (err) {
      notifyError(err)
    } finally {
      setLoadingBuckets(false)
    }
  }, [])

  const loadObjects = useCallback(async (bucket: string) => {
    setLoadingObjects(true)
    try {
      setObjects(await driveRepository.listObjects(bucket))
    } catch (err) {
      notifyError(err)
    } finally {
      setLoadingObjects(false)
    }
  }, [])

  useEffect(() => {
    let active = true
    driveRepository
      .listBuckets()
      .then((list) => {
        if (active) {
          setBuckets(list)
        }
      })
      .catch((err) => {
        if (active) {
          notifyError(err)
        }
      })
      .finally(() => {
        if (active) {
          setLoadingBuckets(false)
        }
      })

    return () => {
      active = false
    }
  }, [])

  const selectBucket = useCallback(
    (name: string) => {
      setSelectedBucket(name)
      setObjects([])
      loadObjects(name)
    },
    [loadObjects],
  )

  const refreshObjects = useCallback(() => {
    if (selectedBucket) {
      loadObjects(selectedBucket)
    }
  }, [selectedBucket, loadObjects])

  const createBucket = useCallback(
    async (name: string) => {
      try {
        await driveRepository.createBucket(name)
        await loadBuckets()
        selectBucket(name)
      } catch (err) {
        notifyError(err)
      }
    },
    [loadBuckets, selectBucket],
  )

  const removeBucket = useCallback(
    async (name: string) => {
      const previous = buckets
      setBuckets((current) => current.filter((bucket) => bucket.name !== name))
      if (selectedBucket === name) {
        setSelectedBucket(null)
        setObjects([])
      }
      try {
        await driveRepository.deleteBucket(name)
        toast.success(`Bucket "${name}" deleted`)
      } catch (err) {
        notifyError(err)
        setBuckets(previous)
      }
    },
    [buckets, selectedBucket],
  )

  const uploadFiles = useCallback(
    async (files: File[]) => {
      if (!selectedBucket) {
        return
      }
      const totalBytes = files.reduce((sum, file) => sum + file.size, 0)
      setUploading(true)
      setUploadProgress({
        done: 0,
        total: files.length,
        loadedBytes: 0,
        totalBytes,
      })
      const startedAt = Date.now()
      try {
        let doneBytes = 0
        for (let index = 0; index < files.length; index++) {
          const file = files[index]
          await driveRepository.uploadObject(selectedBucket, file, (loaded) => {
            setUploadProgress({
              done: index,
              total: files.length,
              loadedBytes: doneBytes + Math.min(loaded, file.size),
              totalBytes,
            })
          })
          doneBytes += file.size
          setUploadProgress({
            done: index + 1,
            total: files.length,
            loadedBytes: doneBytes,
            totalBytes,
          })
        }
        await loadObjects(selectedBucket)
        toast.success(
          files.length === 1 ? "File uploaded" : `${files.length} files uploaded`,
        )
      } catch (err) {
        notifyError(err)
      } finally {
        const elapsed = Date.now() - startedAt
        if (elapsed < MIN_UPLOAD_FEEDBACK_MS) {
          await new Promise((resolve) =>
            setTimeout(resolve, MIN_UPLOAD_FEEDBACK_MS - elapsed),
          )
        }
        setUploading(false)
        setUploadProgress(null)
      }
    },
    [selectedBucket, loadObjects],
  )

  const removeObject = useCallback(
    async (fileName: string) => {
      if (!selectedBucket) {
        return
      }
      setObjects((current) =>
        current.filter((object) => object.fileName !== fileName),
      )
      try {
        await driveRepository.deleteObject(selectedBucket, fileName)
        toast.success("File deleted")
      } catch (err) {
        notifyError(err)
        await loadObjects(selectedBucket)
      }
    },
    [selectedBucket, loadObjects],
  )

  const toggleVisibility = useCallback(
    async (object: DriveObject) => {
      if (!selectedBucket) {
        return
      }
      try {
        await driveRepository.setVisibility(
          selectedBucket,
          object.fileName,
          !object.isPublic,
        )
        await loadObjects(selectedBucket)
      } catch (err) {
        notifyError(err)
      }
    },
    [selectedBucket, loadObjects],
  )

  const createShare = useCallback(
    async (fileName: string): Promise<ShareLink | null> => {
      if (!selectedBucket) {
        return null
      }
      try {
        return await driveRepository.createShare(selectedBucket, fileName)
      } catch (err) {
        notifyError(err)
        return null
      }
    },
    [selectedBucket],
  )

  const download = useCallback(
    async (fileName: string) => {
      if (!selectedBucket) {
        return
      }
      try {
        const directUrl = await driveRepository.getObjectUrl(
          selectedBucket,
          fileName,
          "attachment",
        )
        if (directUrl) {
          const anchor = document.createElement("a")
          anchor.href = directUrl
          anchor.download = fileName
          anchor.rel = "noopener"
          document.body.appendChild(anchor)
          anchor.click()
          anchor.remove()
          return
        }

        const blob = await driveRepository.downloadObject(
          selectedBucket,
          fileName,
        )
        const url = URL.createObjectURL(blob)
        const anchor = document.createElement("a")
        anchor.href = url
        anchor.download = fileName
        document.body.appendChild(anchor)
        anchor.click()
        anchor.remove()
        URL.revokeObjectURL(url)
      } catch (err) {
        notifyError(err)
      }
    },
    [selectedBucket],
  )

  const getObjectBlob = useCallback(
    (fileName: string): Promise<Blob> => {
      if (!selectedBucket) {
        return Promise.reject(new Error("No bucket selected"))
      }
      return driveRepository.downloadObject(selectedBucket, fileName)
    },
    [selectedBucket],
  )

  return {
    buckets,
    objects,
    selectedBucket,
    loadingBuckets,
    loadingObjects,
    uploading,
    uploadProgress,
    selectBucket,
    refreshObjects,
    createBucket,
    removeBucket,
    uploadFiles,
    removeObject,
    toggleVisibility,
    createShare,
    download,
    getObjectBlob,
  }
}
