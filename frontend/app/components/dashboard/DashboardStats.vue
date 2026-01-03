<script setup lang="ts">
import { Clock, HardDrive, Music } from "lucide-vue-next";
import { formatTimeAgo } from "~/lib/format";

const props = defineProps<{
  summary: any;
  pending: boolean;
  storagePercentage: number;
}>();

const stats = computed(() => [
  {
    title: "Total Songs",
    value: props.summary?.data?.total_songs ?? 0,
    unit: "Tracks",
    icon: Music,
    color: "text-blue-500",
    bg: "bg-blue-500/10",
  },
  {
    title: "Storage Used",
    value: props.summary?.data?.total_size ?? "0 B",
    unit: "of 10 GB limit",
    icon: HardDrive,
    color: "text-purple-500",
    bg: "bg-purple-500/10",
    progress: props.storagePercentage,
  },
  {
    title: "Last Upload",
    value: formatTimeAgo(props.summary?.data?.last_upload),
    unit: "Recently added",
    icon: Clock,
    color: "text-orange-500",
    bg: "bg-orange-500/10",
  },
]);
</script>

<template>
  <div class="grid gap-4 md:grid-cols-3">
    <Card v-for="stat in stats" :key="stat.title" class="overflow-hidden">
      <CardHeader
        class="flex flex-row items-center justify-between space-y-0 pb-2"
      >
        <CardTitle class="text-sm font-medium">
          {{ stat.title }}
        </CardTitle>
        <div :class="[stat.bg, stat.color]" class="p-2 rounded-md">
          <component :is="stat.icon" class="h-4 w-4" />
        </div>
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">
          <Skeleton v-if="pending" class="h-8 w-24" />
          <template v-else>
            {{ stat.value }}
          </template>
        </div>
        <div class="mt-1 space-y-2">
          <p class="text-xs text-muted-foreground">
            {{ stat.unit }}
          </p>
          <Progress
            v-if="stat.progress !== undefined"
            :model-value="stat.progress"
            class="h-1"
          />
        </div>
      </CardContent>
    </Card>
  </div>
</template>
