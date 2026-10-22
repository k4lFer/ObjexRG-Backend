import { login as loginRequest } from "../../../core/api/auth/login"
import { register as registerRequest } from "../../../core/api/auth/register"
import type { AuthRepository } from "../domain/auth.repository"
import type {
  AuthSession,
  LoginCredentials,
  RegisterInput,
} from "../domain/models"

export class HttpAuthRepository implements AuthRepository {
  async login(credentials: LoginCredentials): Promise<AuthSession> {
    const payload = await loginRequest({
      email: credentials.email,
      password: credentials.password,
    })

    return {
      user: {
        id: payload.id,
        username: payload.username,
        email: payload.email,
        firstName: payload.first_name,
        lastName: payload.last_name,
      },
      accessToken: payload.access_token,
      refreshToken: payload.refresh_token,
      tokenType: payload.token_type,
      expiresIn: payload.expires_in,
    }
  }

  async register(input: RegisterInput): Promise<void> {
    await registerRequest({
      email: input.email,
      username: input.username,
      password: input.password,
      first_name: input.firstName,
      last_name: input.lastName,
      birth_date: input.birthDate,
    })
  }
}
