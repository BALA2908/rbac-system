import React, { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import KanbanBoard from '../components/KanbanBoard'
import { decodeJWT } from '../utils/jwt'
import { createTask } from '../api/tasks'

type Project = {
  id: string
  name: string
  description?: string
  assigned_employees?: string[]
  created_by?: string
}

export default function ProjectDetails() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [project, setProject] = useState<Project | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!id) return
    const fetchProject = async () => {
      try {
        const token = localStorage.getItem('token')
        const res = await fetch('http://localhost:8080/projects', { headers: { Authorization: `Bearer ${token}` } })
        if (!res.ok) throw new Error('Failed to fetch projects')
        const data = await res.json()
        const found = (data || []).find((p: any) => p.id === id)
        if (!found) {
          setError('Project not found or not accessible')
        } else {
          setProject(found)
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error')
      } finally {
        setLoading(false)
      }
    }
    fetchProject()
  }, [id])

  const tokenRole = (() => {
    const t = localStorage.getItem('token')
    const d = t ? decodeJWT(t) : null
    return d?.role || null
  })()

  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [assignees, setAssignees] = useState<string[]>([])
  const [usersMap, setUsersMap] = useState<Record<string,string>>({})
  const [creating, setCreating] = useState(false)
  const [allUsers, setAllUsers] = useState<Array<{id: string, name?: string, email: string}>>([])

  const canCreate = tokenRole === 'ADMIN' || tokenRole === 'MANAGER' || tokenRole === 'EDITOR'

  const onCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!project) return
    if (!title) return alert('title required')
    setCreating(true)
    try {
      await createTask({ project_id: project.id, title, description, assignees })
      setTitle('')
      setDescription('')
      setAssignees([])
      window.location.reload()
    } catch (err) {
      alert('create failed')
    } finally {
      setCreating(false)
    }
  }

  useEffect(() => {
    // Fetch all users to display names/emails for assignee selection and display
    if (!project) return

    (async () => {
      try {
        const token = localStorage.getItem('token')
        const res = await fetch('http://localhost:8080/api/users', { headers: { Authorization: `Bearer ${token}` } })
        if (res.ok) {
          const data = await res.json()
          const users: any[] = data.users || []
          if (users.length > 0) {
            setAllUsers(users)
            const map: Record<string,string> = {}
            for (const u of users) {
              map[u.id] = u.name || u.email || u.id
            }
            setUsersMap(map)
            return
          }
        } else {
          console.warn('Failed to fetch /api/users', res.status)
        }

        // Fallback: use project.assigned_employees when API is not available or returns empty
        const fallback = (project.assigned_employees || []).map((uid: string) => ({ id: uid, name: usersMap[uid] || uid, email: '' }))
        setAllUsers(fallback)
        const map2: Record<string,string> = {}
        for (const u of fallback) {
          map2[u.id] = u.name || u.email || u.id
        }
        setUsersMap(map2)
      } catch (err) {
        console.warn('Error fetching users', err)
        const fallback = (project.assigned_employees || []).map((uid: string) => ({ id: uid, name: usersMap[uid] || uid, email: '' }))
        setAllUsers(fallback)
        const map2: Record<string,string> = {}
        for (const u of fallback) {
          map2[u.id] = u.name || u.email || u.id
        }
        setUsersMap(map2)
      }
    })()
  }, [project])

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-3xl font-bold text-white">Project Details</h2>
            <p className="text-slate-400">View tasks and information for this project</p>
          </div>
          <div className="flex gap-2">
            <button onClick={() => navigate('/dashboard')} className="px-4 py-2 bg-slate-700/50 text-white rounded-lg">Back</button>
          </div>
        </div>

        <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6 mb-8">
          {loading ? (
            <div className="text-slate-400">Loading...</div>
          ) : error ? (
            <div className="text-red-400">{error}</div>
          ) : project ? (
            <div>
              <h3 className="text-2xl text-white font-semibold">{project.name}</h3>
              <p className="text-slate-400 mb-4">{project.description || 'No description'}</p>
              <div className="mb-4">
                <span className="inline-block px-2 py-1 bg-blue-500/10 text-blue-300 rounded-full text-xs">{(project.assigned_employees || []).length} assignees</span>
                <span className="ml-3 text-xs text-slate-400">Created by {project.created_by || 'system'}</span>
              </div>
              <div className="mt-6">
                {canCreate && (
                  <div className="bg-slate-800/40 p-4 rounded-md mb-6">
                    <h4 className="text-white font-semibold mb-2">Create Task</h4>
                    <form onSubmit={onCreate} className="grid grid-cols-1 md:grid-cols-4 gap-3">
                      <input value={title} onChange={e => setTitle(e.target.value)} placeholder="Title" className="px-3 py-2 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white" required />
                      <input value={description} onChange={e => setDescription(e.target.value)} placeholder="Description" className="px-3 py-2 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white" />
                      <select multiple value={assignees} onChange={e => setAssignees(Array.from(e.target.selectedOptions).map(o => o.value))} className="px-3 py-2 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white">
                        {allUsers.length === 0 && <option value="">No users available</option>}
                        {allUsers.map(u => (
                          <option key={u.id} value={u.id}>{u.name || u.email}</option>
                        ))}
                      </select>
                      <div>
                        <button type="submit" disabled={creating} className="px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg">{creating ? 'Creating...' : 'Create'}</button>
                      </div>
                    </form>
                  </div>
                )}

                <KanbanBoard projectId={project.id} usersMap={usersMap} />
              </div>
            </div>
          ) : null}
        </div>
      </div>
    </div>
  )
}
