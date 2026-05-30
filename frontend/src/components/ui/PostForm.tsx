import { useState } from 'react'
import { createPost, uploadPhoto } from '../../api/endpoints'
import type { Post } from '../../api/endpoints'
import PhotoUploadModal from './PhotoUploadModal'

interface PostFormProps {
  onPosted: (post: Post) => void
}

export default function PostForm({ onPosted }: PostFormProps) {
  const [texto, setTexto] = useState('')
  const [loading, setLoading] = useState(false)
  const [showPhoto, setShowPhoto] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!texto.trim() || loading) return
    setLoading(true)
    try {
      const res = await createPost(texto.trim())
      onPosted(res.data)
      setTexto('')
    } finally {
      setLoading(false)
    }
  }

  async function handlePhotoPost(fd: FormData) {
    setLoading(true)
    try {
      await uploadPhoto(fd)
      setShowPhoto(false)
      // Refresh is handled by parent via re-fetch
      window.dispatchEvent(new Event('feed:refresh'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="bg-white rounded-xl shadow p-4 space-y-3">
        <textarea
          value={texto}
          onChange={(e) => setTexto(e.target.value)}
          placeholder="Diga o que você pensa…"
          rows={3}
          className="w-full resize-none border border-gray-200 rounded-lg p-3 text-sm outline-none focus:ring-2 focus:ring-indigo-400"
        />
        <div className="flex items-center justify-between">
          <button
            type="button"
            onClick={() => setShowPhoto(true)}
            className="text-sm text-gray-500 hover:text-indigo-600"
          >
            📷 Foto
          </button>
          <button
            type="submit"
            disabled={!texto.trim() || loading}
            className="px-4 py-1.5 bg-indigo-600 text-white text-sm rounded-full disabled:opacity-50 hover:bg-indigo-700"
          >
            {loading ? 'Publicando…' : 'Publicar'}
          </button>
        </div>
      </form>

      {showPhoto && (
        <PhotoUploadModal
          onClose={() => setShowPhoto(false)}
          onUpload={handlePhotoPost}
          loading={loading}
        />
      )}
    </>
  )
}
