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
  limit: number
  page: number
  next_page: number
  previous_page: number
  count: number
  total_page: number
}

export interface MusicResponse {
  data: Track[]
  meta: MusicMeta
}
