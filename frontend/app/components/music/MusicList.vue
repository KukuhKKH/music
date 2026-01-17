<script setup lang="ts">
import { useDebounceFn } from "@vueuse/core";
import { Music, Search, Trash2, X } from "lucide-vue-next";
import { toast } from "vue-sonner";
import type { MusicMeta, MusicResponse, Track } from "~/types/music";

const { $api } = useNuxtApp();

const tracks = ref<Track[]>([]);
const meta = ref<MusicMeta | null>(null);
const isLoading = ref(false);
const searchQuery = ref("");
const currentPage = ref(1);

// Selection state
const selectedIds = ref<number[]>([]);
const selectedCount = computed(() => selectedIds.value.length);

async function fetchTracks() {
  isLoading.value = true;
  try {
    const response = await $api<MusicResponse>("/music", {
      query: {
        page: currentPage.value,
        search: searchQuery.value,
        limit: 10,
      },
    });

    if (response) {
      tracks.value = response.data || [];
      meta.value = response.meta || null;
      selectedIds.value = [];
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
  const allOnPage = tracks.value.map((t) => t.id);
  const allCurrentSelected =
    allOnPage.length > 0 &&
    allOnPage.every((id) => selectedIds.value.includes(id));

  if (allCurrentSelected) {
    selectedIds.value = selectedIds.value.filter(
      (id) => !allOnPage.includes(id)
    );
  } else {
    const newSelection = [...selectedIds.value];
    allOnPage.forEach((id) => {
      if (!newSelection.includes(id)) newSelection.push(id);
    });
    selectedIds.value = newSelection;
  }
}

function toggleSelect(id: number) {
  const index = selectedIds.value.indexOf(id);
  if (index > -1) {
    selectedIds.value = selectedIds.value.filter((i) => i !== id);
  } else {
    selectedIds.value = [...selectedIds.value, id];
  }
}

async function bulkDelete() {
  const ids = [...selectedIds.value];
  if (ids.length === 0) return;

  isLoading.value = true;
  const count = ids.length;

  try {
    for (const id of ids) {
      await $api(`/music/${id}`, {
        method: "DELETE",
      });
    }
    toast.success(`${count} tracks deleted`);
    selectedIds.value = [];
    fetchTracks();
  } catch (err) {
    console.error("Bulk delete error:", err);
    toast.error("Partial deletion occurred");
    fetchTracks();
  } finally {
    isLoading.value = false;
  }
}

async function deleteTrack(id: number) {
  try {
    await $api(`/music/${id}`, {
      method: "DELETE",
    });
    tracks.value = tracks.value.filter((t) => t.id !== id);
    selectedIds.value = selectedIds.value.filter((i) => i !== id);
    toast.success("Track deleted");
    fetchTracks();
  } catch {
    toast.error("Failed to delete track");
  }
}
</script>

<template>
  <Card class="border-primary/10 shadow-sm overflow-hidden">
    <!-- Card Header -->
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
            <CardDescription>Library Management</CardDescription>
          </div>
        </div>

        <div class="relative w-full max-w-[280px]">
          <Search
            class="absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-muted-foreground"
          />
          <Input
            v-model="searchQuery"
            type="search"
            placeholder="Search tracks..."
            class="pl-9 bg-background/50 border-none shadow-none h-8 text-xs focus-visible:ring-1"
          />
        </div>
      </div>
    </CardHeader>

    <!-- NEW BULK ACTION BAR - Simplified logic for testing -->
    <div
      v-if="selectedCount > 0"
      class="bg-black text-white px-4 py-3 flex items-center justify-between border-y border-white/10"
      style="background: #1e40af !important; color: white !important"
    >
      <div class="flex items-center gap-3">
        <Trash2 class="h-4 w-4" />
        <span class="text-xs font-bold uppercase tracking-widest">
          {{ selectedCount }} Items Selected
        </span>
      </div>
      <div class="flex items-center gap-3">
        <AlertDialog>
          <AlertDialogTrigger as-child>
            <Button
              size="sm"
              variant="secondary"
              class="bg-white text-blue-800 hover:bg-white/90 h-7"
            >
              Delete All Selected
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure?</AlertDialogTitle>
              <AlertDialogDescription>
                Delete {{ selectedCount }} selected tracks permanently?
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
        <button
          class="text-white/60 hover:text-white"
          @click="selectedIds = []"
        >
          <X class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Table -->
    <CardContent class="p-0">
      <MusicTrackTable
        :search-query="searchQuery"
        :current-page="currentPage"
        :tracks="tracks"
        :is-loading="isLoading"
        :meta="meta"
        :selected-ids="selectedIds"
        @update:search-query="(val: string) => (searchQuery = val)"
        @update:current-page="(val: number) => (currentPage = val)"
        @toggle-select-all="toggleSelectAll"
        @toggle-select="toggleSelect"
        @delete-track="deleteTrack"
        @bulk-delete="bulkDelete"
        @fetch="fetchTracks"
      />
    </CardContent>
  </Card>
</template>
