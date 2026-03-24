<template>
  <div v-if="eventStore.currentEvent" v-loading="eventStore.loading" class="card event-details-container">

    <!-- Левая колонка: Информация о мероприятии -->
    <div class="event-info">
      <h2 class="page-title">{{ eventStore.currentEvent.title }}</h2>
      <p>{{ eventStore.currentEvent.description }}</p>
      <p><strong>Дата:</strong> {{ formatDate(eventStore.currentEvent.event_date) }}</p>
      <p><strong>Локация:</strong> {{ eventStore.currentEvent.location }}</p>
      <p><strong>Формат:</strong> {{ eventStore.currentEvent.format }}</p>
      <p><strong>Направление:</strong> {{ eventStore.currentEvent.direction }}</p>
      <p><strong>Баллы:</strong> {{ eventStore.currentEvent.points_for_participation }}</p>
      <p><strong>Коэффициент сложности:</strong> {{ eventStore.currentEvent.difficulty_coefficient }}</p>
      <p><strong>Призы:</strong> {{ eventStore.currentEvent.prizes?.map(p => p.name).join(', ') || '—' }}</p>

      <!-- Обновленная ссылка на профиль с query-параметрами -->
      <p><strong>Организатор:</strong>
        <router-link :to="{
          path: `/organizer/${eventStore.currentEvent.organizer.id}`,
          query: { event_id: eventStore.currentEvent.id, event_title: eventStore.currentEvent.title }
        }">
          {{ eventStore.currentEvent.organizer?.full_name }}
        </router-link>
      </p>

      <div class="actions" style="margin-top: 24px; display: flex; gap: 12px; flex-wrap: wrap;">
        <button v-if="!eventStore.currentEvent.is_registered" class="btn btn-primary" @click="register">
          Зарегистрироваться
        </button>
        <button v-else class="btn btn-danger" @click="cancel">
          Отменить регистрацию
        </button>

        <!-- Кнопка для перехода к написанию отзыва -->
        <button class="btn btn-secondary" @click="goToReview">
          Оценить организатора
        </button>
      </div>
    </div>

    <!-- Правая колонка: QR-код -->
    <div class="qr-section" ref="qrWrapper">
      <p class="qr-title">Поделиться событием</p>
      <qrcode-vue
          :value="qrUrl"
          :size="160"
          level="H"
          render-as="canvas"
      />
      <el-button type="primary" plain class="download-btn" @click="downloadQR">
        Скачать QR-код
      </el-button>
    </div>

  </div>
</template>

<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useEventStore } from '../stores/eventStore'
import { ElMessage, ElMessageBox } from 'element-plus'
import QrcodeVue from 'qrcode.vue'

const route = useRoute()
const router = useRouter()
const eventStore = useEventStore()
const eventId = route.params.id

const qrWrapper = ref(null)

const qrUrl = computed(() => {
  return window.location.origin + route.fullPath
})

onMounted(() => {
  eventStore.fetchEventById(eventId)
})

const formatDate = (iso) => new Date(iso).toLocaleString()

const register = async () => {
  try {
    await eventStore.registerForEvent(eventId)
    await eventStore.fetchEventById(eventId)
  } catch (err) {}
}

const cancel = async () => {
  try {
    await ElMessageBox.confirm('Вы уверены, что хотите отменить регистрацию?', 'Подтверждение')
    await eventStore.cancelRegistration(eventId)
    await eventStore.fetchEventById(eventId)
  } catch (err) {
    if (err !== 'cancel') ElMessage.error('Ошибка отмены')
  }
}

// Функция перехода к написанию отзыва
const goToReview = () => {
  router.push({
    path: `/organizer/${eventStore.currentEvent.organizer.id}`,
    query: {
      event_id: eventStore.currentEvent.id,
      event_title: eventStore.currentEvent.title
    }
  })
}

const downloadQR = () => {
  const canvas = qrWrapper.value?.querySelector('canvas')

  if (canvas) {
    const imageUrl = canvas.toDataURL("image/png")
    const link = document.createElement('a')
    link.href = imageUrl
    link.download = `event-${eventId}-qr.png`

    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    ElMessage.success('QR-код успешно скачан!')
  } else {
    ElMessage.error('Не удалось найти QR-код для скачивания')
  }
}
</script>

<style scoped>
.event-details-container { display: flex; flex-direction: column; gap: 30px; }
@media (min-width: 768px) {
  .event-details-container { flex-direction: row; justify-content: space-between; align-items: flex-start; }
}
.event-info { flex: 1; }
.qr-section {
  display: flex; flex-direction: column; align-items: center; gap: 16px; padding: 20px;
  background-color: #f9f9fc; border-radius: 12px; border: 1px solid #ebeef5; min-width: 200px;
}
.qr-title { margin: 0; font-weight: bold; color: #606266; text-align: center; }
.download-btn { width: 100%; }

/* Добавлен стиль для новой кнопки */
.btn-secondary {
  background: #f1f5f9;
  color: #1e293b;
  border: 1px solid #cbd5e1;
}
.btn-secondary:hover {
  background: #e2e8f0;
}
</style>