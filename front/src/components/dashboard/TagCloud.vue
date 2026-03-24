<template>
  <div class="tag-cloud">
    <div
      v-for="tag in tags"
      :key="tag.tag"
      class="tag-item"
      :style="{
        fontSize: getFontSize(tag.count) + 'px',
        opacity: getOpacity(tag.count)
      }"
    >
      {{ tag.tag }}
      <span class="tag-count">{{ tag.count }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  tags: {
    type: Array,
    default: () => []
  }
})

const maxCount = computed(() => {
  if (!props.tags.length) return 1
  return Math.max(...props.tags.map(t => t.count))
})

const minCount = computed(() => {
  if (!props.tags.length) return 1
  return Math.min(...props.tags.map(t => t.count))
})

const getFontSize = (count) => {
  const minSize = 14
  const maxSize = 36
  if (maxCount.value === minCount.value) return minSize + (maxSize - minSize) / 2
  const ratio = (count - minCount.value) / (maxCount.value - minCount.value)
  return minSize + ratio * (maxSize - minSize)
}

const getOpacity = (count) => {
  const minOpacity = 0.6
  const maxOpacity = 1
  if (maxCount.value === minCount.value) return (minOpacity + maxOpacity) / 2
  const ratio = (count - minCount.value) / (maxCount.value - minCount.value)
  return minOpacity + ratio * (maxOpacity - minOpacity)
}
</script>

<style scoped>
.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding: 8px;
  min-height: 250px;
  align-items: center;
  justify-content: center;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 40px;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 500;
  color: #1e293b;
  border: 1px solid #e2e8f0;
}

.tag-item:hover {
  transform: scale(1.05);
  background: linear-gradient(135deg, #eef2ff 0%, #e0e7ff 100%);
  border-color: #1d4ed8;
  color: #1d4ed8;
}

.tag-count {
  font-size: 11px;
  background: rgba(0, 0, 0, 0.05);
  padding: 2px 6px;
  border-radius: 20px;
  font-weight: 600;
  color: #64748b;
}

.tag-item:hover .tag-count {
  background: rgba(29, 78, 216, 0.1);
  color: #1d4ed8;
}
</style>