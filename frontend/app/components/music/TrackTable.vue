<script setup lang="ts">
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
  Trash2,
  X,
} from "lucide-vue-next";
import { toast } from "vue-sonner";
import { usePlayerStore } from "~/stores/player";
import type { MusicMeta, Track } from "~/types/music";
import { formatDuration, formatFileSize } from "~/lib/format";

const props = defineProps<{
  searchQuery: string;
  tracks: Track[];
  isLoading: boolean;
  meta: MusicMeta | null;
  selectedIds: Set<number>;
  currentPage: number;
}>();

const emit = defineEmits<{
  (e: "update:searchQuery", value: string): void;
  (e: "update:currentPage", value: number): void;
  (e: "toggleSelectAll"): void;
  (e: "toggleSelect", id: number): void;
  (e: "deleteTrack", id: number): void;
  (e: "bulkDelete"): void;
  (e: "fetch"): void;
}>();

const player = usePlayerStore();
const editingId = ref<number | null>(null);
const editingTitle = ref("");

const config = useRuntimeConfig();

const isAllSelected = computed(
  () =>
    props.tracks.length > 0 && props.selectedIds.size === props.tracks.length
);

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
</script>

<template>
  <div class="overflow-x-auto">
    <Table>
      <TableHeader>
        <TableRow class="hover:bg-transparent border-b bg-muted/30">
          <TableHead class="w-[48px] pl-4">
            <Checkbox
              :checked="isAllSelected"
              @update:checked="emit('toggleSelectAll')"
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
          <TableCell colspan="8" class="h-48 text-center">
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

        <TableRow v-else-if="tracks.length === 0" class="hover:bg-transparent">
          <TableCell colspan="8" class="h-64 text-center">
            <div class="flex flex-col items-center justify-center space-y-3">
              <div class="p-6 rounded-full bg-muted/30">
                <FileMusic class="h-10 w-10 text-muted-foreground/30" />
              </div>
              <div class="space-y-1">
                <p class="text-sm font-bold">No results found</p>
                <p class="text-xs text-muted-foreground max-w-[200px] mx-auto">
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
              ? 'bg-primary/4 hover:bg-primary/8'
              : 'hover:bg-muted/30'
          "
        >
          <TableCell class="pl-4">
            <Checkbox
              :checked="selectedIds.has(track.id)"
              @update:checked="emit('toggleSelect', track.id)"
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
                  class="font-bold text-[13px] sm:text-sm tracking-tight truncate font-sans"
                  :class="
                    player.currentTrack?.id === track.id
                      ? 'text-primary'
                      : 'text-foreground'
                  "
                >
                  {{ track.title }}
                </span>
                <span
                  class="text-[10px] text-muted-foreground sm:hidden truncate leading-none mt-1 uppercase tracking-tight font-medium"
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
              class="text-xs font-semibold text-muted-foreground/60 truncate uppercase tracking-tight"
            >
              {{ track.artist || "Unknown Artist" }}
            </span>
          </TableCell>
          <TableCell
            class="hidden md:table-cell text-xs text-muted-foreground/80 tabular-nums font-mono"
          >
            {{ formatDuration(track.duration) }}
          </TableCell>
          <TableCell
            class="hidden lg:table-cell text-xs text-muted-foreground/80 tabular-nums font-mono"
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
                        <span class="font-bold text-foreground">
                          "{{ track.title }}"
                        </span>
                        ?
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                        @click="emit('deleteTrack', track.id)"
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

    <!-- Pagination -->
    <div
      v-if="meta && meta.last_page > 1"
      class="flex flex-col gap-4 items-center justify-between p-4 md:flex-row border-t bg-muted/5"
    >
      <div
        class="text-[10px] text-muted-foreground uppercase tracking-widest font-bold"
      >
        Found <span class="text-foreground">{{ meta.total }}</span> items
      </div>
      <div class="flex items-center gap-1">
        <Button
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          :disabled="currentPage === 1"
          @click="emit('update:currentPage', currentPage - 1)"
        >
          <ChevronLeft class="h-4 w-4" />
        </Button>
        <div class="flex items-center gap-0.5">
          <Button
            v-for="p in meta.last_page"
            :key="p"
            size="sm"
            variant="ghost"
            class="h-8 w-8 p-0 text-[11px] font-bold transition-all"
            :class="
              currentPage === p
                ? 'bg-primary text-primary-foreground shadow-sm hover:bg-primary/90'
                : 'hover:bg-primary/10'
            "
            @click="emit('update:currentPage', p)"
          >
            {{ p }}
          </Button>
        </div>
        <Button
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          :disabled="currentPage === meta.last_page"
          @click="emit('update:currentPage', currentPage + 1)"
        >
          <ChevronRight class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
:deep(td),
:deep(th) {
  padding-top: 0.85rem;
  padding-bottom: 0.85rem;
}

.tabular-nums {
  font-variant-numeric: tabular-nums;
}
</style>
