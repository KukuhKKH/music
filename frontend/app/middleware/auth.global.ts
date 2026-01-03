export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  // Guest-only routes (login, register, forgot-password)
  const guestRoutes = ['/login', '/register', '/forgot-password']

  if (!authStore.isLoggedIn && !guestRoutes.includes(to.path)) {
    return navigateTo('/login')
  }

  if (authStore.isLoggedIn && guestRoutes.includes(to.path)) {
    return navigateTo('/')
  }
})
