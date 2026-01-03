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
  if (!match)
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
