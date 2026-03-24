<template>
  <div class="notifications-container">
    <transition-group name="toast">
      <div
          v-for="notif in uiStore.notifications"
          :key="notif.id"
          class="toast"
          :class="`toast-${notif.type}`"
      >
        <span>{{ notif.message }}</span>
        <button class="close-btn" @click="uiStore.removeNotification(notif.id)">×</button>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { useUiStore } from '@/stores/ui.js';

const uiStore = useUiStore();
</script>

<style scoped>
.notifications-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  z-index: 10000;
  pointer-events: none;
}

.toast {
  pointer-events: auto;
  min-width: 250px;
  padding: 12px 16px;
  border-radius: 6px;
  background: #333;
  color: #fff;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
}

.toast-success { background: #28a745; }
.toast-error { background: #dc3545; }
.toast-info { background: #17a2b8; }
.toast-warning { background: #ffc107; color: #333; }

.close-btn {
  background: transparent;
  border: none;
  color: inherit;
  font-size: 1.2rem;
  line-height: 1;
  cursor: pointer;
  margin-left: 15px;
  opacity: 0.8;
}

.close-btn:hover {
  opacity: 1;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
</style>