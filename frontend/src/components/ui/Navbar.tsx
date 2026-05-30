import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import UserAvatar from './UserAvatar'
import NotificationPanel from '../NotificationPanel'
import FriendRequestPanel from '../FriendRequestPanel'

export default function Navbar() {
  const { user, clearAuth } = useAuth()
  const navigate = useNavigate()
  const [search, setSearch] = useState('')
  const [showNotif, setShowNotif] = useState(false)
  const [showFriends, setShowFriends] = useState(false)

  function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    if (search.trim()) navigate(`/search?q=${encodeURIComponent(search.trim())}`)
  }

  function handleLogout() {
    clearAuth()
    navigate('/')
  }

  return (
    <header className="sticky top-0 z-40 bg-white shadow-sm border-b border-gray-200">
      <div className="max-w-4xl mx-auto px-4 h-14 flex items-center gap-4">
        {/* Logo */}
        <Link to="/feed" className="text-xl font-bold text-indigo-600 shrink-0">
          FreakDream
        </Link>

        {/* Search */}
        <form onSubmit={handleSearch} className="flex-1 max-w-xs">
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Pesquisar pessoas…"
            className="w-full rounded-full bg-gray-100 px-4 py-1.5 text-sm outline-none focus:ring-2 focus:ring-indigo-400"
          />
        </form>

        <div className="flex items-center gap-3 ml-auto">
          {/* Notifications */}
          <button
            onClick={() => setShowNotif((v) => !v)}
            className="relative text-gray-500 hover:text-indigo-600"
            title="Notificações"
          >
            🔔
          </button>

          {/* Friend requests */}
          <button
            onClick={() => setShowFriends((v) => !v)}
            className="text-gray-500 hover:text-indigo-600"
            title="Solicitações de amizade"
          >
            👥
          </button>

          {/* User menu */}
          {user && (
            <Link to={`/profile/${user.id}`}>
              <UserAvatar foto={user.foto} nome={user.nome} size="sm" />
            </Link>
          )}

          <button
            onClick={handleLogout}
            className="text-sm text-gray-500 hover:text-red-500"
            title="Sair"
          >
            Sair
          </button>
        </div>
      </div>

      {showNotif && <NotificationPanel onClose={() => setShowNotif(false)} />}
      {showFriends && <FriendRequestPanel onClose={() => setShowFriends(false)} />}
    </header>
  )
}
