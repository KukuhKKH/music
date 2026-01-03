<script setup lang="ts">
defineProps<{
  pending: boolean;
}>();
</script>

<template>
  <Card class="col-span-4">
    <CardHeader>
      <CardTitle>Recent Activity</CardTitle>
      <CardDescription>
        You have uploaded 3 new tracks this week.
      </CardDescription>
    </CardHeader>
    <CardContent>
      <div class="space-y-4">
        <Skeleton v-if="pending" class="h-[200px] w-full" />
        <div
          v-else
          class="h-[200px] flex items-center justify-center border-2 border-dashed rounded-lg border-muted"
        >
          <p class="text-muted-foreground text-sm italic">
            Activity chart coming soon...
          </p>
        </div>
      </div>
    </CardContent>
  </Card>
</template>
