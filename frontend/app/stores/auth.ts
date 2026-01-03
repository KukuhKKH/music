import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const config = useRuntimeConfig()
  const user = ref<any>(null)
  const isLoading = ref(false)

  // fetchUser will now rely on the HttpOnly cookie handled by the browser/proxy
  const fetchUser = async () => {
    isLoading.value = true
    try {
      const { data, error } = await useFetch<any>(`${config.public.apiBase}/auth/me`, {
        // useFetch automatically includes cookies if same-origin (proxied)
        // No manual headers needed for token
      })

      if (error.value) {
        user.value = null
        if (error.value.statusCode === 401) {
          // If we are on a protected route, middleware will handle redirect
        }
        return
      }

      if (data.value && data.value.data) {
        user.value = data.value.data
      }
    }

    catch (err) {
      console.error('Failed to fetch user:', err)
      user.value = null
    }

    finally {
      isLoading.value = false
    }
  }

  const login = async (credentials: { email: string, password: string }) => {
    const { data, error } = await useFetch<any>(`${config.public.apiBase}/auth/login`, {
      method: 'POST',
      body: credentials,
    })

    if (error.value) {
      const messages = error.value.data?.messages
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

    if (data.value && data.value.data) {
      await fetchUser()
      return data.value.data
    }
  }

  const logout = async () => {
    try {
      await useFetch(`${config.public.apiBase}/auth/logout`, { method: 'POST' })
    }
    catch {
      // ignored
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
