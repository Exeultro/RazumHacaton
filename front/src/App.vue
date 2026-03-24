<template>
  <div id="app-root">
    <div v-if="uiStore.globalLoading" class="global-loader">
      <span class="spinner">Загрузка...</span>
    </div>

    <router-view v-else />

    <GlobalNotifications />
  </div>
</template>

<script setup>
import { useUiStore } from './stores/ui';
import GlobalNotifications from "@/components/ui/GlobalNotifications.vue";

const uiStore = useUiStore();
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  background-color: #f4f7f6;
  color: #333;
}

* {
  box-sizing: border-box;
}

#app-root {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.global-loader {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.9);
  z-index: 9999;
  font-size: 1.2rem;
  font-weight: 500;
  color: #007bff;
}
</style>