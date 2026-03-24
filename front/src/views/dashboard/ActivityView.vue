<template>
  <div class="activity-view">
    <div class="page-header">
      <div>
        <h1 class="page-title">📊 Дашборд активности</h1>
        <p class="text-muted">Аналитика участия и популярные направления платформы</p>
      </div>

      <button
          class="btn btn-outline"
          @click="analyticsStore.fetchActivity()"
          :disabled="analyticsStore.isLoading"
      >
        <span v-if="analyticsStore.isLoading">Обновление...</span>
        <span v-else>🔄 Обновить</span>
      </button>
    </div>

    <!-- Общий лоадер -->
    <div v-if="analyticsStore.isLoading && !analyticsStore.ratingHistory.length" class="loading-container card">
      <div class="spinner"></div>
      <p class="text-muted">Собираем аналитику...</p>
    </div>

    <div v-else class="dashboard-grid">
      <!-- График -->
      <div class="card chart-card">
        <div class="card-header">
          <h3>📈 Динамика вовлеченности</h3>
          <span class="badge badge-primary">За 30 дней</span>
        </div>
        <div v-if="analyticsStore.ratingHistory.length > 0">
          <ActivityChart :data="analyticsStore.ratingHistory" />
        </div>
        <div v-else class="empty-data">Нет данных для графика</div>
      </div>

      <!-- Теги -->
      <div class="card tags-card">
        <div class="card-header">
          <h3>🏷️ Популярные форматы</h3>
        </div>
        <div v-if="analyticsStore.trendingTags.length > 0">
          <TagCloud :tags="analyticsStore.trendingTags" />
        </div>
        <div v-else class="empty-data">Теги не найдены</div>
      </div>

      <!-- Лента -->
      <div class="card feed-card">
        <div class="card-header">
          <h3>🔄 Последние события</h3>
        </div>
        <div v-if="analyticsStore.recentEvents.length > 0">
          <EventFeed :events="analyticsStore.recentEvents" />
        </div>
        <div v-else class="empty-data">Лента пуста</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAnalyticsStore } from '@/stores/analyticsStore'
import ActivityChart from '@/components/dashboard/ActivityChart.vue'
import TagCloud from '@/components/dashboard/TagCloud.vue'
import EventFeed from '@/components/dashboard/EventFeed.vue'

const analyticsStore = useAnalyticsStore()

onMounted(() => {
  if (analyticsStore.recentEvents.length === 0) {
    analyticsStore.fetchActivity()
  }
})
</script>

<style scoped>
.activity-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.page-title { margin-bottom: 8px; }

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.chart-card { grid-column: span 2; }

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
}

.card-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-main);
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  min-height: 300px;
}

.empty-data {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-muted);
  background: var(--bg-app);
  border-radius: var(--radius-sm);
  font-size: 0.9rem;
}

@media (max-width: 900px) {
  .dashboard-grid { grid-template-columns: 1fr; }
  .chart-card { grid-column: span 1; }
}

@media (max-width: 600px) {
  .page-header { flex-direction: column; align-items: flex-start; }
}
</style>