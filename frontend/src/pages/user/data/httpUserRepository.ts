import { getProfile } from "../../../core/api/user/profile"
import type { UserProfile } from "../domain/models"
import type { UserRepository } from "../domain/user.repository"

export class HttpUserRepository implements UserRepository {
  async getMe(): Promise<UserProfile> {
    const payload = await getProfile()

    return {
      id: payload.id,
      email: payload.email,
      username: payload.username,
      firstName: payload.first_name,
      lastName: payload.last_name,
      birthDate: payload.birth_date,
      avatarUrl: payload.avatar_url ?? null,
      createdAt: payload.created_at,
    }
  }
}
