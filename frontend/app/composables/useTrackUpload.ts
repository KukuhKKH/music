import { toast } from 'vue-sonner'

export interface UploadOptions {
  url: string
  file: File
  metadata: {
    title: string
    artist: string
    album: string
    duration: number
  }
  onProgress?: (progress: number) => void
}

export function useTrackUpload() {
  const uploadTrack = (options: UploadOptions): Promise<any> => {
    return new Promise((resolve, reject) => {
      const formData = new FormData()
      formData.append('file', options.file)
      formData.append('title', options.metadata.title)
      formData.append('artist', options.metadata.artist)
      formData.append('album', options.metadata.album)
      formData.append('duration', options.metadata.duration.toString())

      const xhr = new XMLHttpRequest()

      let slowNetworkToastId: string | number | null = null
      const slowNetworkTimeout = setTimeout(() => {
        slowNetworkToastId = toast.warning('Sabar dude, internetmu koyo babi 🐷', {
          description: 'Upload lagu lagi proses, ojok dicancel yo.',
          duration: 10000,
        })
      }, 15000)

      const cleanup = () => {
        if (slowNetworkTimeout)
          clearTimeout(slowNetworkTimeout)
        if (slowNetworkToastId)
          toast.dismiss(slowNetworkToastId)
      }

      xhr.open('POST', options.url)
      xhr.withCredentials = true

      if (xhr.upload && options.onProgress) {
        xhr.upload.onprogress = (event) => {
          if (event.lengthComputable) {
            const progress = Math.round((event.loaded / event.total) * 100)
            options.onProgress?.(progress)
          }
        }
      }

      xhr.onload = () => {
        cleanup()
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve(xhr.response)
        }
        else {
          reject(new Error(`Upload failed with status ${xhr.status}`))
        }
      }

      xhr.onerror = () => {
        cleanup()
        reject(new Error('Network error during upload'))
      }
      xhr.onabort = () => {
        cleanup()
        reject(new Error('Upload aborted'))
      }

      xhr.send(formData)
    })
  }

  return {
    uploadTrack,
  }
}
