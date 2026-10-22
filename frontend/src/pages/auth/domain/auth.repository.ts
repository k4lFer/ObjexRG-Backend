import type { AuthSession, LoginCredentials, RegisterInput } from "./models"

/** Contract the application layer depends on. Data layer implements it. */
export interface AuthRepository {
  login(credentials: LoginCredentials): Promise<AuthSession>
  register(input: RegisterInput): Promise<void>
}
