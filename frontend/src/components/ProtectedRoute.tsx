import { Navigate, Outlet, useLocation } from "react-router-dom"

import { Spinner } from "./Spinner"
import { useAuth } from "../pages/auth/application/useAuth"

export function ProtectedRoute() {
  const { status } = useAuth()
  const location = useLocation()

  if (status === "loading") {
    return (
      <div className="flex min-h-svh items-center justify-center">
        <Spinner label="Restoring session..." />
      </div>
    )
  }

  if (status === "anonymous") {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return <Outlet />
}
