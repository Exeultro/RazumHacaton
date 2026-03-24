<template>
  <div class="loader-overlay">
    <div class="loader-content">
      <div class="spinner-large"></div>
      <p>Загрузка данных...</p>
    </div>
  </div>
</template>

<style scoped>
.loader-overlay {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  display: flex; align-items: center; justify-content: center;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(4px); /* Эффект размытия фона */
  z-index: 9999;
}
.loader-content {
  display: flex; flex-direction: column; align-items: center;
  background: white; padding: 24px 40px; border-radius: 16px;
  box-shadow: var(--shadow-md); border: 1px solid var(--border-color);
}
.spinner-large {
  width: 36px; height: 36px; border: 3px solid var(--primary-light);
  border-top-color: var(--primary); border-radius: 50%;
  animation: spin 1s infinite linear; margin-bottom: 12px;
}
.loader-content p { color: var(--text-main); font-weight: 500; font-size: 0.9rem; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>