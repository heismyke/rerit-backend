const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api'
const REQUEST_TIMEOUT_MS = 8000

type LoginPayload = {
  email: string
  password: string
  role: string
}

const request = async <T>(path: string, options: RequestInit = {}): Promise<T> => {
  const controller = new AbortController()
  const timeout = window.setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS)

  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    signal: controller.signal,
    ...options,
  }).finally(() => window.clearTimeout(timeout))

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(body.error || 'Request failed')
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json() as Promise<T>
}

export const api = {
  login: (payload: LoginPayload) =>
    request<{ user: { email: string; name: string; role: string }; token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  getAuditCases: <T>() => request<T>('/audit-cases'),
  createAuditCase: <T>(payload: unknown) =>
    request<T>('/audit-cases', { method: 'POST', body: JSON.stringify(payload) }),
  updateAuditCase: <T>(id: string, payload: unknown) =>
    request<T>(`/audit-cases/${encodeURIComponent(id)}`, { method: 'PUT', body: JSON.stringify(payload) }),
  deleteAuditCase: (id: string) =>
    request<void>(`/audit-cases/${encodeURIComponent(id)}`, { method: 'DELETE' }),
  sendAuditResult: <T>(id: string, payload: { status: string; notes: string }) =>
    request<T>(`/audit-cases/${encodeURIComponent(id)}/result`, { method: 'POST', body: JSON.stringify(payload) }),
  getFlaggedCases: <T>() => request<T>('/flagged-cases'),
  updateFlaggedCase: <T>(id: string, payload: unknown) =>
    request<T>(`/flagged-cases/${encodeURIComponent(id)}`, { method: 'PUT', body: JSON.stringify(payload) }),
  sendFlaggedResult: <T>(id: string, payload: { status: string; notes: string }) =>
    request<T>(`/flagged-cases/${encodeURIComponent(id)}/result`, { method: 'POST', body: JSON.stringify(payload) }),
  getSuccessfulFilings: <T>() => request<T>('/successful-filings'),
  updateSuccessfulFiling: <T>(id: string, payload: unknown) =>
    request<T>(`/successful-filings/${encodeURIComponent(id)}`, { method: 'PUT', body: JSON.stringify(payload) }),
  getProperties: <T>() => request<T>('/properties'),
  createProperty: <T>(payload: unknown) =>
    request<T>('/properties', { method: 'POST', body: JSON.stringify(payload) }),
  updateProperty: <T>(id: string, payload: unknown) =>
    request<T>(`/properties/${encodeURIComponent(id)}`, { method: 'PUT', body: JSON.stringify(payload) }),
  deleteProperty: (id: string) =>
    request<void>(`/properties/${encodeURIComponent(id)}`, { method: 'DELETE' }),
  getNotices: <T>() => request<T>('/notices'),
  respondNotice: <T>(id: string, response: string) =>
    request<T>(`/notices/${encodeURIComponent(id)}/respond`, { method: 'POST', body: JSON.stringify({ response }) }),
  createPayment: <T>(payload: unknown) =>
    request<T>('/payments', { method: 'POST', body: JSON.stringify(payload) }),
}
