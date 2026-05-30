import api from './client'

export interface User {
  id: number
  nome: string
  sobre: string
  email?: string
  foto?: string
  created_at?: string
}

export interface Post {
  post_id: number
  user_id: number
  tipo: 1 | 2 | 3
  texto: string
  data: string
  user?: User
}

export interface Notification {
  id: number
  user_id: number
  tipo: string
  texto: string
  data: string
  lida: boolean
}

export interface FriendRequest {
  user_id: number
  amigo_id: number
  status: 'pending' | 'accepted'
  criado_em: string
}

export interface ChatMessage {
  id: number
  sender_id: number
  receiver_id: number
  texto: string
  data: string
}

// Auth
export const login = (email: string, senha: string) =>
  api.post<{ token: string; user: User }>('/auth/login', { email, senha })

export const register = (nome: string, sobre: string, email: string, senha: string) =>
  api.post<{ token: string; user: User }>('/auth/register', { nome, sobre, email, senha })

export const logout = () => api.post('/auth/logout')

// Users
export const getMe = () => api.get<User>('/users/me')
export const getUser = (id: number) => api.get<User>(`/users/${id}`)
export const searchUsers = (q: string) => api.get<User[]>(`/users/search?q=${encodeURIComponent(q)}`)
export const getUserPosts = (id: number, offset = 0) =>
  api.get<Post[]>(`/users/${id}/posts?limit=20&offset=${offset}`)

// Feed & Posts
export const getFeed = (offset = 0) => api.get<Post[]>(`/feed?limit=20&offset=${offset}`)
export const createPost = (texto: string) => api.post<Post>('/posts', { texto })
export const deletePost = (id: number) => api.delete(`/posts/${id}`)
export const sharePost = (post_id: number) => api.post<Post>('/posts/share', { post_id })

export const uploadPhoto = (formData: FormData) =>
  api.post<{ post_id: number; foto_id: number; foto: string }>('/posts/photo', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })

// Follow
export const follow = (id: number) => api.post(`/follow/${id}`)
export const unfollow = (id: number) => api.delete(`/follow/${id}`)
export const getFollowCount = (id: number) => api.get<{ count: number }>(`/follow/${id}/count`)
export const getFollowStatus = (id: number) => api.get<{ following: boolean }>(`/follow/${id}/status`)

// Friends
export const sendFriendRequest = (id: number) => api.post(`/friends/request/${id}`)
export const acceptFriendRequest = (id: number) => api.post(`/friends/accept/${id}`)
export const removeFriend = (id: number) => api.delete(`/friends/${id}`)
export const getFriendRequests = () => api.get<FriendRequest[]>('/friends/requests')

// Notifications
export const getNotifications = () => api.get<Notification[]>('/notifications')
export const removeNotification = (id: number) => api.delete(`/notifications/${id}`)
export const markNotificationRead = (id: number) => api.patch(`/notifications/${id}/read`)

// Chat history
export const getChatHistory = (id: number) => api.get<ChatMessage[]>(`/chat/${id}/history`)
