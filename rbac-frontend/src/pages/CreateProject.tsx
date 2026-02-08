import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Navbar } from '../components/Navbar';

interface User {
  id: string;
  name?: string;
  email: string;
  role?: string;
}

export const CreateProject = () => {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [users, setUsers] = useState<User[]>([]);
  const [selected, setSelected] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/', { replace: true });
      return;
    }

    const fetchUsers = async () => {
      try {
        const res = await fetch('http://localhost:8080/api/users', {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (!res.ok) throw new Error('Failed to load users');
        const data = await res.json();
        setUsers(data.users || []);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error');
      }
    };
    fetchUsers();
  }, [navigate]);

  const toggleSelect = (id: string) => {
    setSelected((prev) => (prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);

    try {
      const token = localStorage.getItem('token');
      const body = { name, description, assigned_employees: selected };
      const res = await fetch('http://localhost:8080/projects/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (!res.ok) {
        const d = await res.json().catch(() => ({}));
        throw new Error(d.message || 'Failed to create project');
      }

      setSuccess('Project created');
      setName('');
      setDescription('');
      setSelected([]);
      setTimeout(() => navigate('/dashboard'), 1200);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <Navbar />
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <h2 className="text-3xl font-bold text-white mb-2">Create Project</h2>
        <p className="text-slate-400 mb-6">Define project and assign employees</p>

        <div className="bg-slate-800/50 backdrop-blur-md border border-slate-700/50 rounded-2xl p-6">
          {error && <div className="p-3 bg-red-500/20 text-red-200 rounded mb-4">{error}</div>}
          {success && <div className="p-3 bg-green-500/20 text-green-200 rounded mb-4">{success}</div>}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm text-slate-200 mb-2">Name</label>
              <input value={name} onChange={(e) => setName(e.target.value)} required className="w-full px-4 py-3 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white" />
            </div>
            <div>
              <label className="block text-sm text-slate-200 mb-2">Description</label>
              <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={4} className="w-full px-4 py-3 bg-slate-700/50 border border-slate-600/50 rounded-lg text-white" />
            </div>

            <div>
              <label className="block text-sm text-slate-200 mb-2">Assign Employees</label>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-2 max-h-48 overflow-y-auto p-2">
                {users.map((u) => (
                  <label key={u.id} className="flex items-center gap-2 p-2 rounded hover:bg-slate-700/30 cursor-pointer">
                    <input type="checkbox" checked={selected.includes(u.id)} onChange={() => toggleSelect(u.id)} className="w-4 h-4" />
                    <div>
                      <div className="text-slate-200">{u.email}{u.name ? ` — ${u.name}` : ''}</div>
                      <div className="text-xs text-slate-400">{u.role || ''}</div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            <div className="flex gap-4">
              <button type="submit" disabled={loading} className="flex-1 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg">{loading ? 'Creating...' : 'Create Project'}</button>
              <button type="button" onClick={() => navigate('/dashboard')} className="flex-1 py-3 bg-slate-700/50 text-white rounded-lg">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};
