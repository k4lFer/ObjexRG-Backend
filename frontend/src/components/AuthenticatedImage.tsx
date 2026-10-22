import { useEffect, useState } from "react"

interface AuthenticatedImageProps {
  load: () => Promise<Blob>
  alt: string
  className?: string
}

/**
 * Renders a protected image by fetching it with auth headers and turning the
 * blob into an object URL (an <img src> cannot send an Authorization header).
 */
export function AuthenticatedImage({
  load,
  alt,
  className,
}: AuthenticatedImageProps) {
  const [url, setUrl] = useState<string | null>(null)

  useEffect(() => {
    let active = true
    let objectUrl: string | null = null

    load()
      .then((blob) => {
        if (!active) {
          return
        }
        objectUrl = URL.createObjectURL(blob)
        setUrl(objectUrl)
      })
      .catch(() => {
        /* keep placeholder */
      })

    return () => {
      active = false
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl)
      }
    }
  }, [load])

  if (!url) {
    return (
      <div className="flex h-full w-full items-center justify-center px-2 text-center font-mono text-[10px] uppercase tracking-wide text-muted">
        {alt}
      </div>
    )
  }

  return <img src={url} alt={alt} className={className} />
}
