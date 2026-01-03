<script setup lang="ts">
import { Upload } from "lucide-vue-next";

definePageMeta({
  title: "My Music Library",
});

const isUploadOpen = ref(false);
const musicListRef = ref<{ fetchTracks: () => void } | null>(null);

function handleUploadSuccess() {
  musicListRef.value?.fetchTracks();
}
</script>

<template>
  <div class="p-6 space-y-8">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div class="space-y-1">
        <h1 class="text-3xl font-extrabold tracking-tight">My Library</h1>
        <p class="text-sm text-muted-foreground">
          Manage, preview and update your music collection metadata.
        </p>
      </div>
      <Button
        size="lg"
        class="shadow-md bg-primary hover:bg-primary/90 text-primary-foreground font-bold px-8 transition-all hover:scale-105 active:scale-95"
        @click="isUploadOpen = true"
      >
        <Upload class="mr-2 h-5 w-5" />
        Upload New Track
      </Button>
    </div>

    <!-- Main Content Area -->
    <MusicList ref="musicListRef" />

    <!-- Upload Dialog Component -->
    <MusicUploadDialog
      v-model:open="isUploadOpen"
      @success="handleUploadSuccess"
    />
  </div>
</template>
