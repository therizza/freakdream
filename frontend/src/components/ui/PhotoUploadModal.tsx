import { useRef, useState } from 'react'

interface PhotoUploadModalProps {
  onClose: () => void
  onUpload: (fd: FormData) => void
  loading: boolean
}

export default function PhotoUploadModal({ onClose, onUpload, loading }: PhotoUploadModalProps) {
  const [preview, setPreview] = useState<string | null>(null)
  const [texto, setTexto] = useState('')
  const [perfil, setPerfil] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  function handleFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file || !file.type.startsWith('image/')) return
    const reader = new FileReader()
    reader.onload = () => {
      if (typeof reader.result === 'string' && reader.result.startsWith('data:image/')) {
        setPreview(reader.result)
      }
    }
    reader.readAsDataURL(file)
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const file = inputRef.current?.files?.[0]
    if (!file) return
    const fd = new FormData()
    fd.append('imagem', file)
    fd.append('texto', texto)
    fd.append('perfil', perfil ? '1' : '0')
    onUpload(fd)
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white rounded-xl w-full max-w-sm p-6 space-y-4 shadow-xl">
        <div className="flex justify-between items-center">
          <h2 className="font-semibold text-lg">Enviar Imagem</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">✕</button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-3">
          <label className="block cursor-pointer border-2 border-dashed border-gray-200 rounded-lg p-6 text-center hover:border-indigo-400">
            {preview ? (
              <img src={preview} alt="preview" className="mx-auto max-h-48 rounded-lg object-cover" />
            ) : (
              <span className="text-gray-400 text-sm">Clique para selecionar</span>
            )}
            <input
              ref={inputRef}
              type="file"
              accept="image/*"
              className="hidden"
              onChange={handleFile}
            />
          </label>

          <textarea
            value={texto}
            onChange={(e) => setTexto(e.target.value)}
            placeholder="Descrição (opcional)"
            rows={2}
            className="w-full resize-none border border-gray-200 rounded-lg p-2 text-sm outline-none"
          />

          <label className="flex items-center gap-2 text-sm text-gray-600">
            <input type="checkbox" checked={perfil} onChange={(e) => setPerfil(e.target.checked)} />
            Usar como foto de perfil
          </label>

          <button
            type="submit"
            disabled={!preview || loading}
            className="w-full py-2 bg-indigo-600 text-white rounded-full text-sm font-medium disabled:opacity-50 hover:bg-indigo-700"
          >
            {loading ? 'Enviando…' : 'Publicar'}
          </button>
        </form>
      </div>
    </div>
  )
}
