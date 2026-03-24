<template>
  <div class="confirm-container">
    <el-card class="confirm-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <el-icon class="header-icon"><Ticket /></el-icon>
          <h2>Подтверждение участия</h2>
        </div>
      </template>

      <el-tabs v-model="activeTab" class="custom-tabs" @tab-change="handleTabChange">

        <!-- Вкладка 1: QR-код (Камера / Файл) -->
        <el-tab-pane label="QR-код" name="qr">

          <div class="qr-mode-switch">
            <el-radio-group v-model="qrMode" size="large" @change="switchQrMode">
              <el-radio-button label="camera">
                <el-icon><Camera /></el-icon> Камера
              </el-radio-button>
              <el-radio-button label="file">
                <el-icon><Picture /></el-icon> Картинка
              </el-radio-button>
            </el-radio-group>
          </div>

          <!-- Режим: Камера -->
          <div v-show="qrMode === 'camera'" class="scanner-wrapper">
            <div id="qr-reader" class="qr-reader"></div>
            <p class="hint-text">Наведите камеру на QR-код мероприятия</p>
          </div>

          <!-- Режим: Загрузка файла (с компа или галереи телефона) -->
          <div v-show="qrMode === 'file'" class="upload-wrapper">
            <el-upload
                class="qr-uploader"
                drag
                action="#"
                :auto-upload="false"
                :show-file-list="false"
                accept="image/*"
                :on-change="handleFileUpload"
            >
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">
                Перетащите картинку сюда или <em>нажмите для загрузки</em>
              </div>
              <template #tip>
                <div class="el-upload__tip text-center">
                  Поддерживаются форматы JPG, PNG
                </div>
              </template>
            </el-upload>
          </div>
        </el-tab-pane>

        <!-- Вкладка 2: Ввод вручную -->
        <el-tab-pane label="Ввести вручную" name="manual">
          <div class="manual-wrapper">
            <el-form @submit.prevent="confirmByCode">
              <el-form-item>
                <el-input
                    v-model="manualCode"
                    placeholder="Например: EVT-123456"
                    size="large"
                    clearable
                >
                  <template #prefix>
                    <el-icon><Key /></el-icon>
                  </template>
                </el-input>
              </el-form-item>

              <el-button
                  type="primary"
                  size="large"
                  class="full-width-btn"
                  @click="confirmByCode"
                  :loading="isLoading"
                  :disabled="!manualCode"
              >
                Подтвердить код
              </el-button>
            </el-form>
          </div>
        </el-tab-pane>

      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useEventStore } from '../stores/eventStore'
import { ElMessage } from 'element-plus'
import { Ticket, Camera, Picture, UploadFilled, Key } from '@element-plus/icons-vue'
import { Html5Qrcode } from 'html5-qrcode'

const route = useRoute()
const router = useRouter()
const eventStore = useEventStore()

// Состояния UI
const activeTab = ref('qr')
const qrMode = ref('camera') // 'camera' или 'file'
const manualCode = ref('')
const isLoading = ref(false)

// Инстанс сканера
let html5QrCode = null

// === ЛОГИКА КАМЕРЫ ===
const startCamera = async () => {
  if (html5QrCode?.isScanning) return;

  try {
    html5QrCode = new Html5Qrcode("qr-reader");
    await html5QrCode.start(
        { facingMode: "environment" }, // Использовать заднюю камеру на телефонах
        {
          fps: 10,
          qrbox: { width: 250, height: 250 }
        },
        (decodedText) => {
          // Успешное сканирование
          stopCamera();
          processCode(decodedText);
        },
        (errorMessage) => {
          // Игнорируем фоновые ошибки поиска кода (они сыпятся постоянно, пока кода нет в кадре)
        }
    );
  } catch (err) {
    console.error("Ошибка запуска камеры", err);
    ElMessage.warning('Не удалось запустить камеру. Проверьте разрешения или используйте загрузку картинки.');
    qrMode.value = 'file'; // Перекидываем на загрузку файла, если нет камеры
  }
}

const stopCamera = async () => {
  if (html5QrCode && html5QrCode.isScanning) {
    try {
      await html5QrCode.stop();
      html5QrCode.clear();
    } catch (err) {
      console.error("Ошибка остановки камеры", err);
    }
  }
}

// === ЛОГИКА ФАЙЛА ===
const handleFileUpload = async (uploadFile) => {
  const file = uploadFile.raw;
  if (!file) return;

  const html5QrCodeForFile = new Html5Qrcode("qr-reader"); // Используем тот же элемент (он скрыт, но API требует ID)

  try {
    const decodedText = await html5QrCodeForFile.scanFile(file, true);
    processCode(decodedText);
  } catch (err) {
    ElMessage.error('Не удалось найти QR-код на изображении. Попробуйте другую картинку.');
    console.error(err);
  }
}

// === УПРАВЛЕНИЕ ЖИЗНЕННЫМ ЦИКЛОМ ===
const switchQrMode = async (mode) => {
  if (mode === 'camera') {
    await stopCamera(); // на всякий случай останавливаем предыдущий
    await nextTick();
    startCamera();
  } else {
    await stopCamera();
  }
};

const handleTabChange = async (tabName) => {
  if (tabName === 'qr' && qrMode.value === 'camera') {
    await nextTick();
    startCamera();
  } else {
    await stopCamera();
  }
}

onMounted(() => {
  // Запускаем камеру по умолчанию при монтировании
  startCamera();
})

onBeforeUnmount(() => {
  stopCamera();
})

// === ОТПРАВКА ДАННЫХ ===
const processCode = async (code) => {
  const eventId = route.query.eventId;
  if (!eventId) {
    ElMessage.error('Не указан ID мероприятия в ссылке');
    return;
  }

  isLoading.value = true;
  try {
    // Предполагаем, что eventStore.confirmParticipation возвращает данные с полем points_earned
    const result = await eventStore.confirmParticipation(eventId, code);
    const points = result?.points_earned ?? result?.data?.points_earned;
    if (points) {
      // ElMessage.success(`Участие подтверждено! Начислено ${points} баллов.`);
    } else {
      // ElMessage.success('Участие успешно подтверждено!');
    }
    router.push('/profile');
  } catch (err) {
    console.error(err);
    let errorMsg = 'Ошибка подтверждения.';
    if (err.response?.data?.message) {
      errorMsg = err.response.data.message;
    }
    ElMessage.error(errorMsg);
    // Если была активна камера, перезапускаем её
    if (activeTab.value === 'qr' && qrMode.value === 'camera') {
      setTimeout(startCamera, 2000);
    }
  } finally {
    isLoading.value = false;
  }
};

const confirmByCode = () => {
  if (!manualCode.value) return;
  processCode(manualCode.value);
}
</script>

<style scoped>
/* Центрирование и ограничение ширины для красоты */
.confirm-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 100px); /* Подгоните под ваш layout */
  padding: 20px;
}

.confirm-card {
  width: 100%;
  max-width: 480px;
  border-radius: var(--radius-lg);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-header h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-main);
}

.header-icon {
  font-size: 1.5rem;
  color: var(--primary);
}

.qr-mode-switch {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

/* Стили для контейнера камеры */
.scanner-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.qr-reader {
  width: 100%;
  max-width: 350px;
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 2px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  margin-bottom: 12px;
}

/* Скрытие дефолтных страшных ссылок html5-qrcode */
:deep(#qr-reader__dashboard_section_csr span),
:deep(#qr-reader__dashboard_section_swaplink) {
  display: none !important;
}

.hint-text {
  color: var(--text-muted);
  font-size: 0.875rem;
  text-align: center;
}

/* Стили для загрузчика файла */
.upload-wrapper {
  padding: 10px 0;
}

.qr-uploader :deep(.el-upload-dragger) {
  background-color: var(--bg-app);
  border-color: var(--border-color);
  transition: all 0.3s;
}

.qr-uploader :deep(.el-upload-dragger:hover) {
  border-color: var(--primary);
  background-color: #eff6ff;
}

.text-center {
  text-align: center;
}

/* Стили для ручного ввода */
.manual-wrapper {
  padding: 20px 0 10px;
}

.full-width-btn {
  width: 100%;
  margin-top: 10px;
}
</style>