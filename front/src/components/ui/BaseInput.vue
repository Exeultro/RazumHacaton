<template>
  <div class="input-group">
    <label v-if="label" class="input-label">{{ label }}</label>

    <div class="input-wrapper">
      <input
          :type="type"
          :value="modelValue"
          :placeholder="placeholder"
          :disabled="disabled"
          @input="$emit('update:modelValue', $event.target.value)"
          :class="['base-input', { 'has-error': error }]"
      />
    </div>

    <!-- Анимация появления ошибки -->
    <transition name="fade">
      <span v-if="error" class="error-text">{{ error }}</span>
    </transition>
  </div>
</template>

<script setup>
defineProps({
  modelValue: [String, Number],
  label: String,
  type: { type: String, default: 'text' },
  placeholder: String,
  error: String,
  disabled: { type: Boolean, default: false }
});

defineEmits(['update:modelValue']);
</script>

<style scoped>
.input-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
  width: 100%;
}

.input-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-main);
}

.input-wrapper {
  position: relative;
  width: 100%;
}

.base-input {
  width: 100%;
  height: 44px; /* Синхронизировано с высотой кнопки */
  padding: 0 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  background-color: var(--bg-card);
  color: var(--text-main);
  font-family: inherit;
  font-size: 0.95rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  box-shadow: var(--shadow-sm);
}

/* Стилизация плейсхолдера */
.base-input::placeholder {
  color: var(--text-muted);
  opacity: 0.7;
}

/* Наведение и фокус */
.base-input:hover:not(:disabled) {
  border-color: #cbd5e1;
}

.base-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15); /* Мягкое свечение при фокусе */
}

/* Состояние Disabled */
.base-input:disabled {
  background-color: #f8fafc;
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: none;
}

/* Состояние Ошибки */
.has-error {
  border-color: var(--danger) !important;
}
.has-error:focus {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15) !important;
}

.error-text {
  font-size: 0.75rem;
  color: var(--danger);
  margin-top: 2px;
}

/* Плавное появление текста ошибки */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}
</style>