import { NavLink, Outlet, useNavigate } from "react-router-dom"
import { Boxes, LogOut, User } from "lucide-react"

import { useAuth } from "../pages/auth/application/useAuth"
import { Button } from "./Button"

const linkClass = ({ isActive }: { isActive: boolean }) =>
  `inline-flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
    isActive
      ? "bg-verdigris text-white"
      : "text-graphite hover:bg-paper hover:text-ink"
  }`

export function AppLayout() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  async function handleLogout() {
    await logout()
    navigate("/login", { replace: true })
  }

  return (
    <div className="flex min-h-svh flex-col bg-paper">
      <header className="flex items-center justify-between border-b border-line bg-surface px-6 py-3">
        <div className="flex items-center gap-6">
          <NavLink to="/" className="flex items-center gap-2">
            <span className="flex h-7 w-7 items-center justify-center rounded-md bg-verdigris text-white">
              <Boxes className="h-4 w-4" />
            </span>
            <span className="font-display text-lg font-semibold tracking-tight text-ink">
              ObjexRG
            </span>
          </NavLink>
          <nav className="flex items-center gap-1">
            <NavLink to="/drive" className={linkClass}>
              <Boxes className="h-4 w-4" />
              Drive
            </NavLink>
            <NavLink to="/profile" className={linkClass}>
              <User className="h-4 w-4" />
              Profile
            </NavLink>
          </nav>
        </div>
        <div className="flex items-center gap-3">
          <span className="font-mono text-xs text-muted">{user?.username}</span>
          <Button variant="ghost" onClick={handleLogout}>
            <LogOut className="h-4 w-4" />
            Sign out
          </Button>
        </div>
      </header>
      <main className="flex flex-1 flex-col">
        <Outlet />
      </main>
    </div>
  )
}
