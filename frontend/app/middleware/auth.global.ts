export default defineNuxtRouteMiddleware((to) => {
  const { isLoggedIn } = useAuth()

  // Guest-only routes (login, register, forgot-password)
  const guestRoutes = ['/login', '/register', '/forgot-password']

  if (!isLoggedIn.value && !guestRoutes.includes(to.path)) {
    return navigateTo('/login')
  }

  if (isLoggedIn.value && guestRoutes.includes(to.path)) {
    return navigateTo('/')
  }
})
