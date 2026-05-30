import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login, register } from '../../api/endpoints'
import { useAuthStore } from '../../store/authStore'

export default function LoginPage() {
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [form, setForm] = useState({ nome: '', sobre: '', email: '', senha: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { setAuth } = useAuthStore()
  const navigate = useNavigate()

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      let res
      if (mode === 'login') {
        res = await login(form.email, form.senha)
      } else {
        res = await register(form.nome, form.sobre, form.email, form.senha)
      }
      setAuth(res.data.user, res.data.token)
      navigate('/feed')
    } catch (err: any) {
      setError(err.response?.data?.error ?? 'Erro ao autenticar')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 to-purple-50 p-4">
      <div className="w-full max-w-sm bg-white rounded-2xl shadow-lg p-8 space-y-6">
        {/* Logo */}
        <div className="text-center">
          <h1 className="text-3xl font-extrabold text-indigo-600">FreakDream</h1>
          <p className="text-sm text-gray-400 mt-1">Conecte-se ao mundo</p>
        </div>

        {/* Tabs */}
        <div className="flex rounded-full bg-gray-100 p-1">
          {(['login', 'register'] as const).map((m) => (
            <button
              key={m}
              onClick={() => setMode(m)}
              className={`flex-1 py-1.5 text-sm font-medium rounded-full transition ${
                mode === m ? 'bg-white shadow text-indigo-600' : 'text-gray-500'
              }`}
            >
              {m === 'login' ? 'Entrar' : 'Cadastrar'}
            </button>
          ))}
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-3">
          {mode === 'register' && (
            <>
              <input
                name="nome"
                value={form.nome}
                onChange={handleChange}
                placeholder="Nome"
                required
                className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-indigo-400"
              />
              <input
                name="sobre"
                value={form.sobre}
                onChange={handleChange}
                placeholder="Sobre você (opcional)"
                className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-indigo-400"
              />
            </>
          )}
          <input
            name="email"
            type="email"
            value={form.email}
            onChange={handleChange}
            placeholder="Email"
            required
            className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-indigo-400"
          />
          <input
            name="senha"
            type="password"
            value={form.senha}
            onChange={handleChange}
            placeholder="Senha"
            required
            minLength={mode === 'register' ? 8 : undefined}
            className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-indigo-400"
          />

          {error && <p className="text-xs text-red-500">{error}</p>}

          <button
            type="submit"
            disabled={loading}
            className="w-full py-2 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 disabled:opacity-50"
          >
            {loading ? 'Aguarde…' : mode === 'login' ? 'Entrar' : 'Cadastrar'}
          </button>
        </form>
      </div>
    </div>
  )
}
