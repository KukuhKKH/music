export const useAuth = () => {
  const config = useRuntimeConfig()
  const token = useCookie('auth_token')
  const user = useState('user', () => null)

  const login = async (email: string, password: string) => {
    try {
      const { data, error } = await useFetch<any>(`${config.public.apiBase}/auth/login`, {
        method: 'POST',
        body: { email, password }
      })

      if (error.value) {
        const messages = error.value.data?.messages
        let errorMessage = 'Login failed'

        if (Array.isArray(messages)) {
          const firstMessage = messages[0]
          if (typeof firstMessage === 'object' && firstMessage !== null) {
            errorMessage = Object.values(firstMessage).join(', ')
          } else {
            errorMessage = String(firstMessage)
          }
        } else if (typeof messages === 'string') {
          errorMessage = messages
        }

        throw new Error(errorMessage)
      }

      if (data.value && data.value.data) {
        token.value = data.value.data.token
        await fetchUser()
        return data.value.data
      }
    } catch (err) {
      throw err
    }
  }

  const fetchUser = async () => {
    if (!token.value) return

    try {
      const { data, error } = await useFetch<any>(`${config.public.apiBase}/auth/me`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      if (error.value) {
        user.value = null
        token.value = null
        return
      }

      if (data.value && data.value.data) {
        user.value = data.value.data
      }
    } catch (err) {
      console.error('Failed to fetch user:', err)
      user.value = null
      token.value = null
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    navigateTo('/login')
  }

  const isLoggedIn = computed(() => !!token.value)

  return {
    token,
    user,
    login,
    logout,
    fetchUser,
    isLoggedIn
  }
}
