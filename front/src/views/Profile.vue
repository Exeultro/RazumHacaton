<template>
    <div class="profile-page">
      <h1 class="page-title">Личный кабинет</h1>
      <div v-if="authStore.user" class="card">
        <h3 class="section-title">Данные профиля</h3>
        <p><strong>ФИО:</strong> {{ authStore.user.full_name }}</p>
        <p><strong>Роль:</strong> {{ authStore.user.role }}</p>
        <p><strong>Email:</strong> {{ authStore.user.email }}</p>
        <p v-if="authStore.user.city"><strong>Город:</strong> {{ authStore.user.city }}</p>
        <p v-if="authStore.user.age"><strong>Возраст:</strong> {{ authStore.user.age }}</p>
        <p v-if="authStore.user.direction"><strong>Направление:</strong> {{ authStore.user.direction }}</p>
      </div>
  
      <div v-if="authStore.user?.role === 'participant'" class="card">
        <h3 class="section-title">Мой рейтинг</h3>
        <p>Баллы: {{ ratingStore.userRating?.total_points || 0 }}</p>
        <p>Место в общем зачёте: {{ ratingStore.userRank?.rank || '—' }}</p>
      </div>
  
      <div v-if="authStore.user?.role === 'organizer'" class="card">
        <h3 class="section-title">Статистика организатора</h3>
        <p>Проведено мероприятий: {{ authStore.user.organizer_stats?.events_count || 0 }}</p>
        <p>Рейтинг доверия: {{ authStore.user.organizer_stats?.trust_rating || 0 }}</p>
      </div>
  
      <div class="actions" style="margin-top: 24px; display: flex; gap: 12px;">
        <router-link to="/events">
          <button class="btn btn-primary">Каталог мероприятий</button>
        </router-link>
        <router-link v-if="authStore.userRole === 'organizer'" to="/events/create">
          <button class="btn btn-success">Создать мероприятие</button>
        </router-link>
      </div>
    </div>
  </template>
  
  <script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRatingStore } from '../stores/ratingStore'
  
  const authStore = useAuthStore()
  const ratingStore = useRatingStore()
  
  onMounted(async () => {
    if (authStore.user?.role === 'participant') {
      await ratingStore.fetchGlobalRating()
    }
  })
  </script>