import { formatDistanceToNow, parseISO } from 'date-fns'

export function parseSizeToBytes(sizeStr: string): number {
  if (!sizeStr || sizeStr === '0 B')
    return 0
  const units: Record<string, number> = {
    B: 1,
    KB: 1024,
    MB: 1024 * 1024,
    GB: 1024 * 1024 * 1024,
    TB: 1024 * 1024 * 1024 * 1024,
  }
  const match = sizeStr.match(/^(\d+(?:\.\d+)?)\s*([A-Z]+)$/i)
  if (!match || !match[1] || !match[2])
    return 0
  const value = Number.parseFloat(match[1])
  const unit = match[2].toUpperCase()
  return value * (units[unit] || 1)
}

export function formatTimeAgo(dateStr: string | null | undefined): string {
  if (!dateStr || dateStr === '-')
    return '-'
  try {
    return formatDistanceToNow(parseISO(dateStr), { addSuffix: true })
  }
  catch {
    return dateStr
  }
}

export function formatDuration(seconds: number): string {
  if (!seconds || Number.isNaN(seconds))
    return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0)
    return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${Number.parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`
}
