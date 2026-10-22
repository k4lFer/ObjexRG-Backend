import { useState } from "react"
import { Link } from "react-router-dom"
import { ArrowLeft, Boxes } from "lucide-react"

import { LoginForm } from "./components/LoginForm"
import { RegisterForm } from "./components/RegisterForm"

type Mode = "login" | "register"

export function AuthPage() {
  const [mode, setMode] = useState<Mode>("login")

  return (
    <div className="bg-grid flex min-h-svh flex-col items-center justify-center gap-6 bg-paper px-4 py-10">
      <Link to="/" className="flex items-center gap-2">
        <span className="flex h-8 w-8 items-center justify-center rounded-md bg-verdigris text-white">
          <Boxes className="h-4 w-4" />
        </span>
        <span className="font-display text-lg font-semibold tracking-tight text-ink">
          ObjexRG
        </span>
      </Link>

      <div className="w-full max-w-sm rounded-2xl border border-line bg-surface p-6 shadow-xl">
        {mode === "login" ? (
          <LoginForm onSwitch={() => setMode("register")} />
        ) : (
          <RegisterForm onSwitch={() => setMode("login")} />
        )}
      </div>

      <Link
        to="/"
        className="inline-flex items-center gap-1 font-mono text-[11px] uppercase tracking-widest text-muted hover:text-ink"
      >
        <ArrowLeft className="h-3.5 w-3.5" />
        Back to home
      </Link>
    </div>
  )
}
