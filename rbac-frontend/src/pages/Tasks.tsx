import React, { useState } from 'react'
import KanbanBoard from '../components/KanbanBoard'
import { createTask } from '../api/tasks'

export default function Tasks() {
  const [projectId, setProjectId] = useState('')
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [loading, setLoading] = useState(false)

  const onCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!projectId || !title) return alert('project and title required')
    setLoading(true)
    try {
      await createTask({ project_id: projectId, title, description })
      setTitle('')
      setDescription('')
      // naive refresh: reload page to fetch tasks
      window.location.reload()
    } catch (err) {
      alert('create failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h2 className="text-3xl font-bold text-white mb-2">Tasks</h2>
        <p className="text-slate-400 mb-6">Create and manage tasks for a project</p>

        <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6 mb-8">
          <form onSubmit={onCreate} className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <input
              placeholder="Project ID"
              value={projectId}
              onChange={e => setProjectId(e.target.value)}
              className="px-4 py-3 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white"
              required
            />
            <input
              placeholder="Title"
              value={title}
              onChange={e => setTitle(e.target.value)}
              className="px-4 py-3 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white"
              required
            />
            <div className="flex gap-2">
              <input
                placeholder="Description (optional)"
                value={description}
                onChange={e => setDescription(e.target.value)}
                className="flex-1 px-4 py-3 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white"
              />
              <button
                type="submit"
                disabled={loading}
                className="px-4 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg"
              >
                {loading ? 'Creating...' : 'Create'}
              </button>
            </div>
          </form>
        </div>

        {projectId ? (
          <KanbanBoard projectId={projectId} />
        ) : (
          <div className="text-slate-400">Please enter a Project ID above to load tasks.</div>
        )}
      </div>
    </div>
  )
}
