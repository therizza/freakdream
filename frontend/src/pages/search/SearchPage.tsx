import { useEffect, useState } from 'react'
import { useSearchParams, Link } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import Navbar from '../../components/ui/Navbar'
import UserAvatar from '../../components/ui/UserAvatar'
import { searchUsers } from '../../api/endpoints'
import type { User } from '../../api/endpoints'

export default function SearchPage() {
  useAuth()
  const [searchParams] = useSearchParams()
  const q = searchParams.get('q') ?? ''
  const [results, setResults] = useState<User[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!q.trim()) return
    setLoading(true)
    searchUsers(q)
      .then((r) => setResults(r.data))
      .finally(() => setLoading(false))
  }, [q])

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-xl mx-auto px-4 py-6 space-y-4">
        <h2 className="text-lg font-semibold text-gray-700">
          {q ? `Resultados para "${q}"` : 'Busca'}
        </h2>

        {loading && <p className="text-sm text-gray-400">Buscando…</p>}

        {!loading && results.length === 0 && q && (
          <p className="text-sm text-gray-400">Nenhum usuário encontrado.</p>
        )}

        <div className="space-y-3">
          {results.map((user) => (
            <Link
              key={user.id}
              to={`/profile/${user.id}`}
              className="flex items-center gap-4 bg-white rounded-xl shadow p-4 hover:bg-gray-50 transition"
            >
              <UserAvatar foto={user.foto} nome={user.nome} size="md" />
              <div>
                <p className="font-semibold text-gray-900">{user.nome}</p>
                {user.sobre && <p className="text-sm text-gray-500">{user.sobre}</p>}
              </div>
            </Link>
          ))}
        </div>
      </main>
    </div>
  )
}
