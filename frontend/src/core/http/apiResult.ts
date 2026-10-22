/** Envelope returned by every backend endpoint (internal/shared/common/result.go). */
export interface ApiResult<T> {
  success: boolean
  data: T | null
  messages: string[]
  code: number
}
