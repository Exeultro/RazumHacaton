<template>
  <div class="inspector-view">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ isAdmin ? 'Управление пользователями' : 'Инспектор кадрового резерва' }}</h1>
        <p class="page-description">
          {{ isAdmin ? 'Полный список всех пользователей платформы' : 'Поиск и отбор кандидатов для стажировок и должностей' }}
        </p>
      </div>
      <div class="stats-badge">
        {{ displayTotal }} {{ isAdmin ? 'пользователей' : 'кандидатов' }}
      </div>
    </div>

    <!-- Фильтры показываем только для HR, так как админское API их не поддерживает -->
    <AdvancedFilters
        v-if="!isAdmin"
        :filters="hrStore.filters"
        :savedFilters="hrStore.savedFilters"
        :cities="hrStore.cities"
        @update:filters="hrStore.updateFilters"
        @reset="hrStore.resetFilters"
        @save="hrStore.saveFilter"
        @apply-saved="hrStore.applySavedFilter"
        @delete-filter="handleDeleteFilter"
    />

    <div v-if="isLoading" class="loading-container">
      <div class="spinner"></div>
      <p>Загрузка данных...</p>
    </div>

    <div v-else>
      <div class="results-header">
        <div class="results-count">
          <span class="count-number">{{ displayTotal }}</span>
          <span class="count-label">найдено {{ isAdmin ? 'человек' : 'кандидатов' }}</span>
        </div>
        <ExportPdfButton :candidates="displayCandidates" />
      </div>

      <!-- Таблица кандидатов/пользователей -->
      <CandidateTable :candidates="displayCandidates" :is-admin="isAdmin" />

      <!-- Пагинация для админа -->
      <div v-if="isAdmin && adminStore.usersPagination.total > 0" class="pagination-container">
        <el-pagination
            background
            layout="prev, pager, next"
            :total="adminStore.usersPagination.total"
            :page-size="adminStore.usersPagination.limit"
            v-model:current-page="currentAdminPage"
            @current-change="loadAdminUsers"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed, ref } from 'vue'
import { useHrStore } from '@/stores/hrStore'
import { useAdminStore } from '@/stores/adminStore'
import { useAuthStore } from '@/stores/auth'

import AdvancedFilters from '@/components/hr/AdvancedFilters.vue'
import CandidateTable from '@/components/hr/CandidateTable.vue'
import ExportPdfButton from '@/components/hr/ExportPdfButton.vue'

const hrStore = useHrStore()
const adminStore = useAdminStore()
const authStore = useAuthStore()

// Определяем, админ ли сейчас на странице
const isAdmin = computed(() => authStore.user?.role === 'admin')

// Пагинация для админского списка
const currentAdminPage = ref(1)

// Единые вычисляемые свойства, которые отдают данные в зависимости от роли
const displayCandidates = computed(() => {
  return isAdmin.value ? adminStore.users : hrStore.candidates
})

const displayTotal = computed(() => {
  return isAdmin.value ? adminStore.usersPagination?.total || 0 : hrStore.total || 0
})

const isLoading = computed(() => {
  return isAdmin.value ? adminStore.loading : hrStore.isLoading
})

const loadAdminUsers = () => {
  adminStore.fetchUsers(currentAdminPage.value)
}
const handleDeleteFilter = (filterId) => {
  if (hrStore.deleteFilter) {
    hrStore.deleteFilter(filterId)
        .then(() => console.log('Удаление успешно'))
        .catch(err => console.error('Ошибка при удалении:', err));
  } else {
    console.error('hrStore.deleteFilter не определён');
  }
};

onMounted(() => {
  // Ждем пока загрузится пользователь, если мы зашли по прямой ссылке
  if (!authStore.user) return;

  if (isAdmin.value) {
    // Если Админ — грузим всех юзеров через админское API
    loadAdminUsers()
  } else {
    // Если HR — грузим через HR API
    hrStore.fetchCandidates()
    hrStore.fetchSavedFilters()
  }
})
</script>

<style scoped>
.inspector-view {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
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

.stats-badge {
  background: linear-gradient(135deg, #1d4ed8 0%, #7c3aed 100%);
  color: white;
  padding: 8px 20px;
  border-radius: 40px;
  font-weight: 600;
  font-size: 14px;
  white-space: nowrap;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 24px 0 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e2e8f0;
  flex-wrap: wrap;
  gap: 16px;
}

.results-count {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.count-number {
  font-size: 24px;
  font-weight: 700;
  color: #0f172a;
}

.count-label {
  color: #64748b;
  font-size: 14px;
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

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Адаптивность */
@media (max-width: 768px) {
  .inspector-view { padding: 16px; }
  .page-title { font-size: 24px; }
  .page-description { font-size: 14px; }
  .stats-badge { padding: 6px 16px; font-size: 12px; }
  .results-header { flex-direction: column; align-items: stretch; }
  .results-count { justify-content: space-between; }
}

@media (max-width: 480px) {
  .inspector-view { padding: 12px; }
  .page-header { flex-direction: column; }
  .stats-badge { align-self: flex-start; }
  .results-count .count-number { font-size: 20px; }
}
</style>