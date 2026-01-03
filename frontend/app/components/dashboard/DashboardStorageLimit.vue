<script setup lang="ts">
defineProps<{
  summary: any;
  storagePercentage: number;
  storageBytes: number;
}>();
</script>

<template>
  <Card class="col-span-3">
    <CardHeader>
      <CardTitle>Storage Limit</CardTitle>
      <CardDescription>
        Visual distribution of your 10GB quota.
      </CardDescription>
    </CardHeader>
    <CardContent>
      <div
        class="flex flex-col items-center justify-center h-[200px] space-y-4"
      >
        <div class="relative flex items-center justify-center">
          <svg class="h-32 w-32 transform -rotate-90">
            <circle
              cx="64"
              cy="64"
              r="58"
              stroke="currentColor"
              stroke-width="8"
              fill="transparent"
              class="text-muted/20"
            />
            <circle
              cx="64"
              cy="64"
              r="58"
              stroke="currentColor"
              stroke-width="8"
              fill="transparent"
              :stroke-dasharray="2 * Math.PI * 58"
              :stroke-dashoffset="
                2 * Math.PI * 58 * (1 - storagePercentage / 100)
              "
              class="text-purple-500 transition-all duration-500"
            />
          </svg>
          <div class="absolute flex flex-col items-center justify-center">
            <span class="text-2xl font-bold">
              {{ Math.round(storagePercentage) }}%
            </span>
            <span class="text-[10px] uppercase text-muted-foreground">
              Used
            </span>
          </div>
        </div>
        <div class="text-center">
          <p class="text-sm font-medium">
            {{ summary?.data?.total_size ?? "0 B" }} / 10 GB
          </p>
          <p class="text-xs text-muted-foreground">
            {{ (10 - storageBytes / 1024 ** 3).toFixed(2) }} GB remaining
          </p>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
