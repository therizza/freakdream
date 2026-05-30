import { useEffect, useState } from 'react'
import { getFriendRequests, acceptFriendRequest, removeFriend } from '../api/endpoints'
import type { FriendRequest } from '../api/endpoints'

interface FriendRequestPanelProps {
  onClose: () => void
}

export default function FriendRequestPanel({ onClose }: FriendRequestPanelProps) {
  const [items, setItems] = useState<FriendRequest[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getFriendRequests()
      .then((r) => setItems(r.data))
      .finally(() => setLoading(false))
  }, [])

  async function handleAccept(id: number) {
    await acceptFriendRequest(id)
    setItems((prev) => prev.filter((f) => f.user_id !== id))
  }

  async function handleDecline(id: number) {
    await removeFriend(id)
    setItems((prev) => prev.filter((f) => f.user_id !== id))
  }

  return (
    <div className="absolute right-4 top-14 w-80 bg-white rounded-xl shadow-lg z-50 border border-gray-100">
      <div className="flex items-center justify-between px-4 py-3 border-b border-gray-100">
        <h3 className="font-semibold text-sm">Solicitações de amizade</h3>
        <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-xs">✕</button>
      </div>
      <div className="max-h-80 overflow-y-auto">
        {loading && <p className="text-center py-6 text-sm text-gray-400">Carregando…</p>}
        {!loading && items.length === 0 && (
          <p className="text-center py-6 text-sm text-gray-400">Nenhuma solicitação</p>
        )}
        {items.map((f) => (
          <div key={f.user_id} className="flex items-center gap-3 px-4 py-3 border-b border-gray-50 last:border-0">
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-700">Usuário #{f.user_id}</p>
              <p className="text-xs text-gray-400">{new Date(f.criado_em).toLocaleString('pt-BR')}</p>
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => handleAccept(f.user_id)}
                className="text-xs px-2 py-1 bg-indigo-600 text-white rounded"
              >
                Aceitar
              </button>
              <button
                onClick={() => handleDecline(f.user_id)}
                className="text-xs px-2 py-1 border border-gray-300 rounded text-gray-600"
              >
                Recusar
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
