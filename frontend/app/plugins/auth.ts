export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()

  if (import.meta.client || import.meta.server) {
    await authStore.fetchUser()
  }
})
