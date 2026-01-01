import tailwindcss from '@tailwindcss/vite'
// https://nuxt.com/docs/api/configuration/nuxt-config
const apiProxyPath = process.env.NUXT_API_PROXY_PATH || '/backend'
const backendUrl = process.env.NUXT_BACKEND_URL || 'http://localhost:8080'
const devHost = process.env.NUXT_DEV_HOST || 'app-music.local'

export default defineNuxtConfig({
  devtools: { enabled: true },

  css: ['~/assets/css/tailwind.css'],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  devServer: {
    host: devHost,
    port: 3000,
  },

  components: [
    {
      path: '~/components',
      extensions: ['.vue'],
    },
  ],

  modules: [
    'shadcn-nuxt',
    '@vueuse/nuxt',
    '@nuxt/eslint',
    '@nuxt/icon',
    '@pinia/nuxt',
    '@nuxtjs/color-mode',
    '@nuxt/fonts',
  ],

  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: '',
    /**
     * Directory that the component lives in.
     * @default "~/components/ui"
     */
    componentDir: '~/components/ui',
  },

  colorMode: {
    classSuffix: '',
  },

  eslint: {
    config: {
      standalone: false,
    },
  },

  fonts: {
    defaults: {
      weights: [300, 400, 500, 600, 700, 800],
    },
  },

  routeRules: {
    '/settings': { redirect: '/settings/profile' },
    [`${apiProxyPath}/**`]: { proxy: `${backendUrl}/**` },
  },

  imports: {
    dirs: [
      './lib',
    ],
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || apiProxyPath,
    },
  },

  compatibilityDate: '2025-12-31',
})
