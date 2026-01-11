import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const config = useRuntimeConfig()
  const token = useCookie('auth_token', {
    path: '/',
    watch: true,
    maxAge: 3600 * 24 * 7, // 7 days default
  })

  const user = ref<any>(null)
  const isLoading = ref(false)

  const fetchUser = async () => {
    if (!token.value)
      return

    isLoading.value = true
    try {
      const { data, error } = await useFetch<any>(`${config.public.apiBase}/auth/me`, {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })

      if (error.value) {
        user.value = null
        if (error.value.statusCode === 401) {
          token.value = null
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
      token.value = null
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
      token.value = data.value.data.token
      await fetchUser()
      return data.value.data
    }
  }

  const logout = async () => {
    // If you want to notify backend about logout
    try {
      await useFetch(`${config.public.apiBase}/auth/logout`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      })
    }
    catch {
      // ignore
    }

    token.value = null
    user.value = null
    navigateTo('/login')
  }

  const isLoggedIn = computed(() => !!token.value)

  return {
    user,
    token,
    isLoading,
    isLoggedIn,
    fetchUser,
    login,
    logout,
  }
})
