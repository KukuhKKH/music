<script setup lang="ts">
import { Upload } from "lucide-vue-next";

definePageMeta({
  title: "My Music Library",
});
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
      <Button size="lg" class="shadow-md">
        <Upload class="mr-2 h-5 w-5" />
        Upload New Track
      </Button>
    </div>

    <MusicList />
  </div>
</template>
