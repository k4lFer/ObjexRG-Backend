import type { ChangeEvent } from "react"
import { useRef } from "react"
import { Upload } from "lucide-react"

import { Button } from "../../../components"

interface UploadButtonProps {
  disabled: boolean
  uploading: boolean
  onUpload: (files: File[]) => void
}

export function UploadButton({
  disabled,
  uploading,
  onUpload,
}: UploadButtonProps) {
  const inputRef = useRef<HTMLInputElement>(null)

  function handleChange(event: ChangeEvent<HTMLInputElement>) {
    const files = Array.from(event.target.files ?? [])
    if (files.length > 0) {
      onUpload(files)
    }
    event.target.value = ""
  }

  return (
    <>
      <input
        ref={inputRef}
        type="file"
        multiple
        className="hidden"
        onChange={handleChange}
      />
      <Button
        disabled={disabled}
        loading={uploading}
        onClick={() => inputRef.current?.click()}
      >
        <Upload className="h-4 w-4" />
        Upload files
      </Button>
    </>
  )
}
