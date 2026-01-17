<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";

definePageMeta({
  layout: "blank",
});

const auth = useAuthStore();

// Watch for auth status to redirect immediately if user is found
watchEffect(() => {
  if (auth.isLoggedIn) {
    navigateTo("/");
  }
});
</script>

<template>
  <LayoutAuth reverse>
    <div class="relative w-full max-w-sm mx-auto">
      <!-- Transition Container -->
      <Transition name="fade-slide" mode="out-in">
        <!-- Loading State: If still checking session -->
        <div
          v-if="auth.isLoading"
          key="loading"
          class="flex flex-col items-center justify-center space-y-6 py-10"
        >
          <div class="relative">
            <!-- Animated Glow Ring -->
            <div
              class="absolute -inset-4 bg-primary/20 rounded-full blur-xl animate-pulse"
            />
            <div
              class="relative bg-background border border-border p-4 rounded-full shadow-2xl"
            >
              <Loader2 class="size-12 animate-spin text-primary" />
            </div>
          </div>
          <div class="text-center animate-pulse">
            <h2 class="text-xl font-medium tracking-tight">
              Securing Session...
            </h2>
            <p class="text-sm text-muted-foreground mt-2">
              Checking your authorization status
            </p>
          </div>
        </div>

        <!-- Authenticated State: Briefly show welcome before redirect -->
        <div
          v-else-if="auth.isLoggedIn"
          key="authenticated"
          class="flex flex-col items-center justify-center space-y-4 py-8 text-center text-primary"
        >
          <div class="bg-primary/10 p-3 rounded-full">
            <Icon name="mdi:shield-check" class="size-10" />
          </div>
          <h2 class="text-2xl font-bold">Welcome Back!</h2>
          <p class="text-muted-foreground">Redirecting to your dashboard...</p>
        </div>

        <!-- Login Form State: If definitely not logged in -->
        <div v-else key="login" class="grid gap-6">
          <div class="grid gap-2 text-center">
            <div class="flex justify-center mb-2">
              <div
                class="p-3 bg-primary/5 rounded-2xl border border-primary/10 shadow-sm"
              >
                <Icon name="mdi:music-circle" class="size-8 text-primary" />
              </div>
            </div>
            <h1 class="text-2xl font-semibold tracking-tight">Music Portal</h1>
            <p class="text-balance text-sm text-muted-foreground">
              Vibe with the rhythm of simplicity.
            </p>
          </div>

          <div
            class="p-1 rounded-3xl bg-gradient-to-b from-primary/10 to-transparent"
          >
            <div
              class="bg-background border border-border/50 rounded-[22px] p-6 shadow-xl backdrop-blur-sm"
            >
              <AuthSignIn />
            </div>
          </div>

          <div class="text-center">
            <p
              class="text-[10px] uppercase tracking-widest text-muted-foreground/50 font-bold"
            >
              Protected by BangLipai Cryptography
            </p>
          </div>
        </div>
      </Transition>
    </div>
  </LayoutAuth>
</template>

<style scoped>
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Custom glow effect */
.bg-background {
  backdrop-filter: blur(16px);
}
</style>
