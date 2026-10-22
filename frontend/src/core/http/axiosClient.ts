import axios from "axios"

import { env } from "../config/env"
import { attachInterceptors } from "./interceptors"

/**
 * Shared axios instance. Every feature repository should use this one.
 *
 * No default Content-Type is set on purpose: axios adds `application/json`
 * for object payloads, and lets the browser set the correct
 * `multipart/form-data; boundary=...` when the body is a FormData. A global
 * JSON default would make axios serialize FormData uploads to JSON.
 */
export const httpClient = axios.create({
  baseURL: env.apiBaseUrl,
})

attachInterceptors(httpClient)
