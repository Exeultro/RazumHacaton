<template>
  <div class="filters-container">
    <span class="filter-label">Направление:</span>
    <div class="filter-buttons">
      <button
          v-for="option in directionOptions"
          :key="option.value"
          class="filter-btn"
          :class="{ active: currentDirection === option.value }"
          @click="$emit('update:direction', option.value)"
      >
        {{ option.label }}
      </button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  currentDirection: {
    type: String,
    default: 'all'
  }
})

defineEmits(['update:direction'])

const directionOptions = [
  { value: 'all', label: 'Все' },
  { value: 'IT', label: 'IT' },
  { value: 'social', label: 'Социальное' },
  { value: 'media', label: 'Медиа' }
]
</script>


<style scoped>
.filters-container {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  margin-bottom: 24px;
  padding: 16px 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.filter-label {
  font-weight: 600;
  color: #0f172a;
  font-size: 0.875rem;
  letter-spacing: 0.3px;
}

.filter-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 8px 20px;
  border: 1px solid #e2e8f0;
  background: white;
  border-radius: 40px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s;
  color: #475569;
}

.filter-btn:hover {
  border-color: #1d4ed8;
  color: #1d4ed8;
  transform: translateY(-1px);
}

.filter-btn.active {
  background: #1d4ed8;
  color: white;
  border-color: #1d4ed8;
  box-shadow: 0 2px 8px rgba(29, 78, 216, 0.2);
}
</style>