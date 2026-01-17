<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
import { toast } from "vue-sonner";

const isLoading = ref(false);
const { login } = useAuth();

async function onSignIn() {
  isLoading.value = true;
  try {
    await login();
  } catch (error: any) {
    toast.error(error.message || "Login failed");
    isLoading.value = false;
  }
}
</script>

<template>
  <div class="grid gap-6">
    <div class="grid gap-4">
      <Button
        class="w-full gap-3 h-12 bg-primary hover:bg-primary/90 text-primary-foreground shadow-lg shadow-primary/20 transition-all duration-300 hover:scale-[1.02] active:scale-[0.98] rounded-xl font-bold"
        type="button"
        :disabled="isLoading"
        @click="onSignIn"
      >
        <Loader2 v-if="isLoading" class="mr-2 h-5 w-5 animate-spin" />
        <Icon
          v-else
          name="mdi:shield-lock-outline"
          class="size-6 transition-transform group-hover:rotate-12"
        />
        Enter BangLipai Secure Portal
      </Button>

      <p class="text-xs text-center text-muted-foreground px-8">
        By clicking continue, you agree to our
        <NuxtLink
          to="#"
          class="underline underline-offset-4 hover:text-primary"
        >
          Terms of Service
        </NuxtLink>
        and
        <NuxtLink
          to="#"
          class="underline underline-offset-4 hover:text-primary"
        >
          Privacy Policy </NuxtLink
        >.
      </p>
    </div>
  </div>
</template>
