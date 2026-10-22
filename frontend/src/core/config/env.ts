export const env = {
  apiBaseUrl: import.meta.env.VITE_API_URL ?? "http://localhost:7080",
} as const
