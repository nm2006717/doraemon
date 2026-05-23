const BASE_URL = '/api'

function getToken(): string | null {
  return localStorage.getItem('token')
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>)
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(`${BASE_URL}${path}`, { ...options, headers })

  if (res.status === 401) {
    localStorage.removeItem('token')
    window.location.href = '/login'
    throw new Error('Unauthorized')
  }

  const data = await res.json()
  if (!res.ok) {
    throw new Error(data.error || 'Request failed')
  }
  return data as T
}

export const api = {
  getStatus: () => request<{ initialized: boolean }>('/status'),
  setup: (username: string, password: string) =>
    request<{ message: string }>('/setup', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    }),
  login: (username: string, password: string) =>
    request<{ token: string; expires_at: string }>('/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    }),
  getDashboard: () => request<{
    total_entries: number
    by_category: Record<string, number>
    recent_entries: any[]
  }>('/dashboard'),
  listEntries: (params: { category?: string; page?: number; per_page?: number }) => {
    const query = new URLSearchParams()
    if (params.category) query.set('category', params.category)
    if (params.page) query.set('page', String(params.page))
    if (params.per_page) query.set('per_page', String(params.per_page))
    return request<{
      entries: any[]
      total: number
      page: number
      per_page: number
      total_pages: number
    }>(`/entries?${query}`)
  },
  getEntry: (id: number) => request<any>(`/entries/${id}`),
  createEntry: (data: { title: string; content: string; category?: string; tags?: string }) =>
    request<any>('/entries', { method: 'POST', body: JSON.stringify(data) }),
  updateEntry: (id: number, data: { title?: string; content?: string; category?: string; tags?: string }) =>
    request<any>(`/entries/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteEntry: (id: number) => request<any>(`/entries/${id}`, { method: 'DELETE' }),
  searchEntries: (q: string) => request<any[]>(`/entries/search?q=${encodeURIComponent(q)}`)
}
