<template>
  <div class="organizer-page" v-loading="organizerStore.loading">
    <div v-if="organizerStore.organizer" class="content-wrapper">

      <div class="org-hero-card">
        <div class="org-avatar">
          {{ organizerStore.organizer.full_name.charAt(0).toUpperCase() }}
        </div>
        <div class="org-info">
          <h1 class="org-name">{{ organizerStore.organizer.full_name }}</h1>
          <div class="org-meta">
            <span class="meta-tag">📍 {{ organizerStore.organizer.city || 'Город не указан' }}</span>
            <span class="meta-tag">✉️ {{ organizerStore.organizer.email }}</span>
          </div>
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon purple">📅</div>
          <div class="stat-content">
            <div class="stat-value">{{ organizerStore.orgStats?.events_count || 0 }}</div>
            <div class="stat-label">Проведено мероприятий</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon orange">⭐</div>
          <div class="stat-content">
            <div class="stat-value">{{ organizerStore.reviewStats?.average_rating || organizerStore.orgStats?.rating || 0 }}</div>
            <div class="stat-label">Рейтинг доверия</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon green">🎁</div>
          <div class="stat-content">
            <div class="stat-value prize-list">
              {{ organizerStore.orgStats?.common_prizes?.length ? organizerStore.orgStats.common_prizes.join(', ') : 'Нет данных' }}
            </div>
            <div class="stat-label">Частые призы</div>
          </div>
        </div>
      </div>

      <div class="reviews-section">
        <div class="reviews-sidebar">
          <div class="rating-summary card-pro">
            <h3>Сводка отзывов</h3>
            <div class="avg-rating-block">
              <div class="avg-number">{{ organizerStore.reviewStats?.average_rating || 0 }}</div>
              <el-rate :model-value="organizerStore.reviewStats?.average_rating || 0" disabled allow-half />
              <div class="total-reviews">{{ organizerStore.reviewStats?.total_reviews || 0 }} отзывов</div>
            </div>

            <div class="rating-bars" v-if="organizerStore.reviewStats && organizerStore.reviewStats.total_reviews > 0">
              <div v-for="star in 5" :key="star" class="bar-row">
                <span class="star-label">{{ 6 - star }} ⭐</span>
                <div class="bar-bg">
                  <div class="bar-fill" :style="{ width: getRatingPercentage(6 - star) + '%' }"></div>
                </div>
                <span class="bar-count">{{ organizerStore.reviewStats[`rating_${6 - star}_count`] || 0 }}</span>
              </div>
            </div>
          </div>

          <div v-if="isParticipant" class="add-review-card card-pro mt-4">
            <h3>Оставить отзыв</h3>

            <div v-if="selectedEventTitle" class="selected-event-alert">
              <div class="alert-content">
                <span class="alert-label">По мероприятию:</span>
                <span class="alert-title">{{ selectedEventTitle }}</span>
              </div>
              <button class="clear-event-btn" @click="clearSelectedEvent" title="Оставить общий отзыв">✕</button>
            </div>

            <el-form :model="newReview" @submit.prevent="submitReview">
              <div class="form-group">
                <label>Оценка</label>
                <el-rate v-model="newReview.rating" :colors="['#99A9BF', '#F7BA2A', '#FF9900']" />
              </div>
              <div class="form-group">
                <label>Комментарий</label>
                <el-input
                    type="textarea"
                    v-model="newReview.comment"
                    rows="3"
                    placeholder="Поделитесь своими впечатлениями..."
                    resize="none"
                />
              </div>
              <el-button type="primary" native-type="submit" :loading="organizerStore.submitting" class="w-100" :disabled="!newReview.rating || !newReview.comment">
                Отправить отзыв
              </el-button>
            </el-form>
          </div>
        </div>

        <div class="reviews-list-container">
          <h3 class="section-title">Отзывы участников</h3>

          <div v-if="organizerStore.reviews.length > 0" class="reviews-list">
            <div v-for="rev in organizerStore.reviews" :key="rev.id" class="review-card">
              <div class="review-header">
                <div class="reviewer-avatar">
                  {{ rev.participant_name.charAt(0).toUpperCase() }}
                </div>
                <div class="reviewer-info">
                  <div class="reviewer-name">{{ rev.participant_name }}</div>
                  <div class="review-date">{{ formatDate(rev.created_at) }}</div>
                </div>
                <div class="review-rating">
                  <el-rate v-model="rev.rating" disabled />
                </div>
              </div>

              <div v-if="rev.event_title" class="review-event">
                <span class="event-badge">Мероприятие: {{ rev.event_title }}</span>
              </div>

              <p class="review-comment">{{ rev.comment }}</p>
            </div>
          </div>

          <div v-else class="empty-state-card">
            <div class="empty-icon">📝</div>
            <p>Об этом организаторе пока нет отзывов.</p>
            <p v-if="isParticipant" class="text-muted">Будьте первым, кто поделится впечатлениями!</p>
          </div>

          <div class="pagination-center" v-if="organizerStore.pagination && organizerStore.pagination.total > organizerStore.pagination.limit">
            <el-pagination
                :current-page="organizerStore.pagination.page"
                :page-size="organizerStore.pagination.limit"
                :total="organizerStore.pagination.total"
                @current-change="handlePageChange"
                layout="prev, pager, next"
                background
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useOrganizerStore } from '@/stores/organizerStore'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const organizerStore = useOrganizerStore()
const organizerId = route.params.id

const selectedEventTitle = ref(route.query.event_title || '')

const newReview = ref({
  rating: 0,
  comment: '',
  event_id: route.query.event_id || null
})

const isParticipant = computed(() => authStore.userRole === 'participant')

const clearSelectedEvent = () => {
  selectedEventTitle.value = ''
  newReview.value.event_id = null
  router.replace({ query: {} })
}

const submitReview = async () => {
  if (!newReview.value.rating || !newReview.value.comment) return

  const success = await organizerStore.submitReview(organizerId, {
    rating: newReview.value.rating,
    comment: newReview.value.comment,
    event_id: newReview.value.event_id || undefined
  })

  if (success) {
    newReview.value.rating = 0
    newReview.value.comment = ''
  }
}

const handlePageChange = (newPage) => {
  organizerStore.fetchReviews(organizerId, newPage)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const formatDate = (iso) => {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'long', year: 'numeric'
  })
}

const getRatingPercentage = (starLevel) => {
  if (!organizerStore.reviewStats || organizerStore.reviewStats.total_reviews === 0) return 0
  const count = organizerStore.reviewStats[`rating_${starLevel}_count`] || 0
  return (count / organizerStore.reviewStats.total_reviews) * 100
}

onMounted(async () => {
  await Promise.all([
    organizerStore.fetchProfile(organizerId),
    organizerStore.fetchReviews(organizerId, 1)
  ])

  if (route.query.event_id && isParticipant.value) {
    setTimeout(() => {
      window.scrollTo({
        top: document.querySelector('.add-review-card')?.offsetTop - 20 || 0,
        behavior: 'smooth'
      })
    }, 500)
  }
})
</script>

<style scoped>
.organizer-page {
  font-family: 'Inter', sans-serif;
  color: #111827;
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.card-pro {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.w-100 { width: 100%; }
.mt-4 { margin-top: 24px; }
.text-muted { color: #6b7280; font-size: 0.9rem; }

.org-hero-card {
  background: linear-gradient(135deg, #1e3a8a 0%, #3b82f6 100%);
  border-radius: 20px;
  padding: 40px;
  display: flex;
  align-items: center;
  gap: 24px;
  color: white;
  margin-bottom: 24px;
  box-shadow: 0 10px 25px -5px rgba(59, 130, 246, 0.4);
}

.org-avatar {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.4);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 700;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.org-name {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 800;
  letter-spacing: -0.5px;
}

.org-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.meta-tag {
  background: rgba(0, 0, 0, 0.2);
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.stat-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}
.stat-icon.purple { background: #f3e8ff; border: 1px solid #d8b4fe; }
.stat-icon.orange { background: #fff7ed; border: 1px solid #fdba74; }
.stat-icon.green { background: #dcfce7; border: 1px solid #86efac; }

.stat-value {
  font-size: 24px;
  font-weight: 800;
  color: #111827;
}
.stat-value.prize-list {
  font-size: 15px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}
.stat-label {
  font-size: 13px;
  color: #6b7280;
  font-weight: 500;
  margin-top: 2px;
}

.reviews-section {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 24px;
}

@media (max-width: 900px) {
  .reviews-section {
    grid-template-columns: 1fr;
  }
}

.rating-summary h3, .add-review-card h3, .section-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  font-weight: 700;
}

.avg-rating-block {
  text-align: center;
  margin-bottom: 24px;
}
.avg-number {
  font-size: 48px;
  font-weight: 800;
  line-height: 1;
  margin-bottom: 8px;
}
.total-reviews {
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.bar-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}
.star-label {
  font-size: 13px;
  font-weight: 600;
  color: #4b5563;
  width: 35px;
}
.bar-bg {
  flex-grow: 1;
  height: 8px;
  background: #f3f4f6;
  border-radius: 4px;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  background: #f59e0b;
  border-radius: 4px;
  transition: width 0.5s ease;
}
.bar-count {
  font-size: 13px;
  color: #6b7280;
  width: 20px;
  text-align: right;
}

.form-group {
  margin-bottom: 16px;
}
.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.reviews-list-container {
  display: flex;
  flex-direction: column;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.review-card {
  background: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
  transition: box-shadow 0.2s ease;
}

.review-card:hover {
  box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -1px rgba(0,0,0,0.06);
}

.review-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.reviewer-avatar {
  width: 40px;
  height: 40px;
  background: #e0e7ff;
  color: #4338ca;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 16px;
}

.reviewer-info {
  flex-grow: 1;
}

.reviewer-name {
  font-weight: 700;
  font-size: 15px;
  color: #111827;
}

.review-date {
  font-size: 12px;
  color: #9ca3af;
}

.review-event {
  margin-bottom: 12px;
}
.event-badge {
  background: #f1f5f9;
  color: #475569;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  display: inline-block;
}

.review-comment {
  margin: 0;
  font-size: 15px;
  line-height: 1.5;
  color: #374151;
}

.empty-state-card {
  background: #ffffff;
  border: 1px dashed #d1d5db;
  border-radius: 16px;
  padding: 60px 20px;
  text-align: center;
}
.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.pagination-center {
  margin-top: 32px;
  display: flex;
  justify-content: center;
}

.selected-event-alert {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.alert-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.alert-label {
  font-size: 12px;
  color: #3b82f6;
  text-transform: uppercase;
  font-weight: 700;
}

.alert-title {
  font-size: 14px;
  font-weight: 600;
  color: #1e3a8a;
}

.clear-event-btn {
  background: none;
  border: none;
  color: #93c5fd;
  font-size: 16px;
  cursor: pointer;
  padding: 4px;
  border-radius: 50%;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
}

.clear-event-btn:hover {
  background-color: #dbeafe;
  color: #2563eb;
}
</style>