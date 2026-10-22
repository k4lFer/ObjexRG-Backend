import { useState } from "react"

import { Button, Modal } from "../../../components"
import type { ShareLink } from "../domain/models"

interface ShareDialogProps {
  link: ShareLink | null
  onClose: () => void
}

export function ShareDialog({ link, onClose }: ShareDialogProps) {
  const [copied, setCopied] = useState(false)

  async function handleCopy() {
    if (!link) {
      return
    }
    try {
      await navigator.clipboard.writeText(link.url)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 2000)
    } catch {
      /* clipboard may be unavailable */
    }
  }

  return (
    <Modal open={link !== null} title="Share link" onClose={onClose}>
      {link ? (
        <div className="flex flex-col gap-3">
          <p className="text-sm text-graphite">
            Anyone with this link can read the object.
          </p>
          <input
            readOnly
            value={link.url}
            className="rounded-lg border border-line bg-paper px-3 py-2 font-mono text-xs text-ink"
            onFocus={(event) => event.target.select()}
          />
          <div className="flex justify-end gap-2">
            <Button variant="ghost" onClick={onClose}>
              Close
            </Button>
            <Button onClick={handleCopy}>{copied ? "Copied" : "Copy link"}</Button>
          </div>
        </div>
      ) : null}
    </Modal>
  )
}
