export default defineNuxtPlugin(() => {
  // Override global $fetch to always include credentials (important for cross-origin cookies)
  const originalFetch = globalThis.$fetch
  if (originalFetch && typeof originalFetch.create === 'function') {
    globalThis.$fetch = originalFetch.create({
      onRequest({ options }) {
        options.credentials = options.credentials || 'include'
      },
    })
  }

  return {
    provide: {
      api: globalThis.$fetch,
    },
  }
})
