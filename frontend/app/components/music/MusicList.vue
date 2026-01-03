<script setup lang="ts">
import { useDebounceFn } from "@vueuse/core";
import {
  Check,
  ChevronLeft,
  ChevronRight,
  Edit2,
  FileMusic,
  Loader2,
  MoreHorizontal,
  Pause,
  Play,
  Search,
  Trash2,
  X,
} from "lucide-vue-next";
import { toast } from "vue-sonner";
import { usePlayerStore } from "~/stores/player";
import type { MusicMeta, MusicResponse, Track } from "~/types/music";
import { formatDuration, formatFileSize } from "~/lib/format";

const config = useRuntimeConfig();
const player = usePlayerStore();

const tracks = ref<Track[]>([]);
const meta = ref<MusicMeta | null>(null);
const isLoading = ref(false);
const searchQuery = ref("");
const currentPage = ref(1);

const selectedIds = ref<Set<number>>(new Set());
const isAllSelected = computed(
  () =>
    tracks.value.length > 0 && selectedIds.value.size === tracks.value.length
);

const editingId = ref<number | null>(null);
const editingTitle = ref("");

async function fetchTracks() {
  isLoading.value = true;
  try {
    const { data } = await useFetch<MusicResponse>(
      `${config.public.apiBase}/music`,
      {
        query: {
          page: currentPage.value,
          search: searchQuery.value,
          limit: 10,
        },
      }
    );

    if (data.value) {
      tracks.value = data.value.data || [];
      meta.value = data.value.meta || null;
      selectedIds.value.clear();
    }
  } catch {
    toast.error("Failed to fetch tracks");
  } finally {
    isLoading.value = false;
  }
}

const debouncedSearch = useDebounceFn(() => {
  currentPage.value = 1;
  fetchTracks();
}, 500);

watch(searchQuery, debouncedSearch);

onMounted(() => {
  fetchTracks();
});

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value.clear();
  } else {
    selectedIds.value = new Set(tracks.value.map((t) => t.id));
  }
}

function toggleSelect(id: number) {
  if (selectedIds.value.has(id)) {
    selectedIds.value.delete(id);
  } else {
    selectedIds.value.add(id);
  }
}

function startEdit(track: Track) {
  editingId.value = track.id;
  editingTitle.value = track.title;
}

function cancelEdit() {
  editingId.value = null;
  editingTitle.value = "";
}

async function saveTitle(track: Track) {
  if (!editingTitle.value || editingTitle.value === track.title) {
    cancelEdit();
    return;
  }

  try {
    await $fetch(`${config.public.apiBase}/music/${track.id}`, {
      method: "PUT",
      body: {
        title: editingTitle.value,
        artist: track.artist,
        album: track.album || "",
      },
    });

    track.title = editingTitle.value;
    toast.success("Title updated");
  } catch {
    toast.error("Failed to update title");
  } finally {
    cancelEdit();
  }
}

async function deleteTrack(id: number) {
  try {
    await $fetch(`${config.public.apiBase}/music/${id}`, {
      method: "DELETE",
    });

    tracks.value = tracks.value.filter((t) => t.id !== id);
    toast.success("Track deleted");

    if (tracks.value.length === 0 && currentPage.value > 1) {
      currentPage.value--;
      fetchTracks();
    }
  } catch {
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
  } catch {
    toast.error("Failed to perform bulk delete");
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col">
    <!-- Top Bar (Search & Actions) -->
    <div
      class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between p-4 bg-muted/20 border-b"
    >
      <div class="relative w-full max-w-sm">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground"
        />
        <Input
          v-model="searchQuery"
          type="search"
          placeholder="Search tracks..."
          class="pl-9 bg-background border-none shadow-sm h-9"
        />
      </div>

      <div class="flex items-center gap-2">
        <AlertDialog v-if="selectedIds.size > 0">
          <AlertDialogTrigger as-child>
            <Button variant="destructive" size="sm" class="h-9 shadow-sm">
              <Trash2 class="mr-2 h-4 w-4" />
              Delete ({{ selectedIds.size }})
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This will permanently delete {{ selectedIds.size }} selected
                tracks.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction
                class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                @click="bulkDelete"
              >
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>

        <Button size="sm" variant="secondary" class="h-9 shadow-sm">
          <Play class="mr-2 h-4 w-4" />
          Play All
        </Button>
      </div>
    </div>

    <!-- Table Section (Integrated without extra borders) -->
    <div class="overflow-x-auto">
      <Table>
        <TableHeader>
          <TableRow class="hover:bg-transparent border-b">
            <TableHead class="w-[48px] pl-4">
              <Checkbox
                :checked="isAllSelected"
                @update:checked="toggleSelectAll"
              />
            </TableHead>
            <TableHead class="w-[48px]" />
            <TableHead class="font-bold text-foreground"> Title </TableHead>
            <TableHead class="font-bold text-foreground hidden sm:table-cell">
              Artist
            </TableHead>
            <TableHead
              class="font-bold text-foreground hidden md:table-cell w-[100px]"
            >
              Duration
            </TableHead>
            <TableHead
              class="font-bold text-foreground hidden lg:table-cell w-[100px]"
            >
              Size
            </TableHead>
            <TableHead class="w-[80px] text-right pr-4"> Actions </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-if="isLoading && tracks.length === 0"
            class="hover:bg-transparent"
          >
            <TableCell colspan="7" class="h-48 text-center">
              <div class="flex flex-col items-center gap-3">
                <Loader2 class="h-8 w-8 animate-spin text-primary opacity-50" />
                <p
                  class="text-xs text-muted-foreground font-medium uppercase tracking-widest"
                >
                  Fetching Library
                </p>
              </div>
            </TableCell>
          </TableRow>

          <TableRow
            v-else-if="tracks.length === 0"
            class="hover:bg-transparent"
          >
            <TableCell colspan="7" class="h-64 text-center">
              <div class="flex flex-col items-center justify-center space-y-3">
                <div class="p-6 rounded-full bg-muted/30">
                  <FileMusic class="h-10 w-10 text-muted-foreground/30" />
                </div>
                <div class="space-y-1">
                  <p class="text-sm font-bold">No results found</p>
                  <p
                    class="text-xs text-muted-foreground max-w-[200px] mx-auto"
                  >
                    Try refining your search or adding some new music.
                  </p>
                </div>
              </div>
            </TableCell>
          </TableRow>

          <TableRow
            v-for="track in tracks"
            :key="track.id"
            class="group transition-colors h-[64px] border-b last:border-0"
            :class="
              player.currentTrack?.id === track.id
                ? 'bg-primary/[0.04] hover:bg-primary/[0.08]'
                : 'hover:bg-muted/30'
            "
          >
            <TableCell class="pl-4">
              <Checkbox
                :checked="selectedIds.has(track.id)"
                @update:checked="toggleSelect(track.id)"
              />
            </TableCell>
            <TableCell>
              <Button
                variant="ghost"
                size="icon"
                class="h-9 w-9 rounded-full transition-all"
                :class="
                  player.currentTrack?.id === track.id && player.isPlaying
                    ? 'bg-primary text-primary-foreground hover:bg-primary/90 shadow-md scale-105'
                    : 'hover:bg-primary/10 hover:text-primary active:scale-90 hover:scale-105'
                "
                @click="player.toggle(track)"
              >
                <component
                  :is="
                    player.currentTrack?.id === track.id && player.isPlaying
                      ? Pause
                      : Play
                  "
                  class="h-4 w-4"
                  :class="{
                    'fill-current':
                      player.currentTrack?.id === track.id && player.isPlaying,
                  }"
                />
              </Button>
            </TableCell>
            <TableCell>
              <div
                v-if="editingId === track.id"
                class="flex items-center gap-2 animate-in fade-in slide-in-from-left-1"
              >
                <Input
                  v-model="editingTitle"
                  size="sm"
                  class="h-8 min-w-[150px] bg-background border-primary/20"
                  @keyup.enter="saveTitle(track)"
                  @keyup.esc="cancelEdit"
                />
                <Button
                  size="icon"
                  variant="ghost"
                  class="h-8 w-8 text-green-600 hover:bg-green-50"
                  @click="saveTitle(track)"
                >
                  <Check class="h-4 w-4" />
                </Button>
                <Button
                  size="icon"
                  variant="ghost"
                  class="h-8 w-8 text-destructive hover:bg-destructive/5"
                  @click="cancelEdit"
                >
                  <X class="h-4 w-4" />
                </Button>
              </div>
              <div v-else class="flex items-center gap-3 group/title">
                <div class="flex flex-col min-w-0">
                  <span
                    class="font-bold text-[13px] sm:text-sm tracking-tight truncate"
                    :class="
                      player.currentTrack?.id === track.id
                        ? 'text-primary'
                        : 'text-foreground'
                    "
                  >
                    {{ track.title }}
                  </span>
                  <span
                    class="text-[10px] text-muted-foreground sm:hidden truncate leading-none mt-1"
                  >
                    {{ track.artist || "Unknown Artist" }}
                  </span>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  class="h-6 w-6 opacity-0 group-hover/title:opacity-100 transition-opacity"
                  @click="startEdit(track)"
                >
                  <Edit2 class="h-3 w-3" />
                </Button>
              </div>
            </TableCell>
            <TableCell class="hidden sm:table-cell">
              <span
                class="text-sm font-medium text-muted-foreground/60 truncate"
              >
                {{ track.artist || "Unknown Artist" }}
              </span>
            </TableCell>
            <TableCell
              class="hidden md:table-cell text-xs text-muted-foreground/80 tabular-nums"
            >
              {{ formatDuration(track.duration) }}
            </TableCell>
            <TableCell
              class="hidden lg:table-cell text-xs text-muted-foreground/80 tabular-nums"
            >
              {{ formatFileSize(track.file_size) }}
            </TableCell>
            <TableCell class="text-right pr-4">
              <DropdownMenu>
                <DropdownMenuTrigger as-child>
                  <Button
                    variant="ghost"
                    size="icon"
                    class="h-8 w-8 opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <MoreHorizontal class="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" class="w-[180px]">
                  <DropdownMenuItem @click="player.toggle(track)">
                    <component
                      :is="
                        player.currentTrack?.id === track.id && player.isPlaying
                          ? Pause
                          : Play
                      "
                      class="mr-2 h-4 w-4"
                    />
                    {{
                      player.currentTrack?.id === track.id && player.isPlaying
                        ? "Pause"
                        : "Play Preview"
                    }}
                  </DropdownMenuItem>
                  <DropdownMenuItem @click="startEdit(track)">
                    <Edit2 class="mr-2 h-4 w-4" />
                    Edit Title
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <AlertDialog>
                    <AlertDialogTrigger as-child>
                      <DropdownMenuItem
                        class="text-destructive focus:bg-destructive/10 focus:text-destructive"
                        @select.prevent
                      >
                        <Trash2 class="mr-2 h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Delete Track</AlertDialogTitle>
                        <AlertDialogDescription>
                          Are you sure you want to delete
                          <span class="font-bold text-foreground"
                            >"{{ track.title }}"</span
                          >?
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                          class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                          @click="deleteTrack(track.id)"
                        >
                          Delete
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </DropdownMenuContent>
              </DropdownMenu>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <!-- Pagination (Border Top for spacing) -->
    <div
      v-if="meta && meta.last_page > 1"
      class="flex flex-col gap-4 items-center justify-between p-4 md:flex-row border-t bg-muted/5 font-medium"
    >
      <div class="text-[11px] text-muted-foreground uppercase tracking-widest">
        Library Size:
        <span class="text-foreground font-bold">{{ meta.total }}</span> content
        items
      </div>
      <div class="flex items-center gap-1">
        <Button
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          :disabled="currentPage === 1"
          @click="
            () => {
              currentPage--;
              fetchTracks();
            }
          "
        >
          <ChevronLeft class="h-4 w-4" />
        </Button>
        <div class="flex items-center gap-0.5">
          <Button
            v-for="p in meta.last_page"
            :key="p"
            size="sm"
            variant="ghost"
            class="h-8 w-8 p-0 text-xs transition-all"
            :class="
              currentPage === p
                ? 'bg-primary text-primary-foreground shadow-sm hover:bg-primary/90'
                : 'hover:bg-primary/10'
            "
            @click="
              () => {
                currentPage = p;
                fetchTracks();
              }
            "
          >
            {{ p }}
          </Button>
        </div>
        <Button
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          :disabled="currentPage === meta.last_page"
          @click="
            () => {
              currentPage++;
              fetchTracks();
            }
          "
        >
          <ChevronRight class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Standardize cell vertical alignments */
:deep(td),
:deep(th) {
  padding-top: 0.85rem;
  padding-bottom: 0.85rem;
}

/* Ensure font numbers stay consistent */
.tabular-nums {
  font-variant-numeric: tabular-nums;
}
</style>
