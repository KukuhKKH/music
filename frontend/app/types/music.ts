export interface Track {
  id: number
  title: string
  artist: string
  album?: string
  duration: number
  file_size: number
  mime_type: string
  public_url: string
  created_at: string
}

export interface MusicMeta {
  current_page: number
  last_page: number
  per_page: number
  total: number
}

export interface MusicResponse {
  data: Track[]
  meta: MusicMeta
}
