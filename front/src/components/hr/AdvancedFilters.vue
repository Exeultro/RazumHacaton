<template>
  <div class="filters-card">
    <div class="filters-header">
      <h3 class="filters-title">Фильтры поиска</h3>
      <div class="filters-actions">
        <el-button size="small" @click="handleReset" :icon="Refresh">Сбросить</el-button>
        <el-button type="primary" size="small" @click="handleSave" :icon="DocumentAdd">Сохранить фильтр</el-button>
      </div>
    </div>

    <el-form :model="localFilters" label-position="top" class="filters-form">
      <el-row :gutter="16">
        <!-- Направление -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="filter-col">
          <el-form-item label="Направление">
            <el-select
                v-model="localFilters.direction"
                placeholder="Любое"
                clearable
                @change="applyFilters"
            >
              <el-option label="IT" value="IT" />
              <el-option label="Социальное" value="social" />
              <el-option label="Медиа" value="media" />
            </el-select>
          </el-form-item>
        </el-col>

        <!-- Город -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="filter-col">
          <el-form-item label="Город">
            <el-select
                v-model="localFilters.city"
                placeholder="Любой"
                clearable
                filterable
                @change="applyFilters"
            >
              <el-option
                  v-for="city in cities"
                  :key="city"
                  :label="city"
                  :value="city"
              />
            </el-select>
          </el-form-item>
        </el-col>

        <!-- Возраст -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="filter-col">
          <el-form-item label="Возраст">
            <div class="age-range">
              <el-input-number
                  v-model="localFilters.age_min"
                  :min="14"
                  :max="100"
                  placeholder="От"
                  controls-position="right"
                  @change="applyFilters"
              />
              <span class="age-separator">—</span>
              <el-input-number
                  v-model="localFilters.age_max"
                  :min="14"
                  :max="100"
                  placeholder="До"
                  controls-position="right"
                  @change="applyFilters"
              />
            </div>
          </el-form-item>
        </el-col>

        <!-- Баллы (мин.) -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="filter-col">
          <el-form-item label="Баллы (мин.)">
            <el-input-number
                v-model="localFilters.min_points"
                :min="0"
                :step="10"
                placeholder="от 0"
                controls-position="right"
                @change="applyFilters"
            />
          </el-form-item>
        </el-col>

        <!-- Мероприятий (мин.) -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="filter-col">
          <el-form-item label="Мероприятий (мин.)">
            <el-input-number
                v-model="localFilters.min_events"
                :min="0"
                :step="1"
                placeholder="от 0"
                controls-position="right"
                @change="applyFilters"
            />
          </el-form-item>
        </el-col>

        <!-- Средний балл (мин.) -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="filter-col">
          <el-form-item label="Ср. балл (мин.)">
            <el-input-number
                v-model="localFilters.min_avg_points"
                :min="0"
                :step="5"
                placeholder="от 0"
                controls-position="right"
                @change="applyFilters"
            />
          </el-form-item>
        </el-col>

        <!-- Сортировка -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" :xl="4" class="filter-col">
          <el-form-item label="Сортировка">
            <div class="sort-group">
              <el-select
                  v-model="localFilters.sort_by"
                  @change="applyFilters"
                  style="flex: 1"
              >
                <el-option label="По баллам" value="points" />
                <el-option label="По мероприятиям" value="events" />
                <el-option label="По возрасту" value="age" />
              </el-select>
              <el-select
                  v-model="localFilters.sort_order"
                  @change="applyFilters"
                  style="width: 100px"
              >
                <el-option label="↓ убыв." value="desc" />
                <el-option label="↑ возр." value="asc" />
              </el-select>
            </div>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <!-- Сохранённые фильтры -->
    <div v-if="savedFilters.length" class="saved-filters">
      <div class="saved-header">
        <span class="saved-title">Сохранённые фильтры</span>
      </div>
      <div class="saved-list">
        <el-tag
            v-for="filter in savedFilters"
            :key="filter.id"
            class="saved-tag"
            closable
            @close="handleDeleteFilter(filter.id)"
            @click="handleApplySavedFilter(filter)"
        >
          {{ filter.name }}
        </el-tag>
      </div>
    </div>
  </div>

  <!-- Диалог сохранения фильтра -->
  <el-dialog v-model="saveDialogVisible" title="Сохранить фильтр" width="400px">
    <el-input
        v-model="filterName"
        placeholder="Название фильтра"
        @keyup.enter="confirmSave"
    />
    <template #footer>
      <el-button @click="saveDialogVisible = false">Отмена</el-button>
      <el-button type="primary" @click="confirmSave">Сохранить</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Refresh, DocumentAdd } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps({
  filters: {
    type: Object,
    required: true
  },
  savedFilters: {
    type: Array,
    default: () => []
  },
  cities: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:filters', 'reset', 'save', 'apply-saved', 'delete-filter'])

// Локальная копия фильтров
const localFilters = ref({ ...props.filters })

const saveDialogVisible = ref(false)
const filterName = ref('')

const applyFilters = () => {
  emit('update:filters', { ...localFilters.value })
}

const handleReset = () => {
  localFilters.value = {
    age_min: null,
    age_max: null,
    city: null,
    direction: null,
    min_points: null,
    min_events: null,
    min_avg_points: null,
    sort_by: 'points',
    sort_order: 'desc',
    page: 1,
    limit: 20
  }
  emit('reset')
}

const handleSave = () => {
  filterName.value = ''
  saveDialogVisible.value = true
}

const confirmSave = () => {
  if (!filterName.value.trim()) {
    ElMessage.warning('Введите название фильтра')
    return
  }
  emit('save', filterName.value.trim())
  saveDialogVisible.value = false
}

const handleApplySavedFilter = (savedFilter) => {
  localFilters.value = {
    ...localFilters.value,
    ...savedFilter.filters,
    page: 1
  }
  emit('apply-saved', savedFilter)
}

const handleDeleteFilter = async (filterId) => {
  console.log('handleDeleteFilter вызван с id:', filterId); // <-- добавить
  try {
    await ElMessageBox.confirm('Удалить этот фильтр?', 'Подтверждение', {
      type: 'warning'
    });
    console.log('Подтверждено, эмитим delete-filter'); // <-- добавить
    emit('delete-filter', filterId);
  } catch {
    console.log('Отменено');
  }
};

watch(() => props.filters, (newFilters) => {
  localFilters.value = { ...newFilters }
}, { deep: true })
</script>

<style scoped>
.filters-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  padding: 20px;
  margin-bottom: 24px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.filters-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.filters-title {
  font-size: 1rem;
  font-weight: 600;
  color: #0f172a;
  margin: 0;
}

.filters-actions {
  display: flex;
  gap: 8px;
}

.filters-form {
  margin-bottom: 16px;
}

/* Возрастные инпуты */
.age-range {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.age-range :deep(.el-input-number) {
  width: 100%;
}

.age-separator {
  color: #94a3b8;
  font-weight: 500;
}

/* Группа сортировки */
.sort-group {
  display: flex;
  gap: 8px;
  width: 100%;
}

.sort-group .el-select:first-child {
  flex: 1;
}

.sort-group .el-select:last-child {
  width: 100px;
}

/* Сохранённые фильтры */
.saved-filters {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.saved-header {
  margin-bottom: 12px;
}

.saved-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: #475569;
}

.saved-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.saved-tag {
  cursor: pointer;
  transition: all 0.2s;
}

.saved-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* Стили для элементов Element Plus внутри scoped */
:deep(.el-form-item__label) {
  font-weight: 500;
  font-size: 0.875rem;
  color: #334155;
  padding-bottom: 6px;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-select) {
  width: 100%;
}

/* Адаптивность */
@media (max-width: 768px) {
  .filters-card {
    padding: 16px;
  }

  .filters-header {
    flex-direction: column;
    align-items: stretch;
  }

  .filters-actions {
    justify-content: stretch;
  }

  .filters-actions .el-button {
    flex: 1;
    justify-content: center;
  }

  .age-range {
    flex-direction: column;
    gap: 8px;
  }

  .age-separator {
    display: none;
  }

  .sort-group {
    flex-direction: column;
  }

  .sort-group .el-select:last-child {
    width: 100% !important;
  }

  .saved-list {
    justify-content: center;
  }

  .saved-tag {
    font-size: 12px;
    padding: 0 8px;
    height: 28px;
    line-height: 26px;
  }
}

@media (max-width: 480px) {
  .filters-card {
    padding: 12px;
  }

  .filter-col {
    margin-bottom: 8px;
  }

  :deep(.el-form-item) {
    margin-bottom: 16px;
  }

  :deep(.el-input-number) {
    width: 100%;
  }
}

.saved-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px; /* увеличил отступ между тегами */
}

.saved-tag {
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;        /* явный размер шрифта */
  padding: 8px 16px;      /* больше внутренних отступов */
  height: auto;           /* снимаем фиксированную высоту */
  border-radius: 20px;
  background-color: #f1f5f9;
  color: #1e293b;
  border: 1px solid #e2e8f0;
}

.saved-tag:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  background-color: #eef2ff;
  border-color: #cbd5e1;
}

@media (max-width: 768px) {
  .saved-tag {
    font-size: 13px;    /* чуть меньше, но всё ещё читаемо */
    padding: 6px 12px;
    height: auto;
  }
}
</style>