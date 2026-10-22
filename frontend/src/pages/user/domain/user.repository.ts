import type { UserProfile } from "./models"

export interface UserRepository {
  getMe(): Promise<UserProfile>
}
