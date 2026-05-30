import { useState, useEffect, useCallback } from 'react'
import { useAuth } from '../../hooks/useAuth'
import Navbar from '../../components/ui/Navbar'
import PostForm from '../../components/ui/PostForm'
import PostCard from '../../components/ui/PostCard'
import { getFeed } from '../../api/endpoints'
import type { Post } from '../../api/endpoints'

export default function FeedPage() {
  useAuth()
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [offset, setOffset] = useState(0)
  const [hasMore, setHasMore] = useState(true)

  const loadFeed = useCallback(async (reset = false) => {
    const o = reset ? 0 : offset
    setLoading(true)
    try {
      const res = await getFeed(o)
      const newPosts = res.data
      if (reset) {
        setPosts(newPosts)
        setOffset(newPosts.length)
      } else {
        setPosts((prev) => [...prev, ...newPosts])
        setOffset((prev) => prev + newPosts.length)
      }
      setHasMore(newPosts.length === 20)
    } finally {
      setLoading(false)
    }
  }, [offset])

  useEffect(() => {
    loadFeed(true)
    const handler = () => loadFeed(true)
    window.addEventListener('feed:refresh', handler)
    return () => window.removeEventListener('feed:refresh', handler)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  function handlePosted(post: Post) {
    setPosts((prev) => [post, ...prev])
  }

  function handleDeleted(postId: number) {
    setPosts((prev) => prev.filter((p) => p.post_id !== postId))
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-xl mx-auto px-4 py-6 space-y-4">
        <PostForm onPosted={handlePosted} />

        {loading && posts.length === 0 && (
          <p className="text-center text-sm text-gray-400 py-8">Carregando…</p>
        )}

        {posts.length === 0 && !loading && (
          <p className="text-center text-sm text-gray-400 py-8">
            Nenhuma publicação ainda. Siga pessoas para ver o feed!
          </p>
        )}

        {posts.map((post) => (
          <PostCard key={post.post_id} post={post} onDeleted={handleDeleted} />
        ))}

        {hasMore && !loading && (
          <button
            onClick={() => loadFeed(false)}
            className="w-full py-2 text-sm text-indigo-600 hover:underline"
          >
            Carregar mais
          </button>
        )}

        {loading && posts.length > 0 && (
          <p className="text-center text-sm text-gray-400">Carregando…</p>
        )}
      </main>
    </div>
  )
}
