import { useState } from 'react'
import { Link } from 'react-router-dom'
import type { Post } from '../../api/endpoints'
import { deletePost, sharePost } from '../../api/endpoints'
import UserAvatar from './UserAvatar'
import ConfirmDeleteModal from './ConfirmDeleteModal'
import { useAuthStore } from '../../store/authStore'

interface PostCardProps {
  post: Post
  onDeleted?: (postId: number) => void
}

export default function PostCard({ post, onDeleted }: PostCardProps) {
  const { user } = useAuthStore()
  const [confirmDelete, setConfirmDelete] = useState(false)
  const [sharing, setSharing] = useState(false)

  async function handleDelete() {
    await deletePost(post.post_id)
    onDeleted?.(post.post_id)
  }

  async function handleShare() {
    if (sharing) return
    setSharing(true)
    try {
      await sharePost(post.post_id)
      window.dispatchEvent(new Event('feed:refresh'))
    } finally {
      setSharing(false)
    }
  }

  const isOwner = user?.id === post.user_id

  return (
    <article className="bg-white rounded-xl shadow p-4 space-y-3">
      {/* Header */}
      <div className="flex items-center gap-3">
        <Link to={`/profile/${post.user_id}`}>
          <UserAvatar foto={post.user?.foto} nome={post.user?.nome} />
        </Link>
        <div className="flex-1 min-w-0">
          <Link to={`/profile/${post.user_id}`} className="font-semibold text-gray-900 hover:underline">
            {post.user?.nome || `Usuário ${post.user_id}`}
          </Link>
          <p className="text-xs text-gray-400">{new Date(post.data).toLocaleString('pt-BR')}</p>
        </div>
        {isOwner && (
          <button
            onClick={() => setConfirmDelete(true)}
            className="text-gray-400 hover:text-red-500 text-sm"
            title="Excluir"
          >
            🗑️
          </button>
        )}
      </div>

      {/* Content */}
      {post.tipo === 1 && <p className="text-gray-800 text-sm whitespace-pre-wrap">{post.texto}</p>}

      {post.tipo === 2 && (
        <div className="space-y-2">
          <img
            src={`/uploads/${post.texto}`}
            alt="foto"
            className="rounded-lg w-full max-h-96 object-cover"
            onError={(e) => (e.currentTarget.style.display = 'none')}
          />
        </div>
      )}

      {post.tipo === 3 && (
        <div className="border border-gray-100 rounded-lg p-3 bg-gray-50 text-sm text-gray-600">
          🔁 Compartilhou publicação #{post.texto}
        </div>
      )}

      {/* Actions */}
      <div className="flex gap-4 pt-1 border-t border-gray-100">
        {post.tipo !== 3 && (
          <button
            onClick={handleShare}
            disabled={sharing}
            className="text-xs text-gray-500 hover:text-indigo-600"
          >
            🔁 Compartilhar
          </button>
        )}
      </div>

      {confirmDelete && (
        <ConfirmDeleteModal
          onConfirm={handleDelete}
          onCancel={() => setConfirmDelete(false)}
        />
      )}
    </article>
  )
}
