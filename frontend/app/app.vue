<script setup lang="ts">
import { ConfigProvider } from "reka-ui";
import { Toaster } from "@/components/ui/sonner";
import "vue-sonner/style.css";

const colorMode = useColorMode();
const color = computed(() =>
  colorMode.value === "dark" ? "#09090b" : "#ffffff",
);
const { theme } = useAppSettings();
const dir = ref<any>("ltr");

useHead({
  meta: [
    { charset: "utf-8" },
    { name: "viewport", content: "width=device-width, initial-scale=1" },
    { key: "theme-color", name: "theme-color", content: color },
  ],
  link: [{ rel: "icon", href: "/favicon.ico" }],
  htmlAttrs: {
    lang: "en",
  },
  bodyAttrs: {
    class: computed(
      () =>
        `color-${theme.value?.color || "default"} theme-${
          theme.value?.type || "default"
        }`,
    ),
  },
});

const title = "Music";
const description = "Music App";

useSeoMeta({
  title,
  description,
  ogTitle: title,
  ogDescription: description,
  ogUrl: "https://banglipai.tech",
  ogImage: "https://banglipai.tech",
});

const auth = useAuthStore();
const showLoader = ref(true);

onMounted(() => {
  // Even if auth is fast, keep loader for a bit to avoid flicker
  watch(
    () => auth.isLoading,
    (loading) => {
      if (!loading) {
        setTimeout(() => {
          showLoader.value = false;
        }, 500);
      }
    },
    { immediate: true },
  );
});
</script>

<template>
  <Body class="overscroll-none antialiased bg-background text-foreground">
    <ConfigProvider :dir="dir">
      <div id="app" vaul-drawer-wrapper class="relative min-h-screen">
        <Transition name="app-fade" mode="out-in">
          <div
            v-if="showLoader"
            key="app-loading"
            class="fixed inset-0 z-[9999] flex flex-col items-center justify-center bg-background"
          >
            <div class="relative mb-8">
              <!-- Triple Ring Pulse -->
              <div
                class="absolute inset-0 scale-150 bg-primary/20 rounded-full blur-2xl animate-pulse"
              />
              <div
                class="absolute inset-0 scale-125 bg-primary/10 rounded-full blur-xl animate-pulse delay-75"
              />
              <div
                class="relative bg-background border border-border/50 p-6 rounded-3xl shadow-2xl backdrop-blur-xl"
              >
                <Icon
                  name="mdi:music-note-eighth"
                  class="size-16 text-primary animate-bounce-slow"
                />
              </div>
            </div>

            <div class="space-y-3 text-center">
              <h1
                class="text-3xl font-bold tracking-tighter sm:text-4xl animate-pulse"
              >
                Music Portal
              </h1>
              <div
                class="flex items-center justify-center gap-2 text-muted-foreground"
              >
                <div
                  class="h-1 w-12 bg-primary/30 rounded-full overflow-hidden"
                >
                  <div class="h-full bg-primary animate-loading-bar" />
                </div>
                <span class="text-xs font-mono uppercase tracking-[0.2em]"
                  >Synchronizing</span
                >
                <div
                  class="h-1 w-12 bg-primary/30 rounded-full overflow-hidden text-right"
                >
                  <div class="h-full bg-primary animate-loading-bar-rev" />
                </div>
              </div>
            </div>
          </div>

          <div v-else key="app-content">
            <NuxtLayout>
              <NuxtPage />
              <MusicMiniPlayer />
            </NuxtLayout>
          </div>
        </Transition>

        <AppSettings />
      </div>

      <Toaster
        :theme="(colorMode.preference as any) || 'system'"
        position="top-right"
        rich-colors
      />
    </ConfigProvider>
  </Body>
</template>

<style>
.app-fade-enter-active,
.app-fade-leave-active {
  transition: opacity 0.5s ease;
}

.app-fade-enter-from,
.app-fade-leave-to {
  opacity: 0;
}

@keyframes loading-bar {
  0% {
    transform: translateX(-100%);
  }
  50% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(100%);
  }
}

@keyframes loading-bar-rev {
  0% {
    transform: translateX(100%);
  }
  50% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(-100%);
  }
}

.animate-loading-bar {
  animation: loading-bar 2s infinite ease-in-out;
}

.animate-loading-bar-rev {
  animation: loading-bar-rev 2s infinite ease-in-out;
}

.animate-bounce-slow {
  animation: bounce 3s infinite;
}

@keyframes bounce {
  0%,
  100% {
    transform: translateY(-5%);
    animation-timing-function: cubic-bezier(0.8, 0, 1, 1);
  }
  50% {
    transform: translateY(0);
    animation-timing-function: cubic-bezier(0, 0, 0.2, 1);
  }
}
</style>
