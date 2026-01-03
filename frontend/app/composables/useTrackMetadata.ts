import { Buffer } from 'buffer'
import * as musicMetadata from 'music-metadata-browser'

// Ensure Buffer is available for music-metadata-browser in browser
if (typeof window !== 'undefined' && !window.Buffer) {
  (window as any).Buffer = Buffer
}

export interface TrackMetadata {
  title: string
  artist: string
  album: string
  duration: number
}

export function useTrackMetadata() {
  const extractMetadata = async (file: File): Promise<TrackMetadata> => {
    const metadata: TrackMetadata = {
      title: file.name.replace(/\.[^/.]+$/, ''),
      artist: '',
      album: '',
      duration: 0,
    }

    try {
      // 1. Primary: music-metadata-browser
      const parsed = await musicMetadata.parseBlob(file)
      const common = parsed.common

      if (common.title) metadata.title = common.title
      metadata.artist = common.artist || common.albumartist || ''
      metadata.album = common.album || ''
      metadata.duration = Math.round(parsed.format.duration || 0)

      // 2. Fallback for Duration: Browser Audio API
      if (metadata.duration === 0) {
        metadata.duration = await getDurationFromAudioApi(file)
      }

      // 3. Fallback for Artist/Album: Filename parsing
      if (!metadata.artist && file.name.includes(' - ')) {
        const parts = file.name.replace(/\.[^/.]+$/, '').split(' - ')
        if (parts.length >= 2) {
          metadata.artist = parts[0].trim().replace(/_/g, ' ')
          if (parts.length >= 3 && !metadata.album) {
            metadata.album = parts[1].trim().replace(/_/g, ' ')
          }
        }
      }

      // 4. Final cleaning
      metadata.title = cleanMetadataString(metadata.title)
      metadata.artist = cleanMetadataString(metadata.artist)
      metadata.album = cleanMetadataString(metadata.album)

      return metadata
    } catch (err) {
      console.warn('Metadata extraction failed, using defaults:', err)
      return metadata
    }
  }

  const getDurationFromAudioApi = (file: File): Promise<number> => {
    return new Promise((resolve) => {
      try {
        const url = URL.createObjectURL(file)
        const audio = new Audio()
        audio.src = url
        audio.onloadedmetadata = () => {
          const duration = Math.round(audio.duration || 0)
          URL.revokeObjectURL(url)
          resolve(Number.isNaN(duration) ? 0 : duration)
        }
        audio.onerror = () => {
          URL.revokeObjectURL(url)
          resolve(0)
        }
        // Safety timeout
        setTimeout(() => {
          URL.revokeObjectURL(url)
          resolve(0)
        }, 3000)
      } catch {
        resolve(0)
      }
    })
  }

  const cleanMetadataString = (str: string): string => {
    return str.replace(/_/g, ' ').replace(/\s+/g, ' ').trim()
  }

  return {
    extractMetadata,
  }
}
