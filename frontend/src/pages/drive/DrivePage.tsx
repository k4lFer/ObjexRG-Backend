import { useState } from "react"
import { Boxes, Loader2 } from "lucide-react"

import { ConfirmDialog } from "../../components"
import { useDrive } from "./application/useDrive"
import { BucketSidebar } from "./components/BucketSidebar"
import { FileGrid } from "./components/FileGrid"
import { ShareDialog } from "./components/ShareDialog"
import { UploadButton } from "./components/UploadButton"
import type { DriveObject, ShareLink } from "./domain/models"

export function DrivePage() {
  const drive = useDrive()
  const [shareLink, setShareLink] = useState<ShareLink | null>(null)
  const [pendingDelete, setPendingDelete] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  async function handleShare(object: DriveObject) {
    const link = await drive.createShare(object.fileName)
    if (link) {
      setShareLink(link)
    }
  }

  async function handleConfirmDelete() {
    if (!pendingDelete) {
      return
    }
    setDeleting(true)
    await drive.removeObject(pendingDelete)
    setDeleting(false)
    setPendingDelete(null)
  }

  return (
    <div className="flex flex-1">
      <BucketSidebar
        buckets={drive.buckets}
        selectedBucket={drive.selectedBucket}
        loading={drive.loadingBuckets}
        onSelect={drive.selectBucket}
        onCreate={drive.createBucket}
        onDelete={drive.removeBucket}
      />

      <section className="relative flex flex-1 flex-col overflow-hidden">
        <header className="flex items-center justify-between border-b border-line bg-surface px-6 py-3">
          <div>
            <h1 className="flex items-center gap-2 font-display text-lg font-semibold text-ink">
              <Boxes className="h-4 w-4 text-verdigris" />
              {drive.selectedBucket ?? "Select a bucket"}
            </h1>
            <p className="mt-0.5 font-mono text-[11px] uppercase tracking-widest text-muted">
              {drive.selectedBucket
                ? `${drive.objects.length} object(s)`
                : "choose a bucket to begin"}
            </p>
          </div>
          <UploadButton
            disabled={!drive.selectedBucket}
            uploading={drive.uploading}
            onUpload={drive.uploadFiles}
          />
        </header>

        <div className="flex flex-1 flex-col">
          {drive.selectedBucket ? (
            <FileGrid
              objects={drive.objects}
              loading={drive.loadingObjects}
              getObjectBlob={drive.getObjectBlob}
              onDownload={drive.download}
              onShare={handleShare}
              onToggleVisibility={drive.toggleVisibility}
              onDelete={setPendingDelete}
            />
          ) : (
            <div className="flex flex-1 items-center justify-center">
              <p className="font-mono text-xs uppercase tracking-widest text-muted">
                Select a bucket to view its objects
              </p>
            </div>
          )}
        </div>

        {drive.uploadProgress ? (
          <UploadOverlay
            done={drive.uploadProgress.done}
            total={drive.uploadProgress.total}
            loadedBytes={drive.uploadProgress.loadedBytes}
            totalBytes={drive.uploadProgress.totalBytes}
          />
        ) : null}
      </section>

      <ShareDialog link={shareLink} onClose={() => setShareLink(null)} />

      <ConfirmDialog
        open={pendingDelete !== null}
        title="Delete object"
        description={
          <>
            Delete <span className="font-mono">{pendingDelete}</span>? This cannot
            be undone.
          </>
        }
        loading={deleting}
        onConfirm={handleConfirmDelete}
        onCancel={() => setPendingDelete(null)}
      />
    </div>
  )
}

function UploadOverlay({
  done,
  total,
  loadedBytes,
  totalBytes,
}: {
  done: number
  total: number
  loadedBytes: number
  totalBytes: number
}) {
  const ratio =
    totalBytes > 0
      ? loadedBytes / totalBytes
      : total === 0
        ? 0
        : done / total
  const percent = Math.round(Math.min(ratio, 1) * 100)

  return (
    <div className="absolute inset-0 z-10 flex items-center justify-center bg-paper/70 backdrop-blur-sm">
      <div className="flex w-72 flex-col gap-4 rounded-2xl border border-line bg-surface p-6 shadow-2xl">
        <div className="flex items-center gap-3">
          <Loader2 className="h-5 w-5 animate-spin text-verdigris" />
          <div className="flex flex-col">
            <span className="font-display text-sm font-semibold text-ink">
              Uploading objects
            </span>
            <span className="font-mono text-[11px] uppercase tracking-widest text-muted">
              {done} / {total} complete
            </span>
          </div>
        </div>
        <div className="h-1.5 overflow-hidden rounded-full bg-line">
          <div
            className="h-full rounded-full bg-verdigris transition-[width] duration-300"
            style={{ width: `${percent}%` }}
          />
        </div>
        <p className="text-xs text-graphite">
          Keep this tab open until the upload finishes.
        </p>
      </div>
    </div>
  )
}
