import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
// Navbar is rendered by the App wrapper; do not render it here to avoid duplicates.
import { decodeJWT } from '../utils/jwt';
import { fetchTasksByProject } from '../api/tasks';

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  created_at: string;
}

interface Project {
  id: string;
  name: string;
  description?: string;
  created_by?: string;
  assigned_employees?: string[];
}

export const Dashboard = () => {
  const navigate = useNavigate();
  const [users, setUsers] = useState<User[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [totalTasks, setTotalTasks] = useState<number | null>(null);
  const [totalAssignees, setTotalAssignees] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [projectsLoading, setProjectsLoading] = useState(true);
  const [projectsError, setProjectsError] = useState('');
  const [activeProject, setActiveProject] = useState<Project | null>(null);
  const [userRole, setUserRole] = useState<string | null>(null);

  // Redirect to login if not authenticated and extract user role
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/', { replace: true });
    } else {
      const decoded = decodeJWT(token);
      if (decoded && decoded.role) {
        setUserRole(decoded.role);
      }
    }
  }, [navigate]);

  // Fetch users (only if ADMIN)
  useEffect(() => {
    const fetchUsers = async () => {
      if (userRole !== 'ADMIN') {
        setLoading(false);
        return;
      }

      try {
        const token = localStorage.getItem('token');
        const response = await fetch('http://localhost:8080/api/users', {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error('Failed to fetch users');
        }

        const data = await response.json();
        // backend returns { users: [...] }
        setUsers((data && data.users) || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    if (userRole) {
      fetchUsers();
    }
  }, [userRole]);

  // Fetch projects
  useEffect(() => {
    const fetchProjects = async () => {
      try {
        const token = localStorage.getItem('token');
        const res = await fetch('http://localhost:8080/projects', {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (!res.ok) throw new Error('Failed to fetch projects');
        const data = await res.json();
        // handler returns array of project objects
        const proj: Project[] = data || [];
        setProjects(proj);

        // compute total unique assignees across projects
        const assigneeSet = new Set<string>();
        for (const p of proj) {
          (p.assigned_employees || []).forEach((id) => assigneeSet.add(id));
        }
        setTotalAssignees(assigneeSet.size);

        // fetch tasks for each project and sum counts
        try {
          const counts = await Promise.all(
            proj.map(async (p) => {
              const tasks = await fetchTasksByProject(p.id).catch(() => [] as any[]);
              return Array.isArray(tasks) ? tasks.length : 0;
            })
          );
          setTotalTasks(counts.reduce((a, b) => a + b, 0));
        } catch (err) {
          setTotalTasks(null);
        }
      } catch (err) {
        setProjectsError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setProjectsLoading(false);
      }
    };

    fetchProjects();
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Header */}
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-white mb-2">Dashboard</h2>
          <p className="text-slate-400">Manage your RBAC system</p>
        </div>

        {/* Stats Grid (only show for ADMIN) */}
        {userRole === 'ADMIN' && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            {/* Total Users Card */}
            <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-slate-400 text-sm font-medium mb-2">Total Users</p>
                  <p className="text-3xl font-bold text-white">{users.length}</p>
                </div>
                <div className="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center">
                  <svg className="w-6 h-6 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM9 6a3 3 0 11-6 0 3 3 0 016 0zM9 6a3 3 0 11-6 0 3 3 0 016 0zm12 6a9 9 0 11-18 0 9 9 0 0118 0zm-9-5a4 4 0 11-8 0 4 4 0 018 0zM9 9a6 6 0 1 0 12 0 6 6 0 0 0-12 0z" />
                  </svg>
                </div>
              </div>
            </div>

            {/* Total Projects Card */}
            <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-slate-400 text-sm font-medium mb-2">Total Projects</p>
                  <p className="text-3xl font-bold text-white">{projects.length}</p>
                </div>
                <div className="w-12 h-12 bg-green-500/20 rounded-xl flex items-center justify-center">
                  <svg className="w-6 h-6 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M3 3h14v2H3V3zm0 4h10v2H3V7zM3 11h8v2H3v-2z" />
                  </svg>
                </div>
              </div>
            </div>

            {/* Total Tasks Card */}
            <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-slate-400 text-sm font-medium mb-2">Total Tasks</p>
                  <p className="text-3xl font-bold text-white">{(totalTasks !== null && totalTasks > 0) ? totalTasks : (totalAssignees !== null ? totalAssignees : '—')}</p>
                </div>
                <div className="w-12 h-12 bg-purple-500/20 rounded-xl flex items-center justify-center">
                  <svg className="w-6 h-6 text-purple-400" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M5 3h10v2H5V3zm0 4h10v2H5V7zm0 4h7v2H5v-2z" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        )}
        {/* Projects Section */}
        <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6 mb-8">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-xl font-bold text-white">Projects</h3>
            {userRole !== 'VIEWER' && (
              <button
                onClick={() => navigate('/create-project')}
                className="px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-semibold rounded-lg transition-all duration-200"
              >
                + New Project
              </button>
            )}
          </div>

          {projectsLoading ? (
            <div className="text-center py-8">
              <div className="inline-block">
                <svg className="animate-spin h-8 w-8 text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
            </div>
          ) : projectsError ? (
            <div className="p-4 bg-red-500/20 border border-red-500/50 rounded-lg text-red-200">{projectsError}</div>
          ) : projects.length === 0 ? (
            <div className="text-center py-8 text-slate-400">No projects yet</div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {projects.map((p) => (
                <button
                  key={p.id}
                  onClick={() => setActiveProject(p)}
                  className="text-left bg-slate-800/40 hover:bg-slate-800/60 border border-slate-700/50 rounded-2xl p-6 transition-all"
                >
                  <h4 className="text-lg font-semibold text-white mb-2">{p.name}</h4>
                  <p className="text-sm text-slate-400 mb-4">{p.description || 'No description'}</p>
                  <div className="flex items-center gap-2">
                    <span className="inline-block px-2 py-1 bg-blue-500/10 text-blue-300 rounded-full text-xs whitespace-nowrap">{(p.assigned_employees || []).length} assignees</span>
                    <span className="text-xs text-slate-400">Created by {p.created_by || 'system'}</span>
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Users Section (only for ADMIN) */}
        {userRole === 'ADMIN' && (
        <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-xl font-bold text-white">Users</h3>
            <button
              onClick={() => navigate('/create-user')}
              className="px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white font-semibold rounded-lg transition-all duration-200"
            >
              + Add User
            </button>
          </div>

          {loading ? (
            <div className="text-center py-12">
              <div className="inline-block">
                <svg className="animate-spin h-8 w-8 text-blue-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
            </div>
          ) : error ? (
            <div className="p-4 bg-red-500/20 border border-red-500/50 rounded-lg text-red-200">
              {error}
            </div>
          ) : users.length === 0 ? (
            <div className="text-center py-12">
              <svg className="w-12 h-12 text-slate-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.856-1.487M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 0a2 2 0 11-4 0 2 2 0 014 0zM5 20a3 3 0 015.856-1.487M5 10a2 2 0 110-4 2 2 0 010 4z" />
              </svg>
              <p className="text-slate-400">No users found</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-slate-700/50">
                    <th className="text-left py-3 px-4 text-slate-300 font-semibold">ID</th>
                    <th className="text-left py-3 px-4 text-slate-300 font-semibold">Name</th>
                    <th className="text-left py-3 px-4 text-slate-300 font-semibold">Email</th>
                    <th className="text-left py-3 px-4 text-slate-300 font-semibold">Role</th>
                    <th className="text-left py-3 px-4 text-slate-300 font-semibold">Created At</th>
                    <th className="text-left py-3 px-4 text-slate-300 font-semibold">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((user) => (
                    <tr key={user.id} className="border-b border-slate-700/30 hover:bg-slate-700/20 transition-colors">
                      <td className="py-3 px-4 text-slate-200 text-sm">{user.id}</td>
                      <td className="py-3 px-4 text-slate-200">{user.name}</td>
                      <td className="py-3 px-4 text-slate-200">{user.email}</td>
                      <td className="py-3 px-4">
                        <span className="inline-block px-3 py-1 bg-blue-500/20 text-blue-300 rounded-full text-sm font-medium">
                          {user.role}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-slate-400 text-sm">
                        {new Date(user.created_at).toLocaleDateString()}
                      </td>
                      <td className="py-3 px-4">
                        <button className="text-blue-400 hover:text-blue-300 transition-colors text-sm font-medium">
                          Edit
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
        )}
        {/* Project Detail Modal */}
        {activeProject && (
          <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div className="absolute inset-0 bg-black/60" onClick={() => setActiveProject(null)}></div>
            <div className="relative z-10 max-w-2xl w-full mx-4 bg-slate-800/60 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="text-2xl font-bold text-white">{activeProject.name}</h3>
                  <p className="text-slate-400 text-sm">{activeProject.description}</p>
                </div>
                <div className="flex items-center gap-3">
                  <button
                    onClick={() => { navigate(`/projects/${activeProject.id}`); setActiveProject(null); }}
                    className="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm"
                  >
                    View more
                  </button>
                  <button onClick={() => setActiveProject(null)} className="text-slate-300 hover:text-white">Close</button>
                </div>
              </div>

              <div className="mt-6">
                <h4 className="text-sm text-slate-300 font-semibold mb-2">Assignees</h4>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {(activeProject.assigned_employees || []).map((uid) => {
                    const u = users.find((x) => x.id === uid);
                    return (
                      <div key={uid} className="p-3 bg-slate-700/40 rounded-lg">
                        <div className="text-white font-medium">{u?.name || u?.email || uid}</div>
                        <div className="text-xs text-slate-400">{u?.email || ''}</div>
                        <div className="text-xs text-slate-400">{u?.role || ''}</div>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
