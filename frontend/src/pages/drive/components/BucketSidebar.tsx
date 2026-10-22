import { useState } from "react"
import type { FormEvent } from "react"
import { Folder, Plus, X } from "lucide-react"

import { Button, ConfirmDialog, Spinner, TextField } from "../../../components"
import type { Bucket } from "../domain/models"

interface BucketSidebarProps {
  buckets: Bucket[]
  selectedBucket: string | null
  loading: boolean
  onSelect: (name: string) => void
  onCreate: (name: string) => Promise<void>
  onDelete: (name: string) => Promise<void>
}

export function BucketSidebar({
  buckets,
  selectedBucket,
  loading,
  onSelect,
  onCreate,
  onDelete,
}: BucketSidebarProps) {
  const [name, setName] = useState("")
  const [creating, setCreating] = useState(false)
  const [pendingDelete, setPendingDelete] = useState<string | null>(null)
  const [deleting, setDeleting] = useState(false)

  async function handleCreate(event: FormEvent) {
    event.preventDefault()
    const trimmed = name.trim()
    if (!trimmed) {
      return
    }
    setCreating(true)
    await onCreate(trimmed)
    setCreating(false)
    setName("")
  }

  async function handleConfirmDelete() {
    if (!pendingDelete) {
      return
    }
    setDeleting(true)
    await onDelete(pendingDelete)
    setDeleting(false)
    setPendingDelete(null)
  }

  return (
    <aside className="flex w-64 shrink-0 flex-col gap-5 border-r border-line bg-surface p-4">
      <form className="flex flex-col gap-2" onSubmit={handleCreate}>
        <span className="font-mono text-[10px] uppercase tracking-widest text-muted">
          New bucket
        </span>
        <TextField
          label=""
          name="bucket"
          placeholder="e.g. photos"
          value={name}
          onChange={(event) => setName(event.target.value)}
        />
        <Button
          type="submit"
          loading={creating}
          disabled={!name.trim()}
          className="w-full"
        >
          <Plus className="h-4 w-4" />
          Create bucket
        </Button>
      </form>

      <div className="flex flex-1 flex-col gap-1 overflow-y-auto">
        <span className="px-1 font-mono text-[10px] uppercase tracking-widest text-muted">
          Buckets
        </span>
        {loading ? (
          <Spinner label="Loading buckets..." />
        ) : buckets.length === 0 ? (
          <p className="px-1 text-sm text-muted">No buckets yet.</p>
        ) : (
          buckets.map((bucket) => {
            const active = bucket.name === selectedBucket
            return (
              <div
                key={bucket.id}
                className={`group flex items-center justify-between rounded-lg px-2 py-1.5 ${
                  active ? "bg-verdigris-tint" : "hover:bg-paper"
                }`}
              >
                <button
                  type="button"
                  className={`flex flex-1 items-center gap-2 truncate text-left font-mono text-xs ${
                    active
                      ? "font-medium text-verdigris-dark"
                      : "text-graphite"
                  }`}
                  onClick={() => onSelect(bucket.name)}
                >
                  <Folder className="h-4 w-4 shrink-0 text-verdigris" />
                  <span className="truncate">{bucket.name}</span>
                </button>
                <button
                  type="button"
                  className="ml-1 hidden text-muted hover:text-brick group-hover:block"
                  title="Delete bucket"
                  onClick={() => setPendingDelete(bucket.name)}
                >
                  <X className="h-4 w-4" />
                </button>
              </div>
            )
          })
        )}
      </div>

      <ConfirmDialog
        open={pendingDelete !== null}
        title="Delete bucket"
        description={
          <>
            Delete <span className="font-mono">{pendingDelete}</span> and all of its
            files? This cannot be undone.
          </>
        }
        loading={deleting}
        onConfirm={handleConfirmDelete}
        onCancel={() => setPendingDelete(null)}
      />
    </aside>
  )
}
