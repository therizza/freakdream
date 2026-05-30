import { useEffect, useState } from 'react'
import { getNotifications, removeNotification, markNotificationRead } from '../api/endpoints'
import type { Notification } from '../api/endpoints'

interface NotificationPanelProps {
  onClose: () => void
}

export default function NotificationPanel({ onClose }: NotificationPanelProps) {
  const [items, setItems] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getNotifications()
      .then((r) => setItems(r.data))
      .finally(() => setLoading(false))
  }, [])

  async function handleRemove(id: number) {
    await removeNotification(id)
    setItems((prev) => prev.filter((n) => n.id !== id))
  }

  async function handleRead(id: number) {
    await markNotificationRead(id)
    setItems((prev) => prev.map((n) => (n.id === id ? { ...n, lida: true } : n)))
  }

  return (
    <div className="absolute right-4 top-14 w-80 bg-white rounded-xl shadow-lg z-50 border border-gray-100">
      <div className="flex items-center justify-between px-4 py-3 border-b border-gray-100">
        <h3 className="font-semibold text-sm">Notificações</h3>
        <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-xs">✕</button>
      </div>
      <div className="max-h-80 overflow-y-auto">
        {loading && <p className="text-center py-6 text-sm text-gray-400">Carregando…</p>}
        {!loading && items.length === 0 && (
          <p className="text-center py-6 text-sm text-gray-400">Nenhuma notificação</p>
        )}
        {items.map((n) => (
          <div
            key={n.id}
            className={`flex items-start gap-3 px-4 py-3 border-b border-gray-50 last:border-0 ${n.lida ? 'opacity-60' : ''}`}
          >
            <div className="flex-1 min-w-0">
              <p className="text-xs text-gray-700 font-medium">{n.tipo}</p>
              <p className="text-xs text-gray-500 truncate">{n.texto}</p>
              <p className="text-xs text-gray-400">{new Date(n.data).toLocaleString('pt-BR')}</p>
            </div>
            <div className="flex gap-1">
              {!n.lida && (
                <button onClick={() => handleRead(n.id)} className="text-xs text-indigo-500">✓</button>
              )}
              <button onClick={() => handleRemove(n.id)} className="text-xs text-red-400">✕</button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
