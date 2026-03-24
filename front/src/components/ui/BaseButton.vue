<template>
  <button
      :type="type"
      :class="['btn', `btn-${variant}`, { 'is-loading': loading, 'btn-block': block }]"
      :disabled="disabled || loading"
      @click="$emit('click', $event)"
  >
    <span v-if="loading" class="spinner"></span>
    <span :class="['btn-content', { 'is-hidden': loading }]">
      <slot></slot>
    </span>
  </button>
</template>

<script setup>
defineProps({
  type: { type: String, default: 'button' },
  variant: { type: String, default: 'primary' },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  block: { type: Boolean, default: false }
});

defineEmits(['click']);
</script>

<style scoped>
.btn {
  color: var(--text-main);
  position: relative;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  height: 44px; /* Единая высота для кнопок и инпутов */
  padding: 0 20px;
  border-radius: var(--radius-sm);
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s ease;
  white-space: nowrap;
  user-select: none;
}

/* Состояния наведения и нажатия */
.btn:active:not(:disabled) {
  transform: scale(0.98);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Модификаторы */
.btn-block {
  width: 100%;
}

/* --- Варианты кнопок --- */
.btn-primary {
  background-color: var(--primary);
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
}
.btn-primary:hover:not(:disabled) {
  background-color: var(--primary-hover);
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.3);
}

.btn-secondary {
  background-color: #f1f5f9;
  color: var(--text-main);
  border-color: #e2e8f0;
}
.btn-secondary:hover:not(:disabled) {
  background-color: #e2e8f0;
}

.btn-danger {
  background-color: var(--danger);
  color: #ffffff;
}
.btn-danger:hover:not(:disabled) {
  background-color: #dc2626;
}

/* --- Состояние загрузки --- */
.is-loading {
  pointer-events: none;
}
.is-hidden {
  opacity: 0;
}

/* Спиннер наследует цвет текста (белый для primary, темный для secondary) */
.spinner {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 20px;
  height: 20px;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
  animation: spin 0.75s linear infinite;
}

@keyframes spin {
  from { transform: translate(-50%, -50%) rotate(0deg); }
  to { transform: translate(-50%, -50%) rotate(360deg); }
}
</style>