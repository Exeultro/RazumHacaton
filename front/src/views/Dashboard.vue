<template>
  <div class="dashboard">
    <!-- Заголовок загрузится, когда появится user -->
    <h1 class="page-title">Привет, {{ user?.full_name || '' }}!</h1>

    <!-- ================= УЧАСТНИК ================= -->
    <template v-if="isParticipant">
      <div class="stats-row">
        <div class="card stat-card">
          <div class="icon-circle gold">🌟</div>
          <div class="stat-info">
            <label>Всего баллов</label>
            <div class="value">{{ user?.rating?.total_points || 0 }}</div>
            <p class="sub">Ваш текущий баланс</p>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle blue">🏆</div>
          <div class="stat-info">
            <label>Место в рейтинге</label>
            <div class="value">#{{ user?.rating?.global_rank || '—' }}</div>
            <p class="sub">{{ user?.direction || 'Все' }} направления</p>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle green">📅</div>
          <div class="stat-info">
            <label>Мероприятий</label>
            <div class="value">{{ user?.rating?.events_count || "-" }}</div>
            <p class="sub">Посещено всего</p>
          </div>
        </div>
      </div>

      <div class="main-grid participant-grid">
        <!-- Блок прогресса -->
        <div class="card">
          <h3 class="section-title">Твой путь к цели</h3>
          <div class="progress-container">
            <div class="progress-labels">
              <span>{{ progress }}% до уровня</span>
              <span class="text-muted">{{ user?.rating?.total_points || 0 }} / 1500</span>
            </div>
            <div class="progress-bar-bg">
              <div class="progress-bar-fill" :style="{ width: progress + '%' }"></div>
            </div>
            <p class="hint">Осталось <b>{{ 1500 - (user?.rating?.total_points || 0) }} баллов</b> до награды</p>
          </div>
          <div class="actions-vertical">
            <button class="btn-primary w-100" @click="$router.push('/events')">Найти мероприятия</button>
            <button class="btn-primary w-100" @click="$router.push('/record')">Записаться</button>
          </div>
        </div>

        <!-- Достижения -->
        <div class="card">
          <h3 class="section-title">Последние достижения</h3>
          <table class="table" v-if="user?.portfolio?.events?.length">
            <thead>
            <tr>
              <th>Событие</th>
              <th>Дата</th>
              <th style="text-align: right">Баллы</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="event in user.portfolio.events" :key="event.id">
              <td><b>{{ event.title }}</b></td>
              <td>{{ formatDate(event.date) }}</td>
              <td class="points-earned">+{{ event.points_earned }}</td>
            </tr>
            </tbody>
          </table>
          <p v-else class="empty-state">Пока нет достижений</p>
        </div>

        <!-- Мои мероприятия -->
        <div class="card">
          <h3 class="section-title">Мои мероприятия</h3>
          <div v-if="participations.length" class="participations-list">
            <div v-for="part in participations" :key="part.id" class="participation-item">
              <div class="part-header">
                <div class="part-title">{{ part.event.title }}</div>
                <div class="part-status" :class="part.status">
                  {{ statusLabel(part.status) }}
                </div>
              </div>
              <div class="part-details">
                <span>{{ formatDate(part.event.event_date) }}</span>
                <span>{{ formatDirection(part.event.direction) }}</span>
                <span v-if="part.status === 'attended'">⭐️ +{{ part.points_earned }} баллов</span>
              </div>
              <div class="part-actions" v-if="part.status === 'registered'">
                <button class="btn-small outline" @click="showQrCode(part.event.id, part.event.title)">Показать QR</button>
                <button class="btn-small outline danger" @click="cancelRegistration(part.event.id)">Отменить</button>
              </div>
              <div v-if="part.status === 'cancelled'" class="part-cancelled-message">
                Регистрация отменена
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>Вы ещё не зарегистрированы ни на одно мероприятие</p>
            <button class="btn-primary w-100 mt-3" @click="$router.push('/events')">Найти мероприятия</button>
          </div>
        </div>
      </div>
    </template>

    <!-- ================= ОРГАНИЗАТОР ================= -->
    <template v-else-if="isOrganizer">
      <div class="stats-row">
        <div class="card stat-card">
          <div class="icon-circle purple">🎯</div>
          <div class="stat-info">
            <label>Проведено мероприятий</label>
            <div class="value">{{ organizerStats?.events_count || 0 }}</div>
            <p class="sub">Всего создано</p>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle orange">⭐️</div>
          <div class="stat-info">
            <label>Рейтинг доверия</label>
            <div class="value">{{ organizerStats?.rating || 0 }} / 5</div>
            <p class="sub">Средняя оценка участников</p>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle blue">📋</div>
          <div class="stat-info">
            <label>Ближайшие мероприятия</label>
            <div class="value">{{ organizerEvents.length }}</div>
            <p class="sub">В ближайшие 30 дней</p>
          </div>
        </div>
      </div>

      <div class="main-grid">
        <div class="card">
          <h3 class="section-title">Мои мероприятия</h3>
          <div v-if="organizerEvents.length" class="events-list">
            <div v-for="event in organizerEvents.slice(0, 5)" :key="event.id" class="event-item">
              <div class="event-row">
                <div class="event-title">{{ event.title }}</div>
                <div class="event-date">{{ formatDate(event.event_date) }}</div>
                <div class="event-status" :class="event.status">
                  {{ event.status === 'published' ? 'Активно' : 'Завершено' }}
                </div>
              </div>
              <button class="btn-small event-action-btn" @click="goToConfirm(event.id)">Проверить QR-коды</button>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>У вас пока нет активных мероприятий</p>
          </div>
        </div>

        <div class="card">
          <h3 class="section-title">Быстрые действия</h3>
          <div class="quick-actions">
            <router-link to="/events" class="action-link">Каталог мероприятий</router-link>
            <router-link to="/profile" class="action-link">Мои данные</router-link>
          </div>
        </div>
      </div>
    </template>

    <!-- ================= АДМИНИСТРАТОР ================= -->
    <template v-else-if="isAdmin">
      <div class="stats-row">
        <div class="card stat-card">
          <div class="icon-circle purple">👥</div>
          <div class="stat-info">
            <label>Всего пользователей</label>
            <div class="value">{{ adminStats?.total_users || 0 }}</div>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle orange">🎯</div>
          <div class="stat-info">
            <label>Организаторов</label>
            <div class="value">{{ adminStats?.total_organizers || 0 }}</div>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle green">📅</div>
          <div class="stat-info">
            <label>Мероприятий</label>
            <div class="value">{{ adminStats?.total_events || 0 }}</div>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle blue">🏆</div>
          <div class="stat-info">
            <label>Участий</label>
            <div class="value">{{ adminStats?.total_participations || 0 }}</div>
          </div>
        </div>
      </div>

      <div class="main-grid">
        <!-- Модерация -->
        <div class="card">
          <h3 class="section-title">Ожидают модерации</h3>
          <div v-if="pendingOrganizers.length" class="pending-list">
            <div v-for="org in pendingOrganizers.slice(0, 5)" :key="org.user_id" class="pending-item">
              <div class="org-info">
                <strong>{{ org.full_name }}</strong> – {{ org.email }}
              </div>
              <div class="org-actions">
                <button class="btn-small success" @click="approveOrganizer(org.user_id)">Одобрить</button>
                <button class="btn-small danger" @click="rejectOrganizer(org.user_id)">Отклонить</button>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>Нет заявок на модерацию</p>
          </div>
          <button class="btn-primary w-100 mt-3" @click="$router.push('/admin')">Настройки сайта</button>
        </div>

        <!-- Статистика и создание мероприятий -->
        <div class="card">
          <h3 class="section-title">Управление мероприятиями</h3>
          <div v-if="eventsByDirection" class="directions mb-3">
            <div v-for="(count, dir) in eventsByDirection" :key="dir" class="direction-item">
              <span class="dir-name">{{ formatDirection(dir) }}</span>
              <span class="dir-count">{{ count }}</span>
            </div>
          </div>
          <div class="actions-vertical mt-3">
            <!-- Кнопка только для АДМИНА -->
            <button class="btn-primary w-100" @click="$router.push('/events/create')">Создать мероприятие</button>
            <button class="btn-secondary w-100" @click="refreshStats">Обновить статистику</button>
          </div>
        </div>
      </div>
    </template>

    <template v-else-if="isObserver">
      <div class="stats-row">
        <div class="card stat-card">
          <div class="icon-circle purple">👥</div>
          <div class="stat-info">
            <label>Всего пользователей</label>
            <div class="value">{{ stats.total_users || 0 }}</div>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle blue">📅</div>
          <div class="stat-info">
            <label>Мероприятий</label>
            <div class="value">{{ stats.total_events || 0 }}</div>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle green">🎯</div>
          <div class="stat-info">
            <label>Участий</label>
            <div class="value">{{ stats.total_participations || 0 }}</div>
          </div>
        </div>

        <div class="card stat-card">
          <div class="icon-circle orange">⭐️</div>
          <div class="stat-info">
            <label>Средний балл</label>
            <div class="value">{{ stats.average_points || 0 }}</div>
          </div>
        </div>
      </div>

      <div class="main-grid">
        <div class="card">
          <h3 class="section-title">Активность платформы</h3>
          <div class="chart-placeholder">
            <div v-if="ratingHistory.length" class="mini-chart">
              <div v-for="(point, i) in ratingHistory.slice(-7)" :key="i" class="chart-bar" :style="{ height: (point.points / 100) + 'px' }"></div>
            </div>
            <p v-else class="text-muted">Нет данных для отображения</p>
          </div>
          <p class="chart-caption">Динамика баллов за последние 7 дней</p>
        </div>

        <div class="card">
          <h3 class="section-title">Популярные теги</h3>
          <div class="tags-cloud">
        <span v-for="tag in trendingTags" :key="tag.tag" class="tag-item" :style="{ fontSize: (12 + tag.count / 10) + 'px' }">
          {{ tag.tag }}
        </span>
          </div>
          <div class="actions-vertical mt-3">
            <router-link to="/hr/inspector" class="btn-primary w-100 text-center">Кадровый резерв</router-link>
          </div>
        </div>

        <div class="card full-width">
          <h3 class="section-title">Последние мероприятия</h3>
          <div v-if="recentEvents.length" class="events-list">
            <div v-for="event in recentEvents.slice(0, 5)" :key="event.id" class="event-item">
              <div class="event-row">
                <div class="event-title">{{ event.title }}</div>
                <div class="event-date">{{ formatDate(event.event_date) }}</div>
                <div class="event-status" :class="event.status">
                  {{ event.status === 'published' ? 'Активно' : 'Завершено' }}
                </div>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">Нет мероприятий</div>
        </div>
      </div>
    </template>

    <!-- Диалог для QR-кода (Участник) -->
    <el-dialog
        v-model="qrDialogVisible"
        :title="'QR-код: ' + currentEventTitle"
        width="400px"
        center
    >
      <div v-loading="qrLoading" style="text-align: center;">
        <img v-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="QR Code" style="max-width: 100%;" />
        <p v-else-if="!qrLoading">Не удалось загрузить QR-код</p>
      </div>
      <template #footer>
        <el-button @click="qrDialogVisible = false">Закрыть</el-button>
        <el-button type="primary" @click="downloadQR" :disabled="!qrCodeDataUrl">Скачать QR</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useAdminStore } from '@/stores/adminStore';
import api from '@/api/axios';
import { ElMessage, ElMessageBox } from 'element-plus';
import QRCode from 'qrcode';

import { useDashboardStore } from '@/stores/dashboardStore';

const dashboardStore = useDashboardStore();


const router = useRouter();
const authStore = useAuthStore();
const adminStore = useAdminStore();

// Вычисляемые свойства пользователя
const user = computed(() => authStore.user);
console.log(user)
const isParticipant = computed(() => user.value?.role === 'participant');
const isOrganizer = computed(() => user.value?.role === 'organizer');
const isAdmin = computed(() => user.value?.role === 'admin');

// Состояние (Участник)
const participations = ref([]);
const progress = computed(() => Math.min(Math.round(((user.value?.rating?.total_points || 0) / 1500) * 100), 100));

// Состояние (Организатор)
const organizerStats = computed(() => user.value?.organizer_stats || {});
const organizerEvents = ref([]);

// Состояние (Администратор)
const adminStats = computed(() => adminStore.stats);
const pendingOrganizers = computed(() => adminStore.pendingOrganizers);
const eventsByDirection = computed(() => adminStore.stats?.events_by_direction || null);


const isObserver = computed(() => user.value?.role === 'observer');
const stats = computed(() => adminStore.stats || {});
const recentEvents = computed(() => dashboardStore.recentEvents);
const ratingHistory = computed(() => dashboardStore.ratingHistory);
const trendingTags = computed(() => dashboardStore.trendingTags);

// Модалка QR
const qrDialogVisible = ref(false);
const qrLoading = ref(false);
const qrCodeDataUrl = ref(null);
const currentEventTitle = ref('');

// ================= ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ =================
const formatDate = (iso) => {
  if (!iso) return '';
  return new Date(iso).toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
};

const formatDirection = (dir) => {
  const map = { IT: 'IT', social: 'Социальное', media: 'Медиа' };
  return map[dir] || dir;
};

const statusLabel = (status) => {
  const map = { registered: 'Зарегистрирован', attended: 'Посещено', cancelled: 'Отменено' };
  return map[status] || status;
};

// ================= ЛОГИКА УЧАСТНИКА =================
const loadMyParticipations = async () => {
  try {
    const response = await api.get('/my/participations');
    participations.value = response.data.participations || [];
  } catch (err) {
    console.error(err);
    ElMessage.error('Не удалось загрузить список мероприятий');
  }
};

const cancelRegistration = async (eventId) => {
  try {
    await ElMessageBox.confirm('Отменить регистрацию?', 'Внимание', {
      confirmButtonText: 'Да', cancelButtonText: 'Нет', type: 'warning',
    });
    const response = await api.delete(`/events/${eventId}/cancel`);
    if (response.data.success) {
      ElMessage.success('Регистрация отменена');
      await loadMyParticipations();
    }
  } catch (err) {
    if (err !== 'cancel') ElMessage.error('Ошибка отмены');
  }
};

const showQrCode = async (eventId, eventTitle) => {
  currentEventTitle.value = eventTitle;
  qrDialogVisible.value = true;
  qrLoading.value = true;
  qrCodeDataUrl.value = null;

  try {
    const response = await api.get(`/events/${eventId}/my-qr`);
    qrCodeDataUrl.value = await QRCode.toDataURL(response.data.qr_code_token);
  } catch (err) {
    ElMessage.error('Не удалось получить QR-код');
    qrDialogVisible.value = false;
  } finally {
    qrLoading.value = false;
  }
};

const downloadQR = () => {
  if (!qrCodeDataUrl.value) return;
  const link = document.createElement('a');
  link.download = `qr_${currentEventTitle.value.replace(/\s/g, '_')}.png`;
  link.href = qrCodeDataUrl.value;
  link.click();
};

// ================= ЛОГИКА ОРГАНИЗАТОРА =================
const loadOrganizerEvents = async () => {
  try {
    const response = await api.get('/events', {
      params: { organizer_id: user.value.id, limit: 5, status: 'published' }
    });
    organizerEvents.value = response.data.events || [];
  } catch (err) {
    console.error(err);
  }
};

const goToConfirm = (eventId) => router.push(`/record?eventId=${eventId}`);

// ================= ЛОГИКА АДМИНИСТРАТОРА =================
const approveOrganizer = async (userId) => {
  try {
    await adminStore.approveOrganizer(userId);
    ElMessage.success('Организатор одобрен');
    await adminStore.fetchPendingOrganizers();
  } catch (err) {
    ElMessage.error('Ошибка при одобрении');
  }
};

    const rejectOrganizer = async (userId) => {
      try {
        await ElMessageBox.confirm('Точно отклонить заявку?', 'Подтверждение');
        await adminStore.rejectOrganizer(userId);
        ElMessage.success('Заявка отклонена');
        await adminStore.fetchPendingOrganizers();
      } catch (err) {
        if (err !== 'cancel') ElMessage.error('Ошибка при отклонении');
      }
    };

const refreshStats = async () => {
  await adminStore.fetchStats();
  ElMessage.success('Статистика обновлена');
};

watch(
    () => authStore.user,
    async (newUser) => {
      if (!newUser) return;

      if (newUser.role === 'participant') {
        await loadMyParticipations();
      } else if (newUser.role === 'organizer') {
        await loadOrganizerEvents();
      } else if (newUser.role === 'admin') {
        await adminStore.fetchStats();
        await adminStore.fetchPendingOrganizers();
      } else if (newUser.role === 'observer') {
        await adminStore.fetchStats();
        await dashboardStore.fetchActivity();
      }
    },
    { immediate: true }
);
</script>

<style scoped>
.page-title {
  margin-bottom: 32px;
  font-size: 2rem;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
}

.icon-circle {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
}

.icon-circle.gold { background: #fffbeb; border: 1px solid #fcd34d; }
.icon-circle.blue { background: #eff6ff; border: 1px solid #93c5fd; }
.icon-circle.green { background: #f0fdf4; border: 1px solid #86efac; }
.icon-circle.purple { background: #f5f0ff; border: 1px solid #c4b5fd; }
.icon-circle.orange { background: #fff7ed; border: 1px solid #fed7aa; }

.stat-info label {
  color: var(--text-muted);
  font-size: 0.8rem;
  text-transform: uppercase;
  font-weight: 700;
}

.stat-info .value {
  font-size: 1.8rem;
  font-weight: 800;
}

.stat-info .sub {
  margin: 0;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.main-grid {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 24px;
}

.empty-state {
  text-align: center;
  padding: 32px 0;
  color: var(--text-muted);
}

/* Кнопки */
.btn-primary, .btn-secondary, .btn-small {
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
  border: none;
}

.btn-primary { background: var(--primary); color: white; padding: 12px; }
.btn-primary:hover { background: var(--primary-hover); }

.btn-secondary { background: #f1f5f9; color: var(--text-main); padding: 12px; border: 1px solid #e2e8f0; }
.btn-secondary:hover { background: #e2e8f0; }

.btn-small { padding: 4px 12px; font-size: 0.8rem; background: #e2e8f0; color: var(--text-main); }
.btn-small:hover { background: #cbd5e1; }
.btn-small.success { background: #10b981; color: white; }
.btn-small.danger { background: #ef4444; color: white; }

.w-100 { width: 100%; }
.mt-3 { margin-top: 16px; }
.mb-3 { margin-bottom: 16px; }
.actions-vertical { display: flex; flex-direction: column; gap: 8px; }

/* Таблицы */
.table { width: 100%; border-collapse: collapse; }
.table th, .table td { padding: 12px 8px; border-bottom: 1px solid var(--border-color); text-align: left; }
.table th { font-weight: 600; color: var(--text-muted); }
.points-earned { text-align: right; color: #10b981; font-weight: bold; }

/* ================= УЧАСТНИК ================= */
.progress-container { margin: 24px 0; }
.progress-labels { display: flex; justify-content: space-between; margin-bottom: 10px; font-weight: 600; }
.progress-bar-bg { background: #e2e8f0; height: 12px; border-radius: 6px; overflow: hidden; }
.progress-bar-fill { height: 100%; background: var(--primary); transition: 1s ease-out; }
.hint { margin-top: 12px; font-size: 0.9rem; }

                   .participation-item { padding: 12px 0; border-bottom: 1px solid var(--border-color); }
.part-header { display: flex; justify-content: space-between; font-weight: 500; margin-bottom: 4px;}
.part-details { font-size: 0.85rem; color: var(--text-muted); display: flex; gap: 12px; margin-bottom: 8px; }
.part-status { font-size: 0.75rem; padding: 2px 8px; border-radius: 12px; background: #e2e8f0; }
.part-status.attended { background: #dcfce7; color: #15803d; }
.part-status.cancelled { background: #fee2e2; color: #b91c1c; }
.part-actions { display: flex; gap: 8px; }

/* ================= ОРГАНИЗАТОР ================= */
.events-list { margin: 16px 0; }
.event-item { margin-bottom: 16px; padding-bottom: 12px; border-bottom: 1px solid var(--border-color); }
.event-row { display: flex; justify-content: space-between; align-items: center; gap: 12px; margin-bottom: 12px; }
.event-title { font-weight: 500; flex: 2; }
.event-date { flex: 1; color: var(--text-muted); font-size: 0.85rem; }
.event-status { font-size: 0.75rem; padding: 2px 8px; border-radius: 12px; background: #e2e8f0; }
.event-status.published { background: #dcfce7; color: #15803d; }
.event-action-btn { width: 100%; padding: 8px; background: var(--primary); color: white; border-radius: 6px; }
.event-action-btn:hover { background: var(--primary-hover); }

.quick-actions { display: flex; flex-direction: column; gap: 12px; margin-top: 16px; }
.action-link { padding: 10px 0; color: var(--primary); text-decoration: none; font-weight: 500; border-bottom: 1px solid var(--border-color); }

/* ================= АДМИНИСТРАТОР ================= */
.pending-list { margin: 16px 0; }
.pending-item { display: flex; justify-content: space-between; align-items: center; padding: 12px 0; border-bottom: 1px solid var(--border-color); }
.org-info { flex: 2; font-size: 0.9rem; }
.org-actions { display: flex; gap: 8px; }

.directions { margin: 16px 0; }
.direction-item { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid var(--border-color); }
.dir-name { font-weight: 500; }
.dir-count { font-weight: 600; color: var(--primary); }

@media (max-width: 900px) {
  .main-grid { grid-template-columns: 1fr; }
}

/* Для графика */
.chart-placeholder {
  min-height: 120px;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 8px;
  margin: 16px 0;
}
.mini-chart {
  display: flex;
  gap: 8px;
  width: 100%;
  justify-content: center;
}
.chart-bar {
  width: 30px;
  background: var(--primary);
  border-radius: 4px 4px 0 0;
  transition: height 0.3s ease;
}
.chart-caption {
  font-size: 0.75rem;
  color: var(--text-muted);
  text-align: center;
}

/* Облако тегов */
.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  margin: 20px 0;
}
.tag-item {
  display: inline-block;
  padding: 4px 12px;
  background: #f1f5f9;
  border-radius: 20px;
  color: var(--text-main);
  cursor: default;
  transition: all 0.2s;
}
.tag-item:hover {
  background: #e2e8f0;
  transform: scale(1.02);
}

.full-width {
  grid-column: span 2;
}
</style>