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
import type { MusicMeta, Track } from "~/types/music";
import { usePlayerStore } from "~/stores/player";
import { formatDuration, formatFileSize } from "~/lib/format";

const props = defineProps<{
  searchQuery: string;
  tracks: Track[];
  isLoading: boolean;
  meta: MusicMeta | null;
  selectedIds: number[];
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

const isAllSelected = computed(
  () =>
    props.tracks.length > 0 &&
    props.tracks.every((t) => props.selectedIds.includes(t.id))
);

function onCheckboxClick(_e: Event, id: number) {
  emit("toggleSelect", id);
}

function onSelectAllClick(_e: Event) {
  emit("toggleSelectAll");
}

// Inline editing functions
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

  const { $api } = useNuxtApp();

  try {
    await $api(`/music/${track.id}`, {
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
  <div class="overflow-x-auto relative">
    <Table>
      <TableHeader>
        <TableRow class="hover:bg-transparent border-b bg-muted/30">
          <TableHead class="w-[48px] pl-4 text-center">
            <!-- Native checkbox for maximum reliability -->
            <input
              type="checkbox"
              class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 cursor-pointer"
              :checked="isAllSelected"
              @change="onSelectAllClick"
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
        <!-- Loading State -->
        <TableRow v-if="isLoading && tracks.length === 0">
          <TableCell
            colspan="7"
            class="h-48 text-center text-muted-foreground uppercase text-[10px] tracking-widest font-bold"
          >
            <Loader2 class="h-6 w-6 animate-spin mx-auto mb-2 opacity-50" />
            Synchronizing Library
          </TableCell>
        </TableRow>

        <!-- Empty State -->
        <TableRow v-else-if="tracks.length === 0">
          <TableCell colspan="7" class="h-64 text-center">
            <div
              class="flex flex-col items-center justify-center space-y-2 opacity-40"
            >
              <FileMusic class="h-10 w-10" />
              <p class="text-xs font-bold uppercase tracking-tighter">
                No items in view
              </p>
            </div>
          </TableCell>
        </TableRow>

        <!-- Data Rows -->
        <TableRow
          v-for="track in tracks"
          :key="track.id"
          class="group transition-colors h-[64px] border-b last:border-0 relative"
          :class="{
            'bg-blue-600/5 hover:bg-blue-600/10': selectedIds.includes(
              track.id
            ),
            'bg-primary/[0.03] hover:bg-primary/[0.06]':
              player.currentTrack?.id === track.id &&
              !selectedIds.includes(track.id),
            'hover:bg-muted/30':
              !selectedIds.includes(track.id) &&
              player.currentTrack?.id !== track.id,
          }"
        >
          <TableCell class="pl-4 text-center">
            <input
              type="checkbox"
              class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 cursor-pointer"
              :checked="selectedIds.includes(track.id)"
              @change="onCheckboxClick($event, track.id)"
            />
          </TableCell>
          <TableCell>
            <Button
              variant="ghost"
              size="icon"
              class="h-9 w-9 rounded-full transition-all"
              :class="
                player.currentTrack?.id === track.id && player.isPlaying
                  ? 'bg-primary text-primary-foreground'
                  : 'hover:bg-primary/10'
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
              />
            </Button>
          </TableCell>
          <TableCell>
            <div v-if="editingId === track.id" class="flex items-center gap-1">
              <Input
                v-model="editingTitle"
                size="sm"
                class="h-8 min-w-[120px]"
                @keyup.enter="saveTitle(track)"
              />
              <Button
                size="icon"
                variant="ghost"
                class="h-8 w-8 text-green-600"
                @click="saveTitle(track)"
              >
                <Check class="h-4 w-4" />
              </Button>
              <Button
                size="icon"
                variant="ghost"
                class="h-8 w-8 text-destructive"
                @click="cancelEdit"
              >
                <X class="h-4 w-4" />
              </Button>
            </div>
            <div v-else class="flex flex-col min-w-0">
              <div class="flex items-center gap-2 group/title">
                <span
                  class="font-bold text-sm truncate"
                  :class="
                    player.currentTrack?.id === track.id ? 'text-primary' : ''
                  "
                >
                  {{ track.title }}
                </span>
                <Edit2
                  class="h-3 w-3 opacity-0 group-hover/title:opacity-40 cursor-pointer"
                  @click="startEdit(track)"
                />
              </div>
              <span
                class="text-[10px] text-muted-foreground sm:hidden uppercase tracking-tight"
              >
                {{ track.artist }}
              </span>
            </div>
          </TableCell>
          <TableCell
            class="hidden sm:table-cell text-xs font-medium text-muted-foreground uppercase tracking-tight truncate"
          >
            {{ track.artist }}
          </TableCell>
          <TableCell
            class="hidden md:table-cell text-xs font-mono text-muted-foreground/70"
          >
            {{ formatDuration(track.duration) }}
          </TableCell>
          <TableCell
            class="hidden lg:table-cell text-xs font-mono text-muted-foreground/70"
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
              <DropdownMenuContent align="end" class="w-[160px]">
                <DropdownMenuItem @click="startEdit(track)">
                  <Edit2 class="mr-2 h-3.5 w-3.5" /> Rename
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <AlertDialog>
                  <AlertDialogTrigger as-child>
                    <DropdownMenuItem class="text-destructive" @select.prevent>
                      <Trash2 class="mr-2 h-3.5 w-3.5" />
                      Delete
                    </DropdownMenuItem>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Confirm Delete</AlertDialogTitle>
                      <AlertDialogDescription>
                        Permanently remove "{{ track.title }}"?
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        class="bg-destructive"
                        @click="emit('deleteTrack', track.id)"
                      >
                        Delete Now
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
      v-if="meta && meta.total_page > 1"
      class="flex items-center justify-between p-4 border-t bg-muted/5"
    >
      <span class="text-[10px] font-bold uppercase text-muted-foreground">
        Found {{ meta.count }} tracks
      </span>
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
        <Button
          v-for="p in meta.total_page"
          :key="p"
          size="sm"
          variant="ghost"
          class="h-8 w-8 text-[11px] font-bold"
          :class="currentPage === p ? 'bg-primary text-primary-foreground' : ''"
          @click="emit('update:currentPage', p)"
        >
          {{ p }}
        </Button>
        <Button
          variant="ghost"
          size="sm"
          class="h-8 w-8 p-0"
          :disabled="currentPage === meta.total_page"
          @click="emit('update:currentPage', currentPage + 1)"
        >
          <ChevronRight class="h-4 w-4" />
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
:deep(td) {
  padding-top: 0.75rem;
  padding-bottom: 0.75rem;
}
</style>
