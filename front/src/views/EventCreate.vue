<template>
  <div class="content-area">
    <div class="header-actions">
      <h1 class="page-title">Создание нового мероприятия</h1>
      <button class="btn-outline" @click="$router.push('/dashboard')">
        <i class="el-icon-back"></i> Вернуться назад
      </button>
    </div>

    <div class="card form-card">
      <el-form
          :model="form"
          label-position="top"
          class="custom-form"
      >

        <!-- ================= БЛОК 1: Основная информация ================= -->
        <h3 class="section-title">
          <span class="section-icon">📌</span> Основная информация
        </h3>

        <el-row :gutter="20">
          <el-col :xs="24" :md="16">
            <el-form-item label="Название мероприятия" required>
              <el-input
                  v-model="form.title"
                  placeholder="Например: Хакатон по Go 2026"
                  size="large"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="Направление" required>
              <el-select v-model="form.direction" placeholder="Выберите" class="w-100" size="large">
                <el-option label="IT" value="IT" />
                <el-option label="Социальное проектирование" value="social" />
                <el-option label="Медиа" value="media" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Описание мероприятия">
          <el-input
              type="textarea"
              v-model="form.description"
              :rows="4"
              placeholder="Опишите цель, программу и для кого это мероприятие..."
          />
        </el-form-item>

        <!-- ВЫБОР ОРГАНИЗАТОРА (Только для Админа) -->
        <el-form-item label="Назначить организатора" required>
          <el-select
              v-model="form.organizer_id"
              placeholder="Выберите организатора из списка"
              class="w-100"
              size="large"
              filterable
              :loading="loadingOrganizers"
          >
            <el-option
                v-for="org in organizers"
                :key="org.id"
                :label="`${org.full_name} (${org.email})`"
                :value="org.id"
            />
            <template #empty>
              <p style="text-align: center; color: #999; padding: 10px;">
                Организаторы не найдены. Проверьте вкладку "Пользователи".
              </p>
            </template>
          </el-select>
        </el-form-item>

        <el-divider />

        <!-- ================= БЛОК 2: Даты и Локация ================= -->
        <h3 class="section-title">
          <span class="section-icon">📅</span> Время и место
        </h3>

        <el-row :gutter="20">
          <el-col :xs="24" :md="8">
            <el-form-item label="Дата проведения" required>
              <el-date-picker
                  v-model="form.event_date"
                  type="datetime"
                  placeholder="Начало мероприятия"
                  class="w-100"
                  format="DD.MM.YYYY HH:mm"
                  value-format="YYYY-MM-DDTHH:mm:ss[Z]"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="Дедлайн регистрации">
              <el-date-picker
                  v-model="form.registration_deadline"
                  type="datetime"
                  placeholder="Конец приема заявок"
                  class="w-100"
                  format="DD.MM.YYYY HH:mm"
                  value-format="YYYY-MM-DDTHH:mm:ss[Z]"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="Формат участия">
              <el-radio-group v-model="form.format" class="w-100 format-group">
                <el-radio-button label="offline">Офлайн</el-radio-button>
                <el-radio-button label="online">Онлайн</el-radio-button>
                <el-radio-button label="hybrid">Гибрид</el-radio-button>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Локация (Адрес или ссылка на трансляцию)">
          <el-input
              v-model="form.location"
              placeholder="Например: г. Москва, ул. Пушкина, д. 1 / Ссылка на Zoom"
          />
        </el-form-item>

        <el-divider />

        <!-- ================= БЛОК 3: Геймификация ================= -->
        <h3 class="section-title">
          <span class="section-icon">🎮</span> Геймификация и призы
        </h3>

        <el-row :gutter="20">
          <el-col :xs="24" :md="12">
            <el-form-item label="Базовые баллы за участие">
              <el-input-number v-model="form.points_for_participation" :min="0" class="w-100" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="Коэффициент сложности (множитель)">
              <el-input-number v-model="form.difficulty_coefficient" :min="0" :step="0.1" class="w-100" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="Призы для участников (необязательно)">
          <div class="prizes-container">
            <div v-for="(prize, idx) in form.prizes" :key="idx" class="prize-row">
              <el-select v-model="prize.type" placeholder="Тип приза" class="prize-type">
                <el-option label="Мерч" value="merch" />
                <el-option label="Стажировка" value="internship" />
                <el-option label="Билеты" value="tickets" />
                <el-option label="Другое" value="other" />
              </el-select>
              <el-input v-model="prize.name" placeholder="Название (например: Худи с логотипом)" class="prize-name" />
              <button class="btn-small danger" @click.prevent="removePrize(idx)">Удалить</button>
            </div>

            <button class="btn-secondary add-prize-btn" @click.prevent="addPrize">
              + Добавить приз
            </button>
          </div>
        </el-form-item>

        <!-- ================= КНОПКИ ================= -->
        <div class="form-actions">
          <button class="btn-primary submit-btn" @click.prevent="submit" :disabled="submitting">
            <span v-if="submitting" class="spinner"></span>
            <span v-else>Опубликовать мероприятие</span>
          </button>
        </div>

      </el-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue';
import { useEventStore } from '@/stores/eventStore';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import api from '@/api/axios';

const eventStore = useEventStore();
const router = useRouter();

const submitting = ref(false);
const loadingOrganizers = ref(false);
const organizers = ref([]);

// Состояние формы
const form = reactive({
  title: '',
  description: '',
  organizer_id: null, // Обязательное поле для админа
  event_date: null,
  registration_deadline: null,
  location: '',
  format: 'offline',
  direction: '',
  points_for_participation: 10, // Дефолтное значение
  difficulty_coefficient: 1.0,
  prizes: []
});

// Загрузка списка пользователей с ролью organizer при открытии страницы
onMounted(async () => {
  loadingOrganizers.value = true;
  try {
    // Используем админский эндпоинт получения пользователей
    // Берем сразу 100 человек, чтобы точно захватить организаторов (если их не тысячи)
    const response = await api.get('/admin/users?page=1&limit=100');

    // Благодаря вашему интерцептору, данные лежат сразу в response.data
    const allUsers = response.data.users || [];

    // Отфильтровываем только организаторов
    organizers.value = allUsers.filter(user => user.role === 'organizer');

  } catch (error) {
    console.error('Ошибка загрузки организаторов:', error);
    ElMessage.error('Не удалось загрузить список организаторов. Проверьте права доступа.');
  } finally {
    loadingOrganizers.value = false;
  }
});

// Работа с призами
const addPrize = () => form.prizes.push({ type: 'merch', name: '' });
const removePrize = (idx) => form.prizes.splice(idx, 1);

// Отправка формы
const submit = async () => {
  // Валидация
  if (!form.title) return ElMessage.warning('Укажите название мероприятия');
  if (!form.direction) return ElMessage.warning('Выберите направление');
  if (!form.organizer_id) return ElMessage.warning('Выберите организатора');
  if (!form.event_date) return ElMessage.warning('Укажите дату проведения');

  submitting.value = true;
  try {
    await eventStore.createEvent(form);
    ElMessage.success('Мероприятие успешно создано!');

    router.push('/dashboard');
  } catch (err) {
    ElMessage.error('Не удалось создать мероприятие. Проверьте введенные данные.');
    console.error(err);
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped>
.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.form-card {
  max-width: 900px;
  margin: 0 auto;
  padding: 32px;
}

.section-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--text-main);
  margin: 24px 0 16px 0;
  display: flex;
  align-items: center;
  gap: 10px;
  border-bottom: 2px solid var(--bg-app);
  padding-bottom: 8px;
}

.section-title:first-child {
  margin-top: 0;
}

.section-icon {
  font-size: 1.4rem;
}

.w-100 {
  width: 100% !important;
}

.format-group {
  display: flex;
}
.format-group :deep(.el-radio-button__inner) {
  width: 100%;
}

/* Призы */
.prizes-container {
  background: var(--bg-app);
  border: 1px dashed var(--border-color);
  padding: 16px;
  border-radius: var(--radius-md);
}

.prize-row {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
  align-items: center;
  background: var(--bg-card);
  padding: 12px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
}

.prize-type {
  width: 180px;
  flex-shrink: 0;
}

.prize-name {
  flex-grow: 1;
}

.add-prize-btn {
  width: 100%;
  border: 1px dashed var(--secondary);
  background: transparent;
  color: var(--secondary);
}
.add-prize-btn:hover {
  background: var(--bg-card);
  border-color: var(--primary);
  color: var(--primary);
}

/* Кнопки формы */
.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--border-color);
}

.submit-btn {
  padding: 14px 32px;
  font-size: 1rem;
}

.btn-small.danger {
  background: #fee2e2;
  color: #ef4444;
  border: none;
}
.btn-small.danger:hover {
  background: #fca5a5;
  color: #b91c1c;
}

/* Адаптивность */
@media (max-width: 768px) {
  .form-card {
    padding: 20px;
  }

  .prize-row {
    flex-direction: column;
    align-items: stretch;
  }

  .prize-type {
    width: 100%;
  }

  .form-actions {
    justify-content: center;
  }

  .submit-btn {
    width: 100%;
  }
}
</style>