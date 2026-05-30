interface UserAvatarProps {
  foto?: string
  nome?: string
  size?: 'sm' | 'md' | 'lg'
}

const sizes = { sm: 'h-8 w-8', md: 'h-10 w-10', lg: 'h-16 w-16' }

export default function UserAvatar({ foto, nome, size = 'md' }: UserAvatarProps) {
  const initial = nome ? nome.charAt(0).toUpperCase() : '?'

  if (foto) {
    return (
      <img
        src={`/uploads/${foto}`}
        alt={nome || 'avatar'}
        className={`${sizes[size]} rounded-full object-cover ring-2 ring-white`}
        onError={(e) => {
          e.currentTarget.style.display = 'none'
        }}
      />
    )
  }

  return (
    <div
      className={`${sizes[size]} rounded-full bg-indigo-500 flex items-center justify-center text-white font-semibold select-none`}
    >
      {initial}
    </div>
  )
}
