<template>
  <div class="leaderboard-view">
    <div class="page-header">
      <h1 class="page-title">Рейтинг участников</h1>
      <p class="page-description">Топ-100 активных участников платформы</p>
    </div>

    <LeaderboardFilters
        :currentDirection="ratingStore.currentDirection"
        @update:direction="ratingStore.setDirection"
    />

    <div v-if="ratingStore.loading" class="loading-container">
      <div class="spinner"></div>
      <p>Загрузка рейтинга...</p>
    </div>

    <LeaderboardTable
        v-else
        :leaders="ratingStore.globalRating"
        :userRank="ratingStore.userRank"
    />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRatingStore } from '@/stores/ratingStore'
import LeaderboardFilters from '@/components/rating/LeaderboardFilters.vue'
import LeaderboardTable from '@/components/rating/LeaderboardTable.vue'

const ratingStore = useRatingStore()

onMounted(() => {
  ratingStore.fetchLeaders()
})
</script>

<style scoped>
.leaderboard-view {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #1d4ed8 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-description {
  color: #64748b;
  margin: 0;
}

.loading-container {
  text-align: center;
  padding: 60px;
  background: white;
  border-radius: 16px;
  margin-top: 24px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e2e8f0;
  border-top-color: #1d4ed8;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>