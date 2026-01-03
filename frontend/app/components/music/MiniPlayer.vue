<script setup lang="ts">
import { computed, ref } from "vue";
import {
  Music,
  Pause,
  Play,
  RotateCcw,
  RotateCw,
  Volume2,
  VolumeX,
  X,
} from "lucide-vue-next";
import { usePlayerStore } from "~/stores/player";
import { formatDuration } from "~/lib/format";

const player = usePlayerStore();

const isMuted = ref(false);
const lastVolume = ref(0.7);
const isHoveringVolume = ref(false);

function toggleMute() {
  if (isMuted.value) {
    player.setVolume(lastVolume.value);
    isMuted.value = false;
  } else {
    lastVolume.value = player.volume;
    player.setVolume(0);
    isMuted.value = true;
  }
}

const volumeValue = computed({
  get: () => [player.volume * 100],
  set: (val: number[]) => {
    const newVol = val[0] / 100;
    player.setVolume(newVol);
    if (newVol > 0) isMuted.value = false;
  },
});

const progressValue = computed({
  get: () => [player.currentTime],
  set: (val: number[]) => {
    player.seek(val[0]);
  },
});
</script>

<template>
  <Transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="translate-y-full opacity-0"
    enter-to-class="translate-y-0 opacity-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="translate-y-0 opacity-100"
    leave-to-class="translate-y-full opacity-0"
  >
    <div
      v-if="player.currentTrack"
      class="fixed bottom-6 left-1/2 -translate-x-1/2 w-[90%] max-w-2xl z-50 px-4 group"
    >
      <!-- Close Button (Hover only) -->
      <Button
        variant="destructive"
        size="icon"
        class="absolute -top-2 right-0 h-7 w-7 rounded-full shadow-lg z-30 border-2 border-background opacity-0 group-hover:opacity-100 transition-all duration-300 hover:scale-110 active:scale-90"
        @click="player.stop()"
      >
        <X class="h-3.5 w-3.5" />
      </Button>

      <div
        class="bg-card/95 backdrop-blur-3xl border border-primary/20 shadow-2xl rounded-2xl flex flex-col relative overflow-hidden transition-all hover:border-primary/40"
      >
        <!-- Main Content Area -->
        <div class="px-4 py-3 flex items-center gap-4">
          <!-- Left: Track Info -->
          <div class="flex items-center gap-3 min-w-0 flex-1">
            <div
              class="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center flex-shrink-0 relative shadow-sm"
            >
              <div
                v-if="player.isPlaying"
                class="absolute inset-0 bg-primary/20 animate-pulse rounded-xl"
              />
              <Music
                class="h-5 w-5 transition-colors"
                :class="
                  player.isPlaying ? 'text-primary' : 'text-muted-foreground'
                "
              />
            </div>
            <div class="min-w-0">
              <h4
                class="text-sm font-bold truncate leading-none tracking-tight"
              >
                {{ player.currentTrack.title }}
              </h4>
              <p class="text-[11px] text-muted-foreground truncate mt-1">
                {{ player.currentTrack.artist || "Unknown Artist" }}
              </p>
            </div>
          </div>

          <!-- Middle: Playback Controls -->
          <div class="flex items-center gap-1 sm:gap-3 flex-shrink-0">
            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 text-muted-foreground hover:text-primary transition-colors"
              @click="player.skip(-10)"
            >
              <RotateCcw class="h-4 w-4" />
            </Button>

            <Button
              variant="default"
              size="icon"
              class="h-10 w-10 rounded-full hover:scale-105 transition-all shadow-md active:scale-95 bg-primary text-primary-foreground"
              @click="player.toggle(player.currentTrack)"
            >
              <component
                :is="player.isPlaying ? Pause : Play"
                class="h-5 w-5 fill-current"
              />
            </Button>

            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 text-muted-foreground hover:text-primary transition-colors"
              @click="player.skip(10)"
            >
              <RotateCw class="h-4 w-4" />
            </Button>
          </div>

          <!-- Right: Volume (Wider & Distinct) -->
          <div
            class="hidden md:flex items-center gap-2 border-l border-border/50 pl-4 min-w-[140px] relative"
            @mouseenter="isHoveringVolume = true"
            @mouseleave="isHoveringVolume = false"
          >
            <!-- Custom Small Tooltip (Bottom) -->
            <div
              class="absolute -bottom-4 left-1/2 -translate-x-1/2 bg-slate-800/90 text-[8px] text-white font-bold px-1.5 py-0.5 rounded shadow-sm transition-all duration-200 pointer-events-none z-30"
              :class="
                isHoveringVolume
                  ? 'opacity-100 translate-y-1'
                  : 'opacity-0 translate-y-0'
              "
            >
              VOL: {{ Math.round(player.volume * 100) }}%
            </div>

            <Button
              variant="ghost"
              size="icon"
              class="h-8 w-8 text-muted-foreground hover:text-primary p-0 flex-shrink-0"
              @click="toggleMute"
            >
              <VolumeX v-if="isMuted || player.volume === 0" class="h-4 w-4" />
              <Volume2 v-else class="h-4 w-4" />
            </Button>
            <div class="flex-1 py-1">
              <Slider
                v-model="volumeValue"
                :max="100"
                :step="1"
                class="cursor-pointer volume-slider h-4"
              />
            </div>
          </div>
        </div>

        <!-- Progress Tracking Bar -->
        <div class="flex flex-col w-full relative">
          <!-- Time Labels -->
          <div
            class="flex justify-between px-4 pb-0.5 text-[9px] font-bold tabular-nums"
          >
            <span class="text-primary">{{
              formatDuration(player.currentTime)
            }}</span>
            <span class="text-muted-foreground/60">{{
              formatDuration(player.duration)
            }}</span>
          </div>

          <!-- The Actual Slider Container -->
          <div class="w-full relative px-0.5 pb-0.5 group/progress">
            <div
              class="h-1.5 w-full bg-muted/20 relative group-hover:h-2 transition-all duration-300 overflow-hidden rounded-b-2xl"
            >
              <Slider
                v-model="progressValue"
                :max="player.duration"
                :step="0.1"
                class="absolute -top-[7px] left-0 w-full z-20 cursor-pointer progress-slider opacity-0"
              />
              <!-- Custom Highlight Fill -->
              <div
                class="absolute top-0 left-0 h-full bg-primary transition-all duration-150 pointer-events-none rounded-r-full shadow-[0_0_10px_rgba(var(--primary),0.6)]"
                :style="{
                  width: `${(player.currentTime / player.duration) * 100}%`,
                }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
/* Reset Deep Selectors to be safer */
:deep([data-slot="slider-track"]) {
  background-color: transparent !important;
  height: 100% !important;
  width: 100% !important;
}

:deep([data-slot="slider-range"]) {
  background-color: var(--primary) !important;
  opacity: 1 !important;
}

:deep([data-slot="slider-thumb"]) {
  width: 12px !important;
  height: 12px !important;
  border: 2px solid white !important;
  background-color: var(--primary) !important;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3) !important;
  opacity: 1 !important;
  transition: transform 0.2s !important;
}

/* Volume-specific styling */
.volume-slider :deep([data-slot="slider-track"]) {
  background-color: hsl(var(--muted) / 0.2) !important;
  height: 4px !important;
  border-radius: 9999px;
  overflow: hidden;
}

/* Progress-slider specific: hide thumb until hover */
.progress-slider :deep([data-slot="slider-thumb"]) {
  opacity: 0 !important;
}

.group:hover .progress-slider :deep([data-slot="slider-thumb"]) {
  opacity: 1 !important;
}

.tabular-nums {
  font-variant-numeric: tabular-nums;
}

@keyframes pulse-subtle {
  0%,
  100% {
    transform: scale(1);
    opacity: 0.15;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.3;
  }
}

.animate-pulse {
  animation: pulse-subtle 2s ease-in-out infinite;
}
</style>
