<template>
  <div class="avatar" :class="size">
    <img v-if="url" :src="url" alt="Avatar" />
    <span v-else>{{ initials }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  url: String,
  name: { type: String, default: 'User' },
  size: { type: String, default: 'md' }
});

const initials = computed(() => {
  return props.name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase();
});
</script>

<style scoped>
.avatar {
  border-radius: 50%;
  background: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  font-weight: 600;
  flex-shrink: 0;
}
.avatar img { width: 100%; height: 100%; object-fit: cover; }
.sm { width: 32px; height: 32px; font-size: 0.8rem; }
.md { width: 48px; height: 48px; font-size: 1.1rem; }
.lg { width: 80px; height: 80px; font-size: 2rem; }
</style>