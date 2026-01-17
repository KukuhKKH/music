import { toast } from 'vue-sonner'

export default defineNuxtPlugin((_nuxtApp) => {
  const config = useRuntimeConfig()

  // Use Nuxt event to get headers for SSR cookie forwarding
  const event = useRequestEvent()

  const apiFetch = $fetch.create({
    baseURL: config.public.apiBase as string,
    onRequest({ options }) {
      options.credentials = 'include'

      if (import.meta.server && event) {
        const cookies = getHeader(event, 'cookie')
        if (cookies) {
          options.headers = {
            ...options.headers,
            cookie: cookies,
          } as any
        }
      }
    },
    onResponseError({ response }) {
      if (response.status === 401) {
        const authStore = useAuthStore()
        authStore.user = null

        // Only redirect on client side to avoid infinite loops during SSR
        if (import.meta.client) {
          const currentPath = useRoute().path
          if (currentPath !== '/login') {
            toast.error('Session expired, please login again.')
            navigateTo('/login')
          }
        }
      }
    },
  })

  // Provide $api as a global helper
  return {
    provide: {
      api: apiFetch,
    },
  }
})
