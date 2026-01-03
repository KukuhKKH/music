<script setup lang="ts">
import { Music, Upload } from "lucide-vue-next";

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

    <!-- Main Content Area -->
    <Card class="border-primary/10 shadow-sm overflow-hidden">
      <CardHeader class="border-b bg-muted/30 pb-4">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-lg bg-primary/10 text-primary">
            <Music class="h-5 w-5" />
          </div>
          <div>
            <CardTitle class="text-xl"> Tracks </CardTitle>
            <CardDescription>
              A complete list of your uploaded audio files.
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent class="p-0">
        <MusicList />
      </CardContent>
    </Card>
  </div>
</template>
