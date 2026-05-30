import { useState, useEffect, useRef } from 'react'
import type { ChatMessage } from '../api/endpoints'
import { getChatHistory } from '../api/endpoints'
import { useWebSocket } from '../hooks/useWebSocket'
import { useAuthStore } from '../store/authStore'
import UserAvatar from './ui/UserAvatar'

interface ChatPanelProps {
  receiverId: number
  receiverNome: string
  receiverFoto?: string
  onClose: () => void
}

export default function ChatPanel({ receiverId, receiverNome, receiverFoto, onClose }: ChatPanelProps) {
  const { user } = useAuthStore()
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState('')
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    getChatHistory(receiverId).then((r) => setMessages([...r.data].reverse()))
  }, [receiverId])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const { sendMessage } = useWebSocket({
    onMessage: (msg) => {
      if (msg.type === 'message' && msg.sender_id === receiverId) {
        setMessages((prev) => [
          ...prev,
          {
            id: Date.now(),
            sender_id: receiverId,
            receiver_id: user?.id ?? 0,
            texto: msg.texto,
            data: msg.data ?? new Date().toISOString(),
          },
        ])
      }
    },
  })

  function handleSend(e: React.FormEvent) {
    e.preventDefault()
    if (!input.trim() || !user) return
    sendMessage(receiverId, input.trim())
    setMessages((prev) => [
      ...prev,
      {
        id: Date.now(),
        sender_id: user.id,
        receiver_id: receiverId,
        texto: input.trim(),
        data: new Date().toISOString(),
      },
    ])
    setInput('')
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 w-80 bg-white rounded-xl shadow-xl flex flex-col overflow-hidden border border-gray-200">
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-3 bg-indigo-600 text-white">
        <UserAvatar foto={receiverFoto} nome={receiverNome} size="sm" />
        <span className="flex-1 font-medium text-sm">{receiverNome}</span>
        <button onClick={onClose} className="text-white/70 hover:text-white">✕</button>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-3 space-y-2 max-h-72">
        {messages.map((m) => (
          <div
            key={m.id}
            className={`flex ${m.sender_id === user?.id ? 'justify-end' : 'justify-start'}`}
          >
            <p
              className={`max-w-[70%] px-3 py-1.5 rounded-2xl text-sm ${
                m.sender_id === user?.id
                  ? 'bg-indigo-600 text-white rounded-br-sm'
                  : 'bg-gray-100 text-gray-800 rounded-bl-sm'
              }`}
            >
              {m.texto}
            </p>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <form onSubmit={handleSend} className="flex gap-2 px-3 py-2 border-t border-gray-100">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Mensagem…"
          className="flex-1 text-sm border border-gray-200 rounded-full px-3 py-1 outline-none focus:ring-2 focus:ring-indigo-300"
        />
        <button
          type="submit"
          className="px-3 py-1 bg-indigo-600 text-white text-sm rounded-full"
        >
          ➤
        </button>
      </form>
    </div>
  )
}
