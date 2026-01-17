<script setup lang="ts">
import { parseSizeToBytes } from "~/lib/format";

const { $api } = useNuxtApp();
const {
  data: summary,
  pending,
  error,
  refresh,
} = await useAsyncData<any>("dashboard-summary", () => $api("/stats/summary"));

// Constants for storage
const MAX_STORAGE_BYTES = 10 * 1024 * 1024 * 1024; // 10GB

const storageBytes = computed(() =>
  parseSizeToBytes(summary.value?.data?.total_size ?? "0 B")
);

const storagePercentage = computed(() =>
  Math.min((storageBytes.value / MAX_STORAGE_BYTES) * 100, 100)
);
</script>

<template>
  <div class="p-6 space-y-8">
    <!-- Header Section -->
    <DashboardHeader :pending="pending" @refresh="refresh" />

    <!-- Error State -->
    <div
      v-if="error"
      class="p-4 border border-destructive/50 bg-destructive/10 text-destructive rounded-lg flex items-center gap-3"
    >
      <div class="i-lucide-alert-circle h-5 w-5" />
      <p>Failed to load statistics. Please try again later.</p>
    </div>

    <!-- Stats Cards Section -->
    <DashboardStats
      :summary="summary"
      :pending="pending"
      :storage-percentage="storagePercentage"
    />

    <!-- Charts & Activity Section -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <DashboardRecentActivity :pending="pending" />
      <DashboardStorageLimit
        :summary="summary"
        :storage-bytes="storageBytes"
        :storage-percentage="storagePercentage"
      />
    </div>
  </div>
</template>
