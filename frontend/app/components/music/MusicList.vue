<script setup lang="ts">
import { useDebounceFn } from "@vueuse/core";
import { Music, Search, Trash2 } from "lucide-vue-next";
import { toast } from "vue-sonner";
import type { MusicMeta, MusicResponse, Track } from "~/types/music";

const config = useRuntimeConfig();

const tracks = ref<Track[]>([]);
const meta = ref<MusicMeta | null>(null);
const isLoading = ref(false);
const searchQuery = ref("");
const currentPage = ref(1);
const selectedIds = ref<Set<number>>(new Set());

async function fetchTracks() {
  isLoading.value = true;
  try {
    const response = await $fetch<MusicResponse>(
      `${config.public.apiBase}/music`,
      {
        query: {
          page: currentPage.value,
          search: searchQuery.value,
          limit: 10,
        },
      }
    );
    
    if (response) {
      tracks.value = response.data || [];
      meta.value = response.meta || null;
      selectedIds.value.clear();
    }
  } catch (err) {
    console.error("Fetch tracks error:", err);
    toast.error("Failed to fetch tracks");
  } finally {
    isLoading.value = false;
  }
}

defineExpose({ fetchTracks });

const debouncedSearch = useDebounceFn(() => {
  currentPage.value = 1;
  fetchTracks();
}, 500);

watch(searchQuery, debouncedSearch);
watch(currentPage, fetchTracks);

onMounted(fetchTracks);

function toggleSelectAll() {
  const allCurrentOnPageSelected = tracks.value.every((t) =>
    selectedIds.value.has(t.id)
  );
  const newSelected = new Set(selectedIds.value);

  if (allCurrentOnPageSelected) {
    tracks.value.forEach((t) => newSelected.delete(t.id));
  } else {
    tracks.value.forEach((t) => newSelected.add(t.id));
  }
  selectedIds.value = newSelected;
}

function toggleSelect(id: number) {
  const newSelected = new Set(selectedIds.value);
  if (newSelected.has(id)) {
    newSelected.delete(id);
  } else {
    newSelected.add(id);
  }
  selectedIds.value = newSelected;
}

async function deleteTrack(id: number) {
  try {
    await $fetch(`${config.public.apiBase}/music/${id}`, { method: "DELETE" });
    tracks.value = tracks.value.filter((t) => t.id !== id);
    toast.success("Track deleted");
    if (tracks.value.length === 0 && currentPage.value > 1) {
      currentPage.value--;
    }
    fetchTracks();
  } catch (err) {
    console.error("Delete error:", err);
    toast.error("Failed to delete track");
  }
}

async function bulkDelete() {
  const ids = Array.from(selectedIds.value);
  if (ids.length === 0) return;
  isLoading.value = true;
  try {
    for (const id of ids) {
      await $fetch(`${config.public.apiBase}/music/${id}`, {
        method: "DELETE",
      });
    }
    toast.success(`${ids.length} tracks deleted`);
    selectedIds.value.clear();
    fetchTracks();
  } catch (err) {
    console.error("Bulk delete error:", err);
    toast.error("Failed to perform bulk delete");
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <Card class="border-primary/10 shadow-sm overflow-hidden">
    <!-- Card Header with Search & Actions -->
    <CardHeader class="border-b bg-muted/30 pb-4">
      <div
        class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
      >
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

        <!-- Search Bar in Header -->
        <div class="relative w-full max-w-[280px]">
          <Search
            class="absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground"
          />
          <Input
            v-model="searchQuery"
            type="search"
            placeholder="Quick search..."
            class="pl-9 bg-background/50 border-none shadow-none h-8 text-xs focus-visible:ring-1"
          />
        </div>
      </div>
    </CardHeader>

    <!-- Bulk Actions Bar -->
    <div
      v-if="selectedIds.size > 0"
      class="bg-primary/5 border-b px-4 py-2 flex items-center justify-between animate-in fade-in slide-in-from-top-2"
    >
      <span class="text-xs font-bold text-primary italic">
        {{ selectedIds.size }} items selected for action
      </span>
      <div class="flex items-center gap-2">
        <AlertDialog>
          <AlertDialogTrigger as-child>
            <Button
              variant="destructive"
              size="sm"
              class="h-7 text-[10px] font-bold uppercase tracking-wider px-3"
            >
              <Trash2 class="mr-1.5 h-3 w-3" />
              Exec Delete
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Action Confirmation</AlertDialogTitle>
              <AlertDialogDescription>
                Are you sure you want to delete these
                {{ selectedIds.size }} tracks?
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                class="bg-destructive text-destructive-foreground"
                @click="bulkDelete"
              >
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>
    </div>

    <!-- Table Body -->
    <CardContent class="p-0">
      <MusicTrackTable
        :search-query="searchQuery"
        :current-page="currentPage"
        :tracks="tracks"
        :is-loading="isLoading"
        :meta="meta"
        :selected-ids="selectedIds"
        @update:search-query="(val) => (searchQuery = val)"
        @update:current-page="(val) => (currentPage = val)"
        @toggle-select-all="toggleSelectAll"
        @toggle-select="toggleSelect"
        @delete-track="deleteTrack"
        @bulk-delete="bulkDelete"
        @fetch="fetchTracks"
      />
    </CardContent>
  </Card>
</template>
