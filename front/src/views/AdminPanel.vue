<template>
    <div class="admin">
      <h1>Панель администратора</h1>
      <el-tabs>
        <el-tab-pane label="Модерация организаторов">
          <el-table :data="adminStore.pendingOrganizers" v-loading="adminStore.loading">
            <el-table-column prop="full_name" label="Имя" />
            <el-table-column prop="email" label="Email" />
            <el-table-column prop="registered_at" label="Дата регистрации" />
            <el-table-column label="Действия">
              <template #default="{ row }">
                <el-button @click="adminStore.approveOrganizer(row.user_id)" type="success">Одобрить</el-button>
                <el-button @click="adminStore.rejectOrganizer(row.user_id)" type="danger">Отклонить</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="Настройка весов">
          <el-form :model="coefficients" label-width="200px">
            <el-form-item label="IT коэффициент">
              <el-input-number v-model="coefficients.IT" :min="0" :step="0.1" />
            </el-form-item>
            <el-form-item label="Социальное проектирование">
              <el-input-number v-model="coefficients.social" :min="0" :step="0.1" />
            </el-form-item>
            <el-form-item label="Медиа">
              <el-input-number v-model="coefficients.media" :min="0" :step="0.1" />
            </el-form-item>
            <el-form-item>
              <el-button @click="saveCoefficients" type="primary">Сохранить</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="Статистика">
          <el-descriptions :column="2" border v-if="adminStore.stats">
            <el-descriptions-item label="Всего пользователей">{{ adminStore.stats.total_users }}</el-descriptions-item>
            <el-descriptions-item label="Организаторов">{{ adminStore.stats.total_organizers }}</el-descriptions-item>
            <el-descriptions-item label="Мероприятий">{{ adminStore.stats.total_events }}</el-descriptions-item>
            <el-descriptions-item label="Участий">{{ adminStore.stats.total_participations }}</el-descriptions-item>
            <el-descriptions-item label="Мероприятий по направлениям" :span="2">
              <div v-for="(count, dir) in adminStore.stats.events_by_direction" :key="dir">{{ dir }}: {{ count }}</div>
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
      </el-tabs>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useAdminStore } from '../stores/adminStore'
  
  const adminStore = useAdminStore()
  const coefficients = ref({ ...adminStore.coefficients })
  
  onMounted(() => {
    adminStore.fetchPendingOrganizers()
    adminStore.fetchStats()
  })
  
  const saveCoefficients = async () => {
    await adminStore.updateDifficultySettings(adminStore.settings);
  }
  </script>