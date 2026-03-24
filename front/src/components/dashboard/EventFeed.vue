<template>
  <div class="event-feed">
    <div v-if="!events.length" class="empty-state">
      <span class="empty-icon">📭</span>
      <p>Нет событий</p>
    </div>
    
    <div v-else class="events-list">
      <div v-for="event in events" :key="event.id" class="event-item">
        <div class="event-icon" :class="getEventType(event)">
          {{ getEventIcon(event) }}
        </div>
        <div class="event-content">
          <div class="event-title">{{ event.title }}</div>
          <div class="event-meta">
            <span class="event-organizer">{{ event.organizer_name }}</span>
            <span class="event-date">{{ formatDate(event.event_date) }}</span>
            <span class="event-participants">👥 {{ event.participants_count }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  events: {
    type: Array,
    default: () => []
  }
})

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)
  
  if (date.toDateString() === today.toDateString()) {
    return `Сегодня в ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
  } else if (date.toDateString() === yesterday.toDateString()) {
    return `Вчера в ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
  } else {
    return `${date.getDate()}.${date.getMonth() + 1} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
  }
}

const getEventIcon = (event) => {
  if (event.title?.toLowerCase().includes('хакатон')) return '💻'
  if (event.title?.toLowerCase().includes('лекция')) return '📚'
  if (event.title?.toLowerCase().includes('фестиваль')) return '🎉'
  if (event.title?.toLowerCase().includes('форум')) return '🗣️'
  return '📅'
}

const getEventType = (event) => {
  if (event.title?.toLowerCase().includes('хакатон')) return 'hackathon'
  if (event.title?.toLowerCase().includes('лекция')) return 'lecture'
  return 'default'
}
</script>

<style scoped>
.event-feed {
  max-height: 400px;
  overflow-y: auto;
}

.events-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.event-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 12px;
  transition: all 0.2s;
  border: 1px solid #e2e8f0;
}

.event-item:hover {
  background: #f1f5f9;
  transform: translateX(4px);
}

.event-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.event-icon.hackathon {
  background: linear-gradient(135deg, #eef2ff 0%, #e0e7ff 100%);
}

.event-icon.lecture {
  background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
}

.event-icon.default {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
}

.event-content {
  flex: 1;
}

.event-title {
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 4px;
  font-size: 14px;
}

.event-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #64748b;
  flex-wrap: wrap;
}

.event-organizer {
  font-weight: 500;
}

.event-date, .event-participants {
  display: flex;
  align-items: center;
  gap: 4px;
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: #94a3b8;
}

.empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 12px;
  opacity: 0.5;
}
</style>