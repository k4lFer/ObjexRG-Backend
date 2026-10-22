export function formatBytes(bytes: number): string {
  if (bytes <= 0) {
    return "0 B"
  }
  const units = ["B", "KB", "MB", "GB", "TB"]
  const exponent = Math.min(
    Math.floor(Math.log(bytes) / Math.log(1024)),
    units.length - 1,
  )
  const value = bytes / 1024 ** exponent
  return `${value.toFixed(value >= 10 || exponent === 0 ? 0 : 1)} ${units[exponent]}`
}

export type FileKind = "image" | "video" | "audio" | "pdf" | "archive" | "text" | "file"

export function fileKind(mimeType: string, extension: string): FileKind {
  const mime = mimeType.toLowerCase()
  const ext = extension.toLowerCase()

  if (mime.startsWith("image/")) return "image"
  if (mime.startsWith("video/")) return "video"
  if (mime.startsWith("audio/")) return "audio"
  if (mime === "application/pdf" || ext === ".pdf") return "pdf"
  if (mime.startsWith("text/") || [".md", ".json", ".csv", ".log"].includes(ext)) {
    return "text"
  }
  if ([".zip", ".rar", ".7z", ".tar", ".gz"].includes(ext)) return "archive"
  return "file"
}

export function isImage(mimeType: string): boolean {
  return mimeType.toLowerCase().startsWith("image/")
}
