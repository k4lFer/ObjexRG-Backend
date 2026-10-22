import { useCallback, useEffect, useMemo, useState } from "react"
import type { ReactNode } from "react"

import { logout as logoutRequest } from "../../../core/api/auth/logout"
import { setUnauthorizedHandler } from "../../../core/http/interceptors"
import { tokenStorage } from "../../../core/http/tokenStorage"
import { HttpAuthRepository } from "../data/httpAuthRepository"
import type { AuthRepository } from "../domain/auth.repository"
import type { LoginCredentials, RegisterInput, User } from "../domain/models"
import { AuthContext } from "./AuthContext"
import type { AuthContextValue, AuthStatus } from "./AuthContext"
import { HttpUserRepository } from "../../user/data/httpUserRepository"
import type { UserRepository } from "../../user/domain/user.repository"

const authRepository: AuthRepository = new HttpAuthRepository()
const userRepository: UserRepository = new HttpUserRepository()

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [status, setStatus] = useState<AuthStatus>(() =>
    tokenStorage.getRefreshToken() ? "loading" : "anonymous",
  )

  useEffect(() => {
    setUnauthorizedHandler(() => {
      setUser(null)
      setStatus("anonymous")
    })
  }, [])

  useEffect(() => {
    if (status !== "loading") {
      return
    }

    let active = true
    userRepository
      .getMe()
      .then((profile) => {
        if (!active) {
          return
        }
        setUser({
          id: profile.id,
          username: profile.username,
          email: profile.email,
          firstName: profile.firstName,
          lastName: profile.lastName,
        })
        setStatus("authenticated")
      })
      .catch(() => {
        if (!active) {
          return
        }
        tokenStorage.clear()
        setStatus("anonymous")
      })

    return () => {
      active = false
    }
  }, [status])

  const login = useCallback(async (credentials: LoginCredentials) => {
    const session = await authRepository.login(credentials)
    tokenStorage.set(session.accessToken, session.refreshToken)
    setUser(session.user)
    setStatus("authenticated")
  }, [])

  const register = useCallback(async (input: RegisterInput) => {
    await authRepository.register(input)
  }, [])

  const logout = useCallback(async () => {
    const refreshToken = tokenStorage.getRefreshToken()
    if (refreshToken) {
      try {
        await logoutRequest({ refresh_token: refreshToken })
      } catch {
        /* best-effort: clear the session regardless of the API result */
      }
    }
    tokenStorage.clear()
    setUser(null)
    setStatus("anonymous")
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({ user, status, login, register, logout }),
    [user, status, login, register, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
