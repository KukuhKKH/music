import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Track } from '~/types/music'

export const usePlayerStore = defineStore('player', () => {
  const currentTrack = ref<Track | null>(null)
  const isPlaying = ref(false)
  const volume = ref(0.7)
  const currentTime = ref(0)
  const duration = ref(0)
  const audio = ref<HTMLAudioElement | null>(null)

  function play(track: Track) {
    if (currentTrack.value?.id === track.id) {
      resume()
      return
    }

    if (audio.value) {
      audio.value.pause()
      audio.value = null
    }

    currentTrack.value = track
    audio.value = new Audio(track.public_url)
    audio.value.volume = volume.value
    duration.value = track.duration

    audio.value.addEventListener('loadedmetadata', () => {
      if (audio.value)
        duration.value = audio.value.duration
    })

    audio.value.addEventListener('timeupdate', () => {
      if (audio.value)
        currentTime.value = audio.value.currentTime
    })

    audio.value.addEventListener('ended', () => {
      isPlaying.value = false
      currentTime.value = 0
    })

    audio.value.play()
    isPlaying.value = true
  }

  function pause() {
    if (audio.value) {
      audio.value.pause()
      isPlaying.value = false
    }
  }

  function resume() {
    if (audio.value && !isPlaying.value) {
      audio.value.play()
      isPlaying.value = true
    }
  }

  function toggle(track: Track) {
    if (currentTrack.value?.id === track.id) {
      if (isPlaying.value)
        pause()
      else
        resume()
    }
    else {
      play(track)
    }
  }

  function seek(time: number) {
    if (audio.value) {
      audio.value.currentTime = time
      currentTime.value = time
    }
  }

  function skip(seconds: number) {
    if (audio.value) {
      const newTime = audio.value.currentTime + seconds
      seek(Math.max(0, Math.min(newTime, duration.value)))
    }
  }

  function setVolume(val: number) {
    volume.value = val
    if (audio.value)
      audio.value.volume = val
  }

  function stop() {
    if (audio.value) {
      audio.value.pause()
      audio.value = null
    }
    currentTrack.value = null
    isPlaying.value = false
    currentTime.value = 0
  }

  return {
    currentTrack,
    isPlaying,
    volume,
    currentTime,
    duration,
    play,
    pause,
    resume,
    toggle,
    seek,
    skip,
    setVolume,
    stop,
  }
})
