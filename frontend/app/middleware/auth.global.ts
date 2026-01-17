export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore()

  // Wait for initial fetch if it's happening
  if (import.meta.client && authStore.isLoading) {
    // Usually, the plugin handles this, but we can wait here if needed.
    // However, isLoggedIn is computed, so it will update.
  }

  const guestRoutes = ['/login', '/register', '/forgot-password']
  const isGuestRoute = guestRoutes.includes(to.path)

  if (!authStore.isLoggedIn && !isGuestRoute) {
    return navigateTo('/login')
  }

  if (authStore.isLoggedIn && isGuestRoute) {
    return navigateTo('/')
  }
})
