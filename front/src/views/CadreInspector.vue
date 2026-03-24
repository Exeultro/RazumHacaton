<template>
    <div class="cadre">
      <h1>Кадровый резерв</h1>
      <el-form :inline="true" @submit.prevent="applyFilters">
        <el-form-item label="Возраст от">
          <el-input-number v-model="filters.age_min" :min="14" :max="100" />
        </el-form-item>
        <el-form-item label="до">
          <el-input-number v-model="filters.age_max" :min="14" :max="100" />
        </el-form-item>
        <el-form-item label="Город">
          <el-input v-model="filters.city" />
        </el-form-item>
        <el-form-item label="Направление">
          <el-select v-model="filters.direction" clearable>
            <el-option label="IT" value="IT" />
            <el-option label="Социальное" value="social" />
            <el-option label="Медиа" value="media" />
          </el-select>
        </el-form-item>
        <el-form-item label="Мин. баллов">
          <el-input-number v-model="filters.min_points" :min="0" />
        </el-form-item>
        <el-form-item label="Мин. мероприятий">
          <el-input-number v-model="filters.min_events" :min="0" />
        </el-form-item>
        <el-form-item label="Сортировка">
          <el-select v-model="filters.sort_by" style="width: 120px">
            <el-option label="Баллы" value="points" />
            <el-option label="Мероприятия" value="events" />
            <el-option label="Возраст" value="age" />
          </el-select>
          <el-select v-model="filters.sort_order" style="width: 100px">
            <el-option label="По убыванию" value="desc" />
            <el-option label="По возрастанию" value="asc" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit">Поиск</el-button>
          <el-button @click="saveFilterDialog = true">Сохранить фильтр</el-button>
        </el-form-item>
      </el-form>
  
      <el-table :data="cadreStore.candidates" v-loading="cadreStore.loading" stripe>
        <el-table-column prop="full_name" label="ФИО" />
        <el-table-column prop="age" label="Возраст" />
        <el-table-column prop="city" label="Город" />
        <el-table-column prop="direction" label="Направление" />
        <el-table-column prop="total_points" label="Баллы" sortable />
        <el-table-column prop="events_count" label="Мероприятий" />
        <el-table-column label="Действия">
          <template #default="{ row }">
            <el-button @click="downloadReport(row.id)" size="small">PDF</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        :current-page="cadreStore.filters.page"
        :page-size="cadreStore.filters.limit"
        :total="cadreStore.total"
        @current-change="handlePageChange"
      />
  
      <el-dialog v-model="saveFilterDialog" title="Сохранить фильтр">
        <el-input v-model="filterName" placeholder="Название фильтра" />
        <template #footer>
          <el-button @click="saveFilterDialog = false">Отмена</el-button>
          <el-button type="primary" @click="saveFilter">Сохранить</el-button>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup>
  import { reactive, ref, onMounted } from 'vue'
  import { useCadreStore } from '../stores/cadreStore'
  import { ElMessage } from 'element-plus'
  
  const cadreStore = useCadreStore()
  const filters = reactive({ ...cadreStore.filters })
  const saveFilterDialog = ref(false)
  const filterName = ref('')
  
  onMounted(() => {
    cadreStore.fetchCandidates()
    cadreStore.loadFilters()
  })
  
  const applyFilters = () => {
    cadreStore.updateFilters(filters)
  }
  
  const handlePageChange = (page) => {
    cadreStore.updateFilters({ page })
  }
  
  const downloadReport = (userId) => {
    cadreStore.downloadReport(userId)
  }
  
  const saveFilter = async () => {
    if (!filterName.value) {
      ElMessage.warning('Введите название фильтра')
      return
    }
    await cadreStore.saveFilter(filterName.value, filters)
    saveFilterDialog.value = false
    filterName.value = ''
  }
  </script>