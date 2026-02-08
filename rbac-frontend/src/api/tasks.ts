const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

async function request(path: string, opts: RequestInit = {}) {
  const token = localStorage.getItem('token')
  const headers: Record<string,string> = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`
  const res = await fetch(API_BASE + path, { headers, ...opts })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function fetchTasksByProject(projectId: string) {
  return request(`/tasks?project_id=${encodeURIComponent(projectId)}`)
}

export async function updateTask(payload: Record<string, any>) {
  return request('/tasks/update', { method: 'POST', body: JSON.stringify(payload) })
}

export async function createTask(payload: Record<string, any>) {
  return request('/tasks/create', { method: 'POST', body: JSON.stringify(payload) })
}
