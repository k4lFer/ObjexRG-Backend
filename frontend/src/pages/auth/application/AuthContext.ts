import { createContext } from "react"

import type {
  LoginCredentials,
  RegisterInput,
  User,
} from "../domain/models"

export type AuthStatus = "loading" | "authenticated" | "anonymous"

export interface AuthContextValue {
  user: User | null
  status: AuthStatus
  login: (credentials: LoginCredentials) => Promise<void>
  register: (input: RegisterInput) => Promise<void>
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthContextValue | null>(null)
