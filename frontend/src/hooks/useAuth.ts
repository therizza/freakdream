import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/authStore'
import { getMe } from '../api/endpoints'

export function useAuth(required = true) {
  const { user, token, setAuth, clearAuth } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (!token) {
      if (required) navigate('/')
      return
    }
    if (!user) {
      getMe()
        .then((res) => setAuth(res.data, token))
        .catch(() => {
          clearAuth()
          if (required) navigate('/')
        })
    }
  }, [token, user, required, navigate, setAuth, clearAuth])

  return { user, token, clearAuth }
}
