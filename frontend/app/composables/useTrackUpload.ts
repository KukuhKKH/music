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
  const config = useRuntimeConfig()

  const uploadTrack = (options: UploadOptions): Promise<any> => {
    return new Promise((resolve, reject) => {
      const formData = new FormData()
      formData.append('file', options.file)
      formData.append('title', options.metadata.title)
      formData.append('artist', options.metadata.artist)
      formData.append('album', options.metadata.album)
      formData.append('duration', options.metadata.duration.toString())

      const xhr = new XMLHttpRequest()
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
        if (xhr.status >= 200 && xhr.status < 300) {
          resolve(xhr.response)
        } else {
          reject(new Error(`Upload failed with status ${xhr.status}`))
        }
      }

      xhr.onerror = () => reject(new Error('Network error during upload'))
      xhr.onabort = () => reject(new Error('Upload aborted'))

      xhr.send(formData)
    })
  }

  return {
    uploadTrack,
  }
}
