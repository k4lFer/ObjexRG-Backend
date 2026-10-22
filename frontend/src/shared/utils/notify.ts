import { toast } from "sonner"

import { toApiError } from "../../core/http/apiError"

/** Shows every message from a normalized API error as a toast. */
export function notifyError(error: unknown): void {
  const { messages } = toApiError(error)
  const list = messages.length > 0 ? messages : ["Something went wrong"]
  list.forEach((message) => toast.error(message))
}
