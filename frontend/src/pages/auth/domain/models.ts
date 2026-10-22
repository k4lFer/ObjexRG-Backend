export interface User {
  id: string
  username: string
  email: string
  firstName: string
  lastName: string
}

export interface AuthSession {
  user: User
  accessToken: string
  refreshToken: string
  tokenType: string
  expiresIn: number
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterInput {
  email: string
  username: string
  password: string
  firstName: string
  lastName: string
  birthDate: string
}
