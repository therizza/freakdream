import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import Navbar from '../../components/ui/Navbar'
import UserAvatar from '../../components/ui/UserAvatar'
import PostCard from '../../components/ui/PostCard'
import ChatPanel from '../../components/ChatPanel'
import {
  getUser,
  getUserPosts,
  getFollowCount,
  getFollowStatus,
  follow,
  unfollow,
  sendFriendRequest,
} from '../../api/endpoints'
import type { User, Post } from '../../api/endpoints'

export default function ProfilePage() {
  const { user: me } = useAuth()
  const { id } = useParams<{ id: string }>()
  const userId = Number(id)

  const [profile, setProfile] = useState<User | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [followerCount, setFollowerCount] = useState(0)
  const [isFollowing, setIsFollowing] = useState(false)
  const [loading, setLoading] = useState(true)
  const [showChat, setShowChat] = useState(false)

  useEffect(() => {
    if (!userId) return
    Promise.all([
      getUser(userId),
      getUserPosts(userId),
      getFollowCount(userId),
      getFollowStatus(userId),
    ]).then(([userRes, postsRes, countRes, statusRes]) => {
      setProfile(userRes.data)
      setPosts(postsRes.data)
      setFollowerCount(countRes.data.count)
      setIsFollowing(statusRes.data.following)
    }).finally(() => setLoading(false))
  }, [userId])

  async function handleFollow() {
    if (isFollowing) {
      await unfollow(userId)
      setIsFollowing(false)
      setFollowerCount((c) => c - 1)
    } else {
      await follow(userId)
      setIsFollowing(true)
      setFollowerCount((c) => c + 1)
    }
  }

  async function handleFriendRequest() {
    await sendFriendRequest(userId)
  }

  function handleDeleted(postId: number) {
    setPosts((prev) => prev.filter((p) => p.post_id !== postId))
  }

  const isOwnProfile = me?.id === userId

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <p className="text-center py-16 text-gray-400">Carregando…</p>
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <p className="text-center py-16 text-gray-400">Usuário não encontrado.</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <main className="max-w-xl mx-auto px-4 py-6 space-y-4">
        {/* Profile card */}
        <div className="bg-white rounded-xl shadow p-6 flex items-center gap-5">
          <UserAvatar foto={profile.foto} nome={profile.nome} size="lg" />
          <div className="flex-1 min-w-0">
            <h1 className="text-xl font-bold text-gray-900">{profile.nome}</h1>
            {profile.sobre && <p className="text-sm text-gray-500 mt-0.5">{profile.sobre}</p>}
            <p className="text-sm text-gray-400 mt-1">{followerCount} seguidores</p>
          </div>
          {!isOwnProfile && (
            <div className="flex flex-col gap-2">
              <button
                onClick={handleFollow}
                className={`px-4 py-1.5 rounded-full text-sm font-medium ${
                  isFollowing
                    ? 'border border-indigo-400 text-indigo-600 hover:bg-indigo-50'
                    : 'bg-indigo-600 text-white hover:bg-indigo-700'
                }`}
              >
                {isFollowing ? 'Seguindo' : 'Seguir'}
              </button>
              <button
                onClick={handleFriendRequest}
                className="px-4 py-1.5 rounded-full text-sm border border-gray-300 text-gray-600 hover:bg-gray-50"
              >
                + Amigo
              </button>
              <button
                onClick={() => setShowChat(true)}
                className="px-4 py-1.5 rounded-full text-sm border border-gray-300 text-gray-600 hover:bg-gray-50"
              >
                💬 Chat
              </button>
            </div>
          )}
        </div>

        {/* Posts */}
        {posts.length === 0 && (
          <p className="text-center text-sm text-gray-400 py-4">Nenhuma publicação ainda.</p>
        )}
        {posts.map((post) => (
          <PostCard key={post.post_id} post={post} onDeleted={handleDeleted} />
        ))}
      </main>

      {showChat && profile && (
        <ChatPanel
          receiverId={profile.id}
          receiverNome={profile.nome}
          receiverFoto={profile.foto}
          onClose={() => setShowChat(false)}
        />
      )}
    </div>
  )
}
