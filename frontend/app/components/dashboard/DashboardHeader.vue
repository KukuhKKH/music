<script setup lang="ts">
import { RefreshCcw } from "lucide-vue-next"

defineProps<{
  pending: boolean
}>()

defineEmits<{
  refresh: []
}>()
</script>

<template>
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Dashboard</h1>
      <p class="text-muted-foreground">
        Overview of your music library statistics.
      </p>
    </div>
    <Button
      variant="outline"
      size="icon"
      :disabled="pending"
      @click="$emit('refresh')"
    >
      <RefreshCcw :class="{ 'animate-spin': pending }" class="h-4 w-4" />
    </Button>
  </div>
</template>
