interface ConfirmDeleteModalProps {
  onConfirm: () => void
  onCancel: () => void
}

export default function ConfirmDeleteModal({ onConfirm, onCancel }: ConfirmDeleteModalProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white rounded-xl p-6 w-full max-w-xs shadow-xl space-y-4">
        <h2 className="font-semibold text-lg">Excluir publicação</h2>
        <p className="text-sm text-gray-600">Tem certeza de que deseja excluir isso?</p>
        <div className="flex gap-3 justify-end">
          <button
            onClick={onCancel}
            className="px-4 py-1.5 rounded-full text-sm border border-gray-300 hover:bg-gray-50"
          >
            Cancelar
          </button>
          <button
            onClick={onConfirm}
            className="px-4 py-1.5 rounded-full text-sm bg-red-500 text-white hover:bg-red-600"
          >
            Excluir
          </button>
        </div>
      </div>
    </div>
  )
}
