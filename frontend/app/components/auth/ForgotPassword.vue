<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";
import { toast } from "vue-sonner";

const isLoading = ref(false);

function handlePajangan() {
  toast.info("Hayooo ngapain, ini cuma pajangan aja yaaaaa");
}

async function onSubmit(event: Event) {
  event.preventDefault();
  handlePajangan();
}
</script>

<template>
  <form @submit="onSubmit">
    <div class="grid gap-4">
      <div class="grid gap-2">
        <Label for="email"> Email </Label>
        <Input
          id="email"
          placeholder="name@example.com"
          type="email"
          auto-capitalize="none"
          auto-complete="email"
          auto-correct="off"
          :disabled="isLoading"
        />
      </div>
      <Button :disabled="isLoading" type="submit">
        <Loader2 v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
        Submit
      </Button>
    </div>
  </form>
</template>
