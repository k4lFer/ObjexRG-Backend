import { Loader2 } from "lucide-react"

interface SpinnerProps {
  label?: string
}

export function Spinner({ label = "Loading..." }: SpinnerProps) {
  return (
    <div className="flex items-center justify-center gap-2 text-sm text-graphite">
      <Loader2 className="h-4 w-4 animate-spin text-verdigris" />
      <span>{label}</span>
    </div>
  )
}
