import type { InputHTMLAttributes } from "react"

interface TextFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string
  error?: string
}

export function TextField({
  label,
  error,
  id,
  className = "",
  ...props
}: TextFieldProps) {
  const inputId = id ?? props.name

  return (
    <label className="flex flex-col gap-1.5 text-left" htmlFor={inputId}>
      <span className="text-xs font-medium uppercase tracking-wide text-graphite">
        {label}
      </span>
      <input
        id={inputId}
        className={`rounded-lg border bg-surface px-3 py-2 text-sm text-ink outline-none transition-colors placeholder:text-muted focus:border-verdigris focus:ring-2 focus:ring-verdigris/20 ${
          error ? "border-brick" : "border-line"
        } ${className}`}
        {...props}
      />
      {error ? <span className="text-xs text-brick">{error}</span> : null}
    </label>
  )
}
