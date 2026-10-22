import { useState } from "react"
import type { FormEvent } from "react"
import { toast } from "sonner"

import { Button, TextField } from "../../../components"
import { notifyError } from "../../../shared/utils/notify"
import { useAuth } from "../application/useAuth"
import type { RegisterInput } from "../domain/models"

const EMPTY_FORM: RegisterInput = {
  email: "",
  username: "",
  password: "",
  firstName: "",
  lastName: "",
  birthDate: "",
}

export function RegisterForm({ onSwitch }: { onSwitch: () => void }) {
  const { register } = useAuth()

  const [form, setForm] = useState<RegisterInput>(EMPTY_FORM)
  const [loading, setLoading] = useState(false)

  function update<K extends keyof RegisterInput>(key: K, value: RegisterInput[K]) {
    setForm((current) => ({ ...current, [key]: value }))
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setLoading(true)
    try {
      await register(form)
      toast.success("Account created. Please sign in.")
      onSwitch()
    } catch (error) {
      notifyError(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
      <h1 className="font-display text-2xl font-semibold text-ink">
        Create account
      </h1>
      <TextField
        label="Email"
        name="email"
        type="email"
        required
        value={form.email}
        onChange={(e) => update("email", e.target.value)}
      />
      <TextField
        label="Username"
        name="username"
        required
        value={form.username}
        onChange={(e) => update("username", e.target.value)}
      />
      <TextField
        label="Password"
        name="password"
        type="password"
        required
        value={form.password}
        onChange={(e) => update("password", e.target.value)}
      />
      <div className="grid grid-cols-2 gap-3">
        <TextField
          label="First name"
          name="firstName"
          required
          value={form.firstName}
          onChange={(e) => update("firstName", e.target.value)}
        />
        <TextField
          label="Last name"
          name="lastName"
          required
          value={form.lastName}
          onChange={(e) => update("lastName", e.target.value)}
        />
      </div>
      <TextField
        label="Birth date"
        name="birthDate"
        type="date"
        required
        value={form.birthDate}
        onChange={(e) => update("birthDate", e.target.value)}
      />
      <Button type="submit" loading={loading}>
        Create account
      </Button>
      <p className="text-sm text-graphite">
        Already registered?{" "}
        <button
          type="button"
          className="font-medium text-verdigris hover:underline"
          onClick={onSwitch}
        >
          Sign in
        </button>
      </p>
    </form>
  )
}
