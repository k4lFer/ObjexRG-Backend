import { useState } from "react"
import type { FormEvent } from "react"
import { useNavigate } from "react-router-dom"

import { Button, TextField } from "../../../components"
import { notifyError } from "../../../shared/utils/notify"
import { useAuth } from "../application/useAuth"

export function LoginForm({ onSwitch }: { onSwitch: () => void }) {
  const { login } = useAuth()
  const navigate = useNavigate()

  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [loading, setLoading] = useState(false)

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setLoading(true)
    try {
      await login({ email, password })
      navigate("/drive", { replace: true })
    } catch (error) {
      notifyError(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
      <h1 className="font-display text-2xl font-semibold text-ink">Sign in</h1>
      <TextField
        label="Email"
        name="email"
        type="email"
        autoComplete="email"
        required
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <TextField
        label="Password"
        name="password"
        type="password"
        autoComplete="current-password"
        required
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <Button type="submit" loading={loading}>
        Sign in
      </Button>
      <p className="text-sm text-graphite">
        No account?{" "}
        <button
          type="button"
          className="font-medium text-verdigris hover:underline"
          onClick={onSwitch}
        >
          Create one
        </button>
      </p>
    </form>
  )
}
