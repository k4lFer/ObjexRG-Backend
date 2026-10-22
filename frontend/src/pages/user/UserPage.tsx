import { useEffect, useState } from "react"

import { Spinner } from "../../components"
import { notifyError } from "../../shared/utils/notify"
import { ProfileCard } from "./components/ProfileCard"
import { HttpUserRepository } from "./data/httpUserRepository"
import type { UserProfile } from "./domain/models"

const userRepository = new HttpUserRepository()

export function UserPage() {
  const [profile, setProfile] = useState<UserProfile | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let active = true
    userRepository
      .getMe()
      .then((data) => {
        if (active) {
          setProfile(data)
        }
      })
      .catch((error) => {
        if (active) {
          notifyError(error)
        }
      })
      .finally(() => {
        if (active) {
          setLoading(false)
        }
      })

    return () => {
      active = false
    }
  }, [])

  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col gap-6 px-6 py-10">
      <div className="flex flex-col gap-1">
        <span className="font-mono text-[10px] uppercase tracking-widest text-copper">
          Account
        </span>
        <h1 className="font-display text-2xl font-semibold text-ink">
          Profile
        </h1>
      </div>

      {loading ? (
        <Spinner />
      ) : profile ? (
        <ProfileCard profile={profile} />
      ) : null}
    </div>
  )
}
