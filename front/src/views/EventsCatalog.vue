<template>
  <div class="catalog-container">
    <div class="catalog-header">
      <div class="title-group">
        <h1 class="page-title">Мероприятия</h1>
        <span class="badge-total" v-if="eventStore.pagination">
          Всего: {{ eventStore.pagination.total }}
        </span>
      </div>
    </div>

    <!-- Фильтры -->
    <EventFilters @filter="applyFilters" class="filters-wrapper" />

    <div v-loading="eventStore.loading" class="events-grid">
      <article
          v-for="event in eventStore.events"
          :key="event.id"
          class="event-card-clean"
          @click="goToDetail(event.id)"
      >
        <!-- Верхняя часть с тегами -->
        <div class="card-header-tags">
          <div class="tag-row">
            <span class="tag tag-format">
              {{ translateFormat(event.format) }}
            </span>
            <span class="tag tag-direction">
              {{ translateDirection(event.direction) }}
            </span>
          </div>

          <div v-if="event.is_registered" class="status-going">
            <svg viewBox="0 0 24 24" width="14" height="14" stroke="currentColor" stroke-width="2.5" fill="none" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
            Вы записаны
          </div>
        </div>

        <!-- Тело карточки -->
        <div class="card-body">
          <h3 class="card-title" :title="event.title">{{ event.title }}</h3>

          <div class="info-list">
            <div class="info-row">
              <svg class="info-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
              <span class="info-text">{{ event.format === 'online' ? 'Онлайн' : (event.location || 'Место не указано') }}</span>
            </div>
            <div class="info-row">
              <svg class="info-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
              <span class="info-text">{{ formatTimeRange(event.event_date, event.end_date) }}</span>
            </div>
          </div>

          <!-- Баллы -->
          <div class="points-wrap" v-if="event.points_for_participation">
            <svg class="points-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
            <span class="points-text"><strong>{{ event.points_for_participation }}</strong> баллов за участие</span>
          </div>
        </div>

        <!-- Подвал карточки: Организатор и Участники -->
        <div class="card-footer">
          <div class="org-block">
            <div class="org-name" :title="event.organizer?.full_name">
              {{ event.organizer?.full_name || 'Организатор не указан' }}
            </div>
            <div class="org-rating" v-if="event.organizer?.trust_rating !== undefined">
              <svg class="rating-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
              Рейтинг: <span>{{ event.organizer.trust_rating }}</span>
            </div>
          </div>

          <div class="users-block" title="Участники">
            <svg viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
            {{ event.participants_count || 0 }}
          </div>
        </div>
      </article>

      <!-- Состояние: Ничего не найдено -->
      <div v-if="!eventStore.loading && eventStore.events.length === 0" class="empty-state">
        <svg class="empty-state-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect><line x1="3" y1="9" x2="21" y2="9"></line><line x1="9" y1="21" x2="9" y2="9"></line></svg>
        <h3>Ничего не найдено</h3>
        <p>Попробуйте сбросить фильтры или выбрать другие даты.</p>
      </div>
    </div>

    <!-- Пагинация -->
    <div class="pagination-center" v-if="eventStore.pagination && eventStore.pagination.total > 0">
      <el-pagination
          :current-page="eventStore.filters.page"
          :page-size="eventStore.filters.limit"
          :total="eventStore.pagination.total"
          @current-change="handlePageChange"
          layout="prev, pager, next"
          background
      />
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useEventStore } from '../stores/eventStore'
import { useAuthStore } from '../stores/auth'
import EventFilters from '../components/EventFilters.vue'

const router = useRouter()
const eventStore = useEventStore()
const authStore = useAuthStore()

const isOrganizer = computed(() => authStore.userRole === 'organizer')

onMounted(() => {
  eventStore.fetchEvents()
})

const applyFilters = (filters) => {
  eventStore.updateFilters({ ...filters, page: 1 })
}

const handlePageChange = (page) => {
  eventStore.updateFilters({ page })
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const goToDetail = (id) => {
  router.push(`/events/${id}`)
}

const translateFormat = (f) => ({ online: 'Онлайн', offline: 'Офлайн', hybrid: 'Гибрид' }[f] || f || 'Формат')
const translateDirection = (d) => ({ media: 'Медиа', it: 'IT', sport: 'Спорт', science: 'Наука', art: 'Творчество' }[d] || d || 'Направление')

const formatTimeRange = (startIso, endIso) => {
  if (!startIso) return 'Время не указано'

  const start = new Date(startIso)
  const end = endIso ? new Date(endIso) : new Date(start.getTime() + 2 * 60 * 60 * 1000)

  const dateOpts = { day: 'numeric', month: 'short', timeZone: 'Europe/Moscow' }
  const timeOpts = { hour: '2-digit', minute: '2-digit', timeZone: 'Europe/Moscow' }

  const startDate = start.toLocaleDateString('ru-RU', dateOpts)
  const startTime = start.toLocaleTimeString('ru-RU', timeOpts)

  const endDate = end.toLocaleDateString('ru-RU', dateOpts)
  const endTime = end.toLocaleTimeString('ru-RU', timeOpts)

  if (startDate === endDate) {
    return `${startDate}, ${startTime} – ${endTime} (МСК)`
  }

  return `${startDate}, ${startTime} – ${endDate}, ${endTime} (МСК)`
}
</script>

<style scoped>
.catalog-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
}

.catalog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.title-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: -0.02em;
}

.badge-total {
  background: #eff6ff;
  color: #2563eb;
  padding: 4px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  border: 1px solid #bfdbfe;
}

/* Фильтры */
:deep(.filters-wrapper) {
  width: 100%;
}
:deep(.filters-wrapper .el-select),
:deep(.filters-wrapper .el-input) {
  min-width: 180px;
}

/* Сетка карточек */
.events-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
  margin-top: 24px;
}

/* КРАСИВАЯ КАРТОЧКА */
.event-card-clean {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  padding: 0 !important;
  position: relative;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
}

.event-card-clean:hover {
  transform: translateY(-4px);
  border-color: #93c5fd; /* Синяя рамка при наведении */
  box-shadow: 0 12px 24px -8px rgba(59, 130, 246, 0.15); /* Мягкая синяя тень */
}

.card-header-tags {
  padding: 20px 20px 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.tag-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* Теги с приятными оттенками */
.tag {
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
}
.tag-format {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}
.tag-direction {
  background: #f8fafc;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.status-going {
  display: flex;
  align-items: center;
  gap: 4px;
  background: #dcfce7;
  color: #166534;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  border: 1px solid #bbf7d0;
}

.card-body {
  padding: 16px 20px 20px 20px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.card-title {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}
.info-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}
.info-icon {
  width: 16px;
  height: 16px;
  color: #64748b; /* Мягкий сизо-серый */
  margin-top: 2px;
  flex-shrink: 0;
}
.info-text {
  font-size: 13.5px;
  color: #475569;
  line-height: 1.5;
}

/* Красивые баллы */
.points-wrap {
  margin-top: auto;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  width: fit-content;
}
.points-icon {
  width: 16px;
  height: 16px;
  color: #d97706; /* Золотистый */
}
.points-text {
  font-size: 13px;
  color: #92400e;
}
.points-text strong {
  color: #78350f;
  font-weight: 700;
}

/* Подвал со светло-серым фоном */
.card-footer {
  padding: 16px 20px;
  background: #f8fafc;
  border-top: 1px solid #f1f5f9;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.org-block {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-width: 75%;
}
.org-name {
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.org-rating {
  font-size: 12px;
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 4px;
}
.rating-icon {
  width: 12px;
  height: 12px;
  color: #94a3b8;
}
.org-rating span {
  font-weight: 700;
  color: #334155;
}

.users-block {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #64748b;
  font-size: 13px;
  font-weight: 600;
}

/* Пустое состояние */
.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 20px;
  background: #ffffff;
  border-radius: 16px;
  border: 1px dashed #cbd5e1;
}
.empty-state-icon {
  width: 48px;
  height: 48px;
  color: #94a3b8;
  margin-bottom: 16px;
}
.empty-state h3 { margin: 0 0 8px 0; font-size: 18px; color: #0f172a; font-weight: 600; }
.empty-state p { margin: 0; color: #64748b; font-size: 14px; }

.pagination-center {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}
</style>