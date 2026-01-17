import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const { $api } = useNuxtApp()
  const user = ref<any>(null)
  const isLoading = ref(true)

  const fetchUser = async () => {
    console.log('🔄 Fetching user session...')
    isLoading.value = true
    try {
      const response = await $api<any>('/auth/me')
      // Small artificial delay for premium feel and to prevent flicker
      await new Promise(resolve => setTimeout(resolve, 800))

      if (response && response.data) {
        user.value = response.data
      }
    }
    catch {
      user.value = null
    }
    finally {
      isLoading.value = false
    }
  }

  const login = () => {
    const config = useRuntimeConfig()
    // Directly redirect to backend login endpoint
    window.location.href = `${config.public.apiBase}/auth/login`
  }

  const logout = async () => {
    try {
      const response = await $api<any>('/auth/logout', {
        method: 'POST',
      })

      // If backend provides a global logout URL (Logto), redirect there
      if (response && response.data && response.data.logout_url) {
        user.value = null
        window.location.href = response.data.logout_url
        return
      }
    }
    catch {
      // ignore
    }

    user.value = null
    navigateTo('/login')
  }

  const isLoggedIn = computed(() => !!user.value)

  return {
    user,
    isLoading,
    isLoggedIn,
    fetchUser,
    login,
    logout,
  }
})
