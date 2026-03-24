<template>
  <div class="table-container">
    <div class="table-wrapper">
      <table class="candidates-table">
        <thead>
        <tr>
          <th class="rank-col">#</th>
          <th class="name-col">ФИО</th>
          <th class="age-col">Возраст</th>
          <th class="city-col">Город</th>
          <th class="direction-col">Направление</th>
          <th class="events-col">Мероприятий</th>
          <th class="points-col">Баллы</th>
          <th class="avg-col">Ср. балл</th>
          <th class="actions-col"></th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="candidate in candidates" :key="candidate.id">
          <td class="rank-col">
            <span class="rank-number">{{ candidate.global_rank }}</span>
          </td>
          <td class="name-col">
            <div class="user-info">
              <div class="user-avatar">
                {{ getInitials(candidate.full_name) }}
              </div>
              <div>
                <div class="user-name">{{ candidate.full_name }}</div>
              </div>
            </div>
          </td>
          <td class="age-col">{{ candidate.age }} лет</td>
          <td class="city-col">{{ candidate.city }}</td>
          <td class="direction-col">
            <span class="direction-tag" :class="getDirectionClass(candidate.direction)">
              {{ getDirectionLabel(candidate.direction) }}
            </span>
          </td>
          <td class="events-col">
            <span class="events-badge">{{ candidate.events_count }}</span>
          </td>
          <td class="points-col">
            <span class="points-value">{{ formatPoints(candidate.total_points) }}</span>
          </td>
          <td class="avg-col">
            <span class="avg-value">{{ getAveragePoints(candidate) }}</span>
          </td>
          <td class="actions-col">
            <button @click="exportCandidate(candidate)" class="action-btn" title="Скачать PDF">
              📄
            </button>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <div v-if="!candidates.length" class="empty-state">
      <span class="empty-icon">🔍</span>
      <p>Кандидаты не найдены</p>
      <p class="empty-hint">Попробуйте изменить параметры фильтрации</p>
    </div>

    <!-- Модальное окно -->
    <div v-if="selectedCandidate" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2>Профиль кандидата</h2>
          <button class="close-btn" @click="closeModal">×</button>
        </div>
        <div class="modal-content">
          <div class="profile-avatar">
            {{ getInitials(selectedCandidate.full_name) }}
          </div>
          <h3 class="profile-name">{{ selectedCandidate.full_name }}</h3>
          <div class="profile-info">
            <p><strong>🎂 Возраст:</strong> {{ selectedCandidate.age }} лет</p>
            <p><strong>🏙️ Город:</strong> {{ selectedCandidate.city }}</p>
            <p><strong>📌 Направление:</strong> {{ getDirectionLabel(selectedCandidate.direction) }}</p>
          </div>
          <div class="profile-stats">
            <div class="stat">
              <div class="stat-value">{{ selectedCandidate.total_points }}</div>
              <div class="stat-label">Всего баллов</div>
            </div>
            <div class="stat">
              <div class="stat-value">{{ selectedCandidate.events_count }}</div>
              <div class="stat-label">Мероприятий</div>
            </div>
            <div class="stat">
              <div class="stat-value">#{{ selectedCandidate.rank }}</div>
              <div class="stat-label">Место в рейтинге</div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="exportCandidate(selectedCandidate)" class="btn-pdf">
            📄 Скачать отчет PDF
          </button>
          <button @click="closeModal" class="btn-close">Закрыть</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useUiStore } from '@/stores/ui'
import html2pdf from 'html2pdf.js'

const props = defineProps({
  candidates: {
    type: Array,
    default: () => []
  }
})

const uiStore = useUiStore()
const selectedCandidate = ref(null)

const formatPoints = (points) => {
  if (!points) return '0'
  return points.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
}

const getAveragePoints = (candidate) => {
  if (candidate.average_points_per_event) return candidate.average_points_per_event
  if (candidate.events_count && candidate.total_points) {
    return (candidate.total_points / candidate.events_count).toFixed(1)
  }
  return 0
}

const getInitials = (name) => {
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

const getDirectionLabel = (direction) => {
  const labels = {
    'IT': '💻 IT',
    'social': '🤝 Социальное',
    'media': '📱 Медиа'
  }
  return labels[direction] || direction
}

const getDirectionClass = (direction) => {
  const classes = {
    'IT': 'it',
    'social': 'social',
    'media': 'media'
  }
  return classes[direction] || ''
}

const openModal = (candidate) => {
  selectedCandidate.value = candidate
}

const closeModal = () => {
  selectedCandidate.value = null
}

const formatDate = (date) => {
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  }).format(date)
}

const generateSinglePdf = async (candidate) => {
  const avg = getAveragePoints(candidate)
  const age = candidate.age !== undefined ? `${candidate.age} лет` : '—'
  const city = candidate.city || '—'
  const direction = getDirectionLabel(candidate.direction) || '—'
  const totalPoints = candidate.total_points !== undefined ? candidate.total_points.toLocaleString() : '—'
  const eventsCount = candidate.events_count !== undefined ? candidate.events_count : '—'
  const rank = candidate.rank !== undefined ? `#${candidate.rank}` : '—'

  const html = `
    <!DOCTYPE html>
    <html>
    <head>
      <meta charset="UTF-8">
      <title>Отчет: ${candidate.full_name}</title>
      <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
          font-family: 'Inter', 'Segoe UI', sans-serif;
          padding: 40px;
          background: white;
        }
        .report {
          max-width: 800px;
          margin: 0 auto;
        }
        .header {
          text-align: center;
          border-bottom: 2px solid #1d4ed8;
          padding-bottom: 20px;
          margin-bottom: 30px;
        }
        .header h1 {
          color: #1d4ed8;
          margin-bottom: 8px;
        }
        .date {
          color: #64748b;
          font-size: 12px;
        }
        h2 {
          color: #0f172a;
          font-size: 20px;
          border-left: 4px solid #1d4ed8;
          padding-left: 12px;
          margin: 20px 0 16px;
        }
        .info-table {
          width: 100%;
          border-collapse: collapse;
          margin-bottom: 20px;
        }
        .info-table td {
          padding: 10px 0;
          border-bottom: 1px solid #e2e8f0;
        }
        .label {
          width: 120px;
          color: #64748b;
          font-weight: 500;
        }
        .value {
          color: #0f172a;
          font-weight: 500;
        }
        .stats {
          display: grid;
          grid-template-columns: repeat(4, 1fr);
          gap: 16px;
          margin: 20px 0;
        }
        .stat-card {
          text-align: center;
          padding: 16px;
          background: #f8fafc;
          border-radius: 12px;
        }
        .stat-value {
          font-size: 28px;
          font-weight: 700;
          color: #1d4ed8;
        }
        .stat-label {
          font-size: 12px;
          color: #64748b;
          margin-top: 4px;
        }
        .recommendation {
          background: #f1f5f9;
          padding: 16px;
          border-radius: 12px;
          margin-top: 20px;
        }
        .footer {
          text-align: center;
          font-size: 10px;
          color: #94a3b8;
          margin-top: 40px;
          padding-top: 16px;
          border-top: 1px solid #e2e8f0;
        }
      </style>
    </head>
    <body>
      <div class="report">
        <div class="header">
          <h1>Отчет по кандидату</h1>
          <div class="date">${formatDate(new Date())}</div>
        </div>

        <h2>${candidate.full_name}</h2>

        <table class="info-table">
          <tr><td class="label">Возраст:</td><td class="value">${age}</td></tr>
          <tr><td class="label">Город:</td><td class="value">${city}</td></tr>
          <tr><td class="label">Направление:</td><td class="value">${direction}</td></tr>
        </table>

        <div class="stats">
          <div class="stat-card">
            <div class="stat-value">${totalPoints}</div>
            <div class="stat-label">Всего баллов</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">${eventsCount}</div>
            <div class="stat-label">Мероприятий</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">${avg}</div>
            <div class="stat-label">Средний балл</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">${rank}</div>
            <div class="stat-label">Место в рейтинге</div>
          </div>
        </div>

        <div class="recommendation">
          <strong>Заключение:</strong><br>
          ${candidate.total_points > 2000
      ? `Кандидат ${candidate.full_name} показывает отличные результаты. Рекомендуется к включению в кадровый резерв высшего уровня.`
      : candidate.total_points > 1000
          ? `Кандидат ${candidate.full_name} показывает стабильные результаты. Рекомендуется к рассмотрению на стажировку.`
          : `Кандидат ${candidate.full_name} имеет потенциал для развития. Рекомендуется дополнительное вовлечение в мероприятия.`
  }
        </div>

        <div class="footer">
          Платформа активности молодежного парламента
        </div>
      </div>
    </body>
    </html>
  `

  const opt = {
    margin: [10, 10, 10, 10],
    filename: `candidate_${candidate.full_name.replace(/\s/g, '_')}_${Date.now()}.pdf`,
    image: { type: 'jpeg', quality: 0.98 },
    html2canvas: { scale: 2 },
    jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' }
  }

  await html2pdf().set(opt).from(html).save()
}

const exportCandidate = async (candidate) => {
  uiStore.addNotification({
    type: 'info',
    message: `Генерация отчета для ${candidate.full_name}...`
  })

  try {
    await generateSinglePdf(candidate)
    uiStore.addNotification({
      type: 'success',
      message: `Отчет для ${candidate.full_name} скачан`
    })
  } catch (error) {
    console.error('Ошибка:', error)
    uiStore.addNotification({
      type: 'error',
      message: `Ошибка при генерации отчета`
    })
  }
}
</script>

<style scoped>
/* стили без изменений, но можно оставить как есть */
.table-container {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}
.table-wrapper {
  overflow-x: auto;
}
.candidates-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}
.candidates-table th {
  text-align: left;
  padding: 14px 16px;
  background: #f8fafc;
  font-weight: 600;
  color: #475569;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #e2e8f0;
}
.candidates-table td {
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  vertical-align: middle;
}
.candidates-table tr:hover {
  background: #f8fafc;
}
.rank-col {
  width: 50px;
  text-align: center;
}
.rank-number {
  font-weight: 600;
  color: #64748b;
}
.name-col {
  min-width: 200px;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1d4ed8 0%, #7c3aed 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}
.user-name {
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 2px;
}
.age-col, .city-col {
  color: #334155;
}
.direction-tag {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 500;
}
.direction-tag.it { background: #eef2ff; color: #1d4ed8; }
.direction-tag.social { background: #ecfdf5; color: #10b981; }
.direction-tag.media { background: #fef3c7; color: #f59e0b; }
.events-badge {
  display: inline-block;
  padding: 4px 8px;
  background: #f1f5f9;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: #475569;
}
.points-value {
  font-weight: 700;
  color: #1d4ed8;
  font-size: 15px;
}
.avg-value {
  font-weight: 500;
  color: #10b981;
}
.actions-col {
  width: 70px;
  text-align: center;
  white-space: nowrap;
}
.action-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 6px;
  margin: 0 2px;
  border-radius: 6px;
  transition: all 0.2s;
  opacity: 0.6;
}
.action-btn:hover {
  opacity: 1;
  background: #f1f5f9;
  transform: scale(1.1);
}
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #94a3b8;
}
.empty-icon {
  font-size: 48px;
  display: block;
  margin-bottom: 16px;
  opacity: 0.5;
}
.empty-hint {
  font-size: 12px;
  margin-top: 8px;
  opacity: 0.7;
}
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
}
.modal-header h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #0f172a;
}
.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #94a3b8;
  padding: 0;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  transition: all 0.2s;
}
.close-btn:hover {
  background: #f1f5f9;
  color: #475569;
}
.modal-content {
  padding: 24px;
}
.profile-avatar {
  width: 80px;
  height: 80px;
  margin: 0 auto 16px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1d4ed8 0%, #7c3aed 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  font-weight: 600;
}
.profile-name {
  text-align: center;
  font-size: 1.25rem;
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 20px;
}
.profile-info {
  background: #f8fafc;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 20px;
}
.profile-info p {
  margin: 8px 0;
  font-size: 14px;
}
.profile-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  text-align: center;
}
.stat {
  padding: 12px;
  background: #f1f5f9;
  border-radius: 12px;
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1d4ed8;
}
.stat-label {
  font-size: 11px;
  color: #64748b;
  margin-top: 4px;
}
.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}
.btn-pdf {
  padding: 10px 20px;
  background: #10b981;
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-pdf:hover {
  background: #059669;
}
.btn-close {
  padding: 10px 20px;
  background: #f1f5f9;
  color: #475569;
  border: none;
  border-radius: 10px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-close:hover {
  background: #e2e8f0;
}

@media (max-width: 768px) {
  .candidates-table th,
  .candidates-table td {
    padding: 10px 8px;
    font-size: 13px;
  }

  .user-avatar {
    width: 32px;
    height: 32px;
    font-size: 12px;
  }

  .user-name {
    font-size: 13px;
  }

  .direction-tag {
    font-size: 10px;
    padding: 2px 6px;
  }

  .events-badge {
    font-size: 12px;
    padding: 2px 6px;
  }

  .points-value {
    font-size: 13px;
  }

  .action-btn {
    font-size: 16px;
    padding: 4px;
  }

  .modal {
    width: 95%;
    margin: 16px;
  }
}

@media (max-width: 480px) {
  .candidates-table th,
  .candidates-table td {
    padding: 8px 6px;
    font-size: 12px;
  }

  .user-avatar {
    width: 28px;
    height: 28px;
    font-size: 10px;
  }

  .profile-avatar {
    width: 64px;
    height: 64px;
    font-size: 24px;
  }

  .stat-value {
    font-size: 20px;
  }

  .modal-footer {
    flex-direction: column;
    gap: 8px;
  }

  .btn-pdf, .btn-close {
    width: 100%;
    text-align: center;
  }
}
</style>