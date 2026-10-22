import { isAxiosError } from "axios"

import type { ApiResult } from "./apiResult"

/** Throws an ApiError when the envelope reports failure, otherwise returns data. */
export function unwrapResult<T>(result: ApiResult<T>): T {
  if (!result.success || result.data === null) {
    throw new ApiError(result.messages, result.code)
  }
  return result.data
}

/** Normalized error thrown by the repository layer, independent of axios. */
export class ApiError extends Error {
  readonly status: number
  readonly messages: string[]

  constructor(messages: string[], status: number) {
    super(messages[0] ?? "Unexpected error")
    this.name = "ApiError"
    this.messages = messages
    this.status = status
  }
}

/** Converts any thrown value (axios, network, unknown) into an ApiError. */
export function toApiError(error: unknown): ApiError {
  if (error instanceof ApiError) {
    return error
  }

  if (isAxiosError<ApiResult<unknown>>(error)) {
    const status = error.response?.status ?? 0
    const body = error.response?.data
    const messages =
      body?.messages && body.messages.length > 0
        ? body.messages
        : [error.message || "Network error"]
    return new ApiError(messages, status)
  }

  if (error instanceof Error) {
    return new ApiError([error.message], 0)
  }

  return new ApiError(["Unexpected error"], 0)
}
