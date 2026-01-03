<script setup lang="ts">
import { computed, ref } from "vue";
import {
  AlertCircle,
  CheckCircle2,
  FileAudio,
  Loader2,
  Music,
  Trash2,
  Upload,
} from "lucide-vue-next";
import { toast } from "vue-sonner";

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "success"): void;
}>();

const { extractMetadata } = useTrackMetadata();
const { uploadTrack } = useTrackUpload();
const config = useRuntimeConfig();

interface UploadFile {
  id: string;
  file: File;
  title: string;
  artist: string;
  album: string;
  duration: number;
  status: "idle" | "parsing" | "uploading" | "success" | "error";
  progress: number;
  errorMessage?: string;
}

const files = ref<UploadFile[]>([]);
const fileInput = ref<HTMLInputElement | null>(null);
const isDragging = ref(false);

const canUpload = computed(
  () =>
    files.value.length > 0 &&
    files.value.every((f) => f.status !== "uploading" && f.title.trim() !== "")
);

const totalProgress = computed(() => {
  if (files.value.length === 0) return 0;
  const sum = files.value.reduce((acc, f) => acc + f.progress, 0);
  return Math.round(sum / files.value.length);
});

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.files) {
    addFiles(Array.from(target.files));
  }
}

async function addFiles(newFiles: File[]) {
  const audioFiles = newFiles.filter((f) => f.type.startsWith("audio/"));

  if (audioFiles.length < newFiles.length) {
    toast.error("Some files were skipped as they are not valid audio.");
  }

  for (const file of audioFiles) {
    const id = Math.random().toString(36).substring(7);

    // Add initial entry
    files.value.push({
      id,
      file,
      title: file.name.replace(/\.[^/.]+$/, ""),
      artist: "",
      album: "",
      duration: 0,
      status: "parsing",
      progress: 0,
    });

    const targetFile = files.value.find((f) => f.id === id);
    if (!targetFile) continue;

    // Extract metadata using standardized composable
    const meta = await extractMetadata(file);
    targetFile.title = meta.title;
    targetFile.artist = meta.artist;
    targetFile.album = meta.album;
    targetFile.duration = meta.duration;
    targetFile.status = "idle";
  }

  // Clear input so same file can be selected again if removed
  if (fileInput.value) fileInput.value.value = "";
}

function removeFile(id: string) {
  files.value = files.value.filter((f) => f.id !== id);
}

function handleDrop(e: DragEvent) {
  isDragging.value = false;
  if (e.dataTransfer?.files) {
    addFiles(Array.from(e.dataTransfer.files));
  }
}

async function startUpload() {
  const toUpload = files.value.filter(
    (f) => f.status === "idle" || f.status === "error"
  );

  for (const item of toUpload) {
    item.status = "uploading";
    item.progress = 0;

    try {
      await uploadTrack({
        url: `${config.public.apiBase}/music`,
        file: item.file,
        metadata: {
          title: item.title,
          artist: item.artist,
          album: item.album,
          duration: item.duration,
        },
        onProgress: (p) => {
          item.progress = p;
        },
      });

      item.status = "success";
      item.progress = 100;
    } catch (err: any) {
      console.error("Upload failed for track:", item.title, err);
      item.status = "error";
      item.errorMessage = err.message || "Upload failed";
    }
  }

  // Check if everything finished successfully
  const allSuccess = files.value.every((f) => f.status === "success");
  if (allSuccess) {
    toast.success("All tracks uploaded successfully!");
    setTimeout(() => {
      resetAndClose();
    }, 1500);
  }
}

function resetAndClose() {
  emit("update:open", false);
  emit("success");
  setTimeout(() => {
    files.value = [];
  }, 300);
}
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent
      class="sm:max-w-[700px] gap-0 p-0 overflow-hidden bg-card border-primary/20 shadow-2xl"
    >
      <!-- Hidden file input for programatic trigger -->
      <input
        ref="fileInput"
        type="file"
        multiple
        accept="audio/*"
        class="hidden"
        @change="handleFileSelect"
      />

      <DialogHeader class="p-6 pb-0">
        <div class="flex items-center gap-3">
          <div class="p-2.5 rounded-xl bg-primary/10 text-primary shadow-inner">
            <Upload class="h-6 w-6" />
          </div>
          <div>
            <DialogTitle class="text-2xl font-black tracking-tight">
              Upload New Tracks
            </DialogTitle>
            <DialogDescription class="text-sm font-medium opacity-70">
              Drag and drop your audio files here.
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <div class="p-6 space-y-6">
        <!-- Empty State / Dropzone -->
        <div
          v-if="files.length === 0"
          class="relative group"
          @dragover.prevent="isDragging = true"
          @dragleave.prevent="isDragging = false"
          @drop.prevent="handleDrop"
        >
          <div
            class="border-2 border-dashed rounded-2xl p-12 transition-all duration-300 flex flex-col items-center justify-center text-center gap-4 cursor-pointer"
            :class="
              isDragging
                ? 'border-primary bg-primary/5 scale-[0.99]'
                : 'border-muted-foreground/20 hover:border-primary/40 hover:bg-primary/[0.02]'
            "
            @click="fileInput?.click()"
          >
            <div
              class="h-20 w-20 rounded-full bg-muted/50 flex items-center justify-center mb-2 group-hover:scale-110 transition-transform duration-500 shadow-sm"
            >
              <Music class="h-10 w-10 text-muted-foreground/40" />
            </div>
            <div class="space-y-1">
              <p class="text-lg font-bold tracking-tight">
                Click or drag audio files
              </p>
              <p
                class="text-xs text-muted-foreground font-medium uppercase tracking-widest px-8"
              >
                MP3, WAV, M4A, FLAC are supported
              </p>
            </div>
            <Button
              variant="secondary"
              size="sm"
              class="mt-2 font-bold px-6 shadow-sm"
            >
              Select Files
            </Button>
          </div>
        </div>

        <!-- File List with scrollable area -->
        <div
          v-else
          class="space-y-4 max-h-[400px] overflow-y-auto pr-2 custom-scrollbar"
        >
          <div
            v-for="item in files"
            :key="item.id"
            class="group relative flex flex-col gap-3 p-4 rounded-xl border bg-muted/30 transition-all duration-300 hover:border-primary/30 hover:bg-muted/50"
          >
            <div class="flex items-start gap-4">
              <div
                class="h-10 w-10 rounded-lg bg-background border flex items-center justify-center shrink-0 shadow-sm"
              >
                <div v-if="item.status === 'parsing'" class="animate-spin">
                  <Loader2 class="h-5 w-5 text-primary opacity-50" />
                </div>
                <FileAudio v-else class="h-5 w-5 text-primary" />
              </div>

              <div class="flex-1 grid grid-cols-1 sm:grid-cols-2 gap-3 min-w-0">
                <div class="space-y-1">
                  <Label
                    class="text-[10px] font-black uppercase text-muted-foreground tracking-tighter"
                  >
                    Title
                  </Label>
                  <Input
                    v-model="item.title"
                    placeholder="Track Title"
                    size="sm"
                    class="h-8 text-sm font-bold bg-background/50"
                    :disabled="
                      item.status === 'uploading' || item.status === 'success'
                    "
                  />
                </div>
                <div class="space-y-1">
                  <Label
                    class="text-[10px] font-black uppercase text-muted-foreground tracking-tighter"
                  >
                    Artist
                  </Label>
                  <Input
                    v-model="item.artist"
                    placeholder="Artist Name"
                    size="sm"
                    class="h-8 text-sm font-medium bg-background/50"
                    :disabled="
                      item.status === 'uploading' || item.status === 'success'
                    "
                  />
                </div>
                <div class="space-y-1 sm:col-span-2">
                  <Label
                    class="text-[10px] font-black uppercase text-muted-foreground tracking-tighter"
                  >
                    Album (Optional)
                  </Label>
                  <Input
                    v-model="item.album"
                    placeholder="Album Name"
                    size="sm"
                    class="h-8 text-sm bg-background/50"
                    :disabled="
                      item.status === 'uploading' || item.status === 'success'
                    "
                  />
                </div>
              </div>

              <!-- Actions column -->
              <div class="flex flex-col gap-2">
                <Button
                  v-if="
                    item.status !== 'uploading' && item.status !== 'success'
                  "
                  variant="ghost"
                  size="icon"
                  class="h-8 w-8 text-muted-foreground hover:text-destructive hover:bg-destructive/10"
                  @click="removeFile(item.id)"
                >
                  <Trash2 class="h-4 w-4" />
                </Button>
                <div
                  v-if="item.status === 'success'"
                  class="h-8 w-8 flex items-center justify-center text-green-500 animate-in zoom-in duration-300"
                >
                  <CheckCircle2 class="h-6 w-6" />
                </div>
                <div
                  v-if="item.status === 'error'"
                  class="h-8 w-8 flex items-center justify-center text-destructive"
                >
                  <AlertCircle class="h-6 w-6" />
                </div>
              </div>
            </div>

            <!-- Individual Progress bar -->
            <div v-if="item.status === 'uploading'" class="space-y-1.5">
              <div
                class="flex justify-between text-[10px] font-black uppercase tracking-widest text-primary italic"
              >
                <span>Uploading...</span>
                <span>{{ item.progress }}%</span>
              </div>
              <Progress
                :model-value="item.progress"
                class="h-1 bg-primary/20"
              />
            </div>

            <!-- Error message display -->
            <div
              v-if="item.status === 'error'"
              class="text-[10px] font-bold text-destructive uppercase tracking-widest"
            >
              {{ item.errorMessage || "Error occurred" }}
            </div>
          </div>

          <!-- Add More Button when files exist -->
          <div
            v-if="!files.some((f) => f.status === 'uploading')"
            class="flex justify-center pb-2"
          >
            <Button
              variant="outline"
              size="sm"
              class="rounded-full h-8 text-xs font-bold border-dashed border-primary/40 hover:border-primary hover:bg-primary/5 shadow-sm"
              @click="fileInput?.click()"
            >
              <Upload class="h-3.5 w-3.5 mr-2" />
              Add More Files
            </Button>
          </div>
        </div>
      </div>

      <DialogFooter
        class="p-6 bg-muted/20 border-t flex flex-col sm:flex-row gap-3"
      >
        <!-- Overall Progress display -->
        <div
          v-if="files.some((f) => f.status === 'uploading')"
          class="flex-1 flex flex-col gap-1.5 mr-4"
        >
          <div
            class="flex justify-between text-[10px] font-black uppercase tracking-widest text-muted-foreground font-mono"
          >
            <span>Overall Status</span>
            <span>{{ totalProgress }}%</span>
          </div>
          <Progress :model-value="totalProgress" class="h-1.5 shadow-sm" />
        </div>

        <div class="flex gap-2 w-full sm:w-auto ml-auto">
          <Button
            variant="ghost"
            size="sm"
            class="font-bold px-6"
            :disabled="files.some((f) => f.status === 'uploading')"
            @click="emit('update:open', false)"
          >
            Cancel
          </Button>
          <Button
            v-if="files.length > 0"
            size="sm"
            class="font-bold px-8 shadow-lg shadow-primary/20"
            :disabled="!canUpload"
            @click="startUpload"
          >
            <template v-if="files.some((f) => f.status === 'uploading')">
              <Loader2 class="h-4 w-4 mr-2 animate-spin" />
              Working...
            </template>
            <template v-else>
              Start Upload
              {{ files.length > 1 ? `(${files.length} Tracks)` : "" }}
            </template>
          </Button>
        </div>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: hsl(var(--primary) / 0.1);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: hsl(var(--primary) / 0.3);
}

:deep([data-slot="progress-indicator"]) {
  background-color: var(--primary) !important;
  transition: transform 0.3s ease-in-out;
}
</style>
