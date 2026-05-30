import { useEffect, useRef, useCallback } from 'react'
import type { ChatMessage } from '../api/endpoints'

interface UseWebSocketOptions {
  onMessage: (msg: ChatMessage & { type: string }) => void
  enabled?: boolean
}

export function useWebSocket({ onMessage, enabled = true }: UseWebSocketOptions) {
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!enabled) return

    const token = localStorage.getItem('token')
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const ws = new WebSocket(`${protocol}://${window.location.host}/ws/chat?token=${token}`)
    wsRef.current = ws

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        onMessage(msg)
      } catch {
        // ignore parse errors
      }
    }

    ws.onerror = () => console.warn('WebSocket error')

    return () => {
      ws.close()
    }
  }, [enabled, onMessage])

  const sendMessage = useCallback((receiverId: number, texto: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: 'message', receiver_id: receiverId, texto }))
    }
  }, [])

  return { sendMessage }
}
