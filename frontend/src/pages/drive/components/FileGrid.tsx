import { Spinner } from "../../../components"
import type { DriveObject } from "../domain/models"
import { FileCard } from "./FileCard"

interface FileGridProps {
  objects: DriveObject[]
  loading: boolean
  getObjectBlob: (fileName: string) => Promise<Blob>
  onDownload: (fileName: string) => void
  onShare: (object: DriveObject) => void
  onToggleVisibility: (object: DriveObject) => void
  onDelete: (fileName: string) => void
}

export function FileGrid({
  objects,
  loading,
  getObjectBlob,
  onDownload,
  onShare,
  onToggleVisibility,
  onDelete,
}: FileGridProps) {
  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center p-6">
        <Spinner label="Loading objects..." />
      </div>
    )
  }

  if (objects.length === 0) {
    return (
      <div className="flex flex-1 flex-col items-center justify-center gap-1 text-center">
        <p className="font-display text-sm font-semibold text-ink">
          This bucket is empty
        </p>
        <p className="text-sm text-muted">Upload a file to get started.</p>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-2 gap-4 p-6 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
      {objects.map((object) => (
        <FileCard
          key={object.id || object.fileName}
          object={object}
          getObjectBlob={getObjectBlob}
          onDownload={() => onDownload(object.fileName)}
          onShare={() => onShare(object)}
          onToggleVisibility={() => onToggleVisibility(object)}
          onDelete={() => onDelete(object.fileName)}
        />
      ))}
    </div>
  )
}
