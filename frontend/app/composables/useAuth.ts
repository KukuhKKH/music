export function useAuth() {
  const authStore = useAuthStore()

  return {
    token: computed(() => authStore.token),
    user: computed(() => authStore.user),
    login: authStore.login,
    logout: authStore.logout,
    fetchUser: authStore.fetchUser,
    isLoggedIn: computed(() => authStore.isLoggedIn),
    isLoading: computed(() => authStore.isLoading),
  }
}
