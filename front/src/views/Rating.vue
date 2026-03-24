<template>
    <div class="rating">
      <h1>Таблица лидеров</h1>
      <el-select v-model="direction" placeholder="Направление" clearable @change="loadRating">
        <el-option label="IT" value="IT" />
        <el-option label="Социальное" value="social" />
        <el-option label="Медиа" value="media" />
      </el-select>
      <el-table :data="ratingStore.globalRating" v-loading="ratingStore.loading" stripe>
        <el-table-column prop="rank" label="Место" width="80" />
        <el-table-column prop="full_name" label="Участник" />
        <el-table-column prop="total_points" label="Баллы" sortable />
        <el-table-column prop="events_count" label="Мероприятий" />
        <el-table-column prop="direction" label="Направление" />
      </el-table>
      <div v-if="ratingStore.userRank" class="my-rank">
        <h3>Мой рейтинг: {{ ratingStore.userRank.rank }} место ({{ ratingStore.userRank.total_points }} баллов)</h3>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRatingStore } from '../stores/ratingStore'
  import { useAuthStore } from '../stores/auth'
  
  const ratingStore = useRatingStore()
  const authStore = useAuthStore()
  const direction = ref(null)
  
  const loadRating = () => {
    ratingStore.fetchGlobalRating(direction.value)
  }
  
  onMounted(() => {
    loadRating()
  })
  </script>