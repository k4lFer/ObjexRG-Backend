import { AtSign, CalendarDays, Cake, Mail } from "lucide-react"
import type { LucideIcon } from "lucide-react"

import type { UserProfile } from "../domain/models"

export function ProfileCard({ profile }: { profile: UserProfile }) {
  const fullName =
    [profile.firstName, profile.lastName].filter(Boolean).join(" ") ||
    profile.username

  return (
    <article className="overflow-hidden rounded-2xl border border-line bg-surface shadow-sm">
      <div className="bg-grid h-28 bg-verdigris-dark" />

      <div className="px-6 pb-6">
        <div className="-mt-12 flex items-end justify-between">
          <Avatar profile={profile} fullName={fullName} />
        </div>

        <div className="mt-4 flex flex-col gap-2">
          <h2 className="font-display text-xl font-semibold text-ink">
            {fullName}
          </h2>
          <div className="flex flex-wrap items-center gap-x-4 gap-y-1">
            <Info icon={AtSign} value={profile.username} />
            <Info icon={Mail} value={profile.email} />
          </div>
        </div>

        <div className="mt-6 grid grid-cols-1 gap-3 border-t border-line pt-5 sm:grid-cols-2">
          <Detail
            icon={Cake}
            label="Birth date"
            value={formatDate(profile.birthDate)}
          />
          <Detail
            icon={CalendarDays}
            label="Member since"
            value={formatDate(profile.createdAt)}
          />
        </div>
      </div>
    </article>
  )
}

function Avatar({
  profile,
  fullName,
}: {
  profile: UserProfile
  fullName: string
}) {
  return (
    <div className="flex h-24 w-24 items-center justify-center overflow-hidden rounded-2xl border-4 border-surface bg-verdigris-tint shadow-sm">
      {profile.avatarUrl ? (
        <img
          src={profile.avatarUrl}
          alt={fullName}
          className="h-full w-full object-cover"
        />
      ) : (
        <span className="font-display text-3xl font-semibold text-verdigris-dark">
          {initials(fullName)}
        </span>
      )}
    </div>
  )
}

function Info({ icon: Icon, value }: { icon: LucideIcon; value: string }) {
  return (
    <span className="inline-flex items-center gap-1.5 font-mono text-xs text-graphite">
      <Icon className="h-3.5 w-3.5 text-muted" />
      {value}
    </span>
  )
}

function Detail({
  icon: Icon,
  label,
  value,
}: {
  icon: LucideIcon
  label: string
  value: string
}) {
  return (
    <div className="flex items-center gap-3 rounded-xl border border-line bg-paper/60 px-4 py-3">
      <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-surface text-verdigris">
        <Icon className="h-4 w-4" />
      </span>
      <div className="flex min-w-0 flex-col">
        <span className="font-mono text-[10px] uppercase tracking-widest text-muted">
          {label}
        </span>
        <span className="truncate text-sm text-ink">{value}</span>
      </div>
    </div>
  )
}

function initials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean)
  if (parts.length === 0) {
    return "?"
  }
  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase()
  }
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
}

function formatDate(value: string): string {
  if (!value) {
    return "—"
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  })
}
