import { useCallback } from "react"
import type { LucideIcon } from "lucide-react"
import {
  Download,
  Eye,
  EyeOff,
  File,
  FileArchive,
  FileImage,
  FileText,
  FileType,
  Film,
  Music,
  Share2,
  Trash2,
} from "lucide-react"

import { AuthenticatedImage } from "../../../components"
import type { FileKind } from "../../../shared/utils/format"
import { fileKind, formatBytes, isImage } from "../../../shared/utils/format"
import type { DriveObject } from "../domain/models"

const kindIcon: Record<FileKind, LucideIcon> = {
  image: FileImage,
  video: Film,
  audio: Music,
  pdf: FileText,
  archive: FileArchive,
  text: FileType,
  file: File,
}

interface FileCardProps {
  object: DriveObject
  getObjectBlob: (fileName: string) => Promise<Blob>
  onDownload: () => void
  onShare: () => void
  onToggleVisibility: () => void
  onDelete: () => void
}

export function FileCard({
  object,
  getObjectBlob,
  onDownload,
  onShare,
  onToggleVisibility,
  onDelete,
}: FileCardProps) {
  const load = useCallback(async () => {
    const blob = await getObjectBlob(object.fileName)
    return new Blob([blob], {
      type: object.mimeType || "application/octet-stream",
    })
  }, [getObjectBlob, object.fileName, object.mimeType])

  const image = isImage(object.mimeType)
  const KindIcon = kindIcon[fileKind(object.mimeType, object.extension)]

  return (
    <div className="group flex flex-col overflow-hidden rounded-xl border border-line bg-surface transition-colors hover:border-verdigris">
      <div className="relative flex h-32 items-center justify-center border-b border-line bg-paper">
        {image ? (
          <AuthenticatedImage
            load={load}
            alt={object.fileName}
            className="h-full w-full object-cover"
          />
        ) : (
          <KindIcon className="h-10 w-10 text-muted" />
        )}
        {object.isPublic ? (
          <span className="absolute left-2 top-2 inline-flex items-center gap-1 rounded-full bg-copper-tint px-2 py-0.5 font-mono text-[10px] uppercase tracking-wide text-copper">
            <Eye className="h-3 w-3" />
            public
          </span>
        ) : null}
      </div>

      <div className="flex flex-1 flex-col gap-2 p-3">
        <span
          className="truncate font-mono text-xs font-medium text-ink"
          title={object.fileName}
        >
          {object.fileName}
        </span>
        <span className="font-mono text-[11px] text-muted">
          {formatBytes(object.sizeBytes)}
        </span>

        <div className="mt-auto flex flex-wrap gap-1 pt-1">
          <Action label="Download" icon={Download} onClick={onDownload} />
          <Action label="Share" icon={Share2} onClick={onShare} />
          <Action
            label={object.isPublic ? "Private" : "Public"}
            icon={object.isPublic ? EyeOff : Eye}
            onClick={onToggleVisibility}
          />
          <Action label="Delete" icon={Trash2} danger onClick={onDelete} />
        </div>
      </div>
    </div>
  )
}

function Action({
  label,
  icon: Icon,
  onClick,
  danger = false,
}: {
  label: string
  icon: LucideIcon
  onClick: () => void
  danger?: boolean
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`inline-flex items-center gap-1 rounded-md border px-2 py-1 font-mono text-[10px] uppercase tracking-wide transition-colors ${
        danger
          ? "border-brick/30 text-brick hover:bg-brick/10"
          : "border-line text-graphite hover:border-verdigris hover:text-verdigris"
      }`}
    >
      <Icon className="h-3.5 w-3.5" />
      {label}
    </button>
  )
}
