import { useEffect, useState } from 'react'
import { fetchTasksByProject, updateTask } from '../api/tasks'

type Task = {
  id: string
  title: string
  description?: string
  status: string
  assignees?: string[]
}

const STATUSES = ['TODO', 'IN_PROGRESS', 'REVIEW', 'DONE']

export default function KanbanBoard({ projectId, usersMap = {} }: { projectId: string; usersMap?: Record<string, string> }) {
  const [tasks, setTasks] = useState<Task[]>([])

  useEffect(() => {
    if (!projectId) return
    fetchTasksByProject(projectId).then(setTasks).catch(console.error)
  }, [projectId])

  const move = async (id: string, status: string) => {
    await updateTask({ id, status })
    setTasks(prev => prev.map(t => t.id === id ? { ...t, status } : t))
  }

  return (
    <div className="grid grid-cols-4 gap-4">
      {STATUSES.map(s => (
        <div key={s} className="bg-slate-800/40 p-3 rounded-md">
          <h3 className="text-sm font-semibold text-slate-200 mb-2">{s.replace('_', ' ')}</h3>
          {tasks?.filter(t => t.status === s).map(t => (
            <div key={t.id} className="bg-white/5 p-3 mb-3 rounded-md">
              <div className="font-semibold text-white">{t.title}</div>
              {t.description && <div className="text-slate-400 text-sm">{t.description}</div>}
              {t.assignees && t.assignees.length > 0 && (
                <div className="mt-2 text-xs text-slate-300">Assigned: {t.assignees.map(id => usersMap[id] || id).join(', ')}</div>
              )}
              <div className="mt-2 grid grid-cols-2 gap-1">
                {STATUSES.map(opt => (
                  <button key={opt} disabled={opt === s} onClick={() => move(t.id, opt)} className="px-2 py-1 text-xs bg-slate-700/50 rounded disabled:opacity-50">
                    {opt.split('_')[0]}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>
      ))}
    </div>
  )
}
