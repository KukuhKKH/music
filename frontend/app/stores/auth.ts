import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const { $api } = useNuxtApp()
  const user = ref<any>(null)
  const isLoading = ref(false)

  const fetchUser = async () => {
    isLoading.value = true
    try {
      const response = await $api<any>('/auth/me')

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

  const login = async (credentials: { email: string, password: string }) => {
    try {
      const response = await $api<any>('/auth/login', {
        method: 'POST',
        body: credentials,
      })

      if (response) {
        await fetchUser()
        return response.data || true
      }
    }
    catch (error: any) {
      const messages = error.data?.messages
      let errorMessage = 'Login failed'

      if (Array.isArray(messages)) {
        const firstMessage = messages[0]
        if (typeof firstMessage === 'object' && firstMessage !== null) {
          errorMessage = Object.values(firstMessage).join(', ')
        }
        else {
          errorMessage = String(firstMessage)
        }
      }
      else if (typeof messages === 'string') {
        errorMessage = messages
      }

      throw new Error(errorMessage)
    }
  }

  const logout = async () => {
    try {
      await $api('/auth/logout', {
        method: 'POST',
      })
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
