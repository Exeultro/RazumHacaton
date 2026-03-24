<template>
  <div class="export-container">
    <button
        @click="exportToPdf"
        :disabled="isGenerating || !candidates.length"
        class="btn-export"
    >
      <span v-if="isGenerating" class="spinner-small"></span>
      <span v-else>📄</span>
      {{ isGenerating ? 'Генерация...' : 'Экспорт всех в PDF' }}
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useUiStore } from '@/stores/ui'

const props = defineProps({
  candidates: {
    type: Array,
    default: () => []
  }
})

const uiStore = useUiStore()
const isGenerating = ref(false)

const formatDate = (date) => {
  return new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  }).format(date)
}

const getDirectionLabel = (direction) => {
  const labels = {
    'IT': 'IT',
    'social': 'Социальное проектирование',
    'media': 'Медиа'
  }
  return labels[direction] || direction
}

const getAveragePoints = (candidate) => {
  if (candidate.average_points_per_event) return candidate.average_points_per_event
  if (candidate.events_count && candidate.total_points) {
    return (candidate.total_points / candidate.events_count).toFixed(1)
  }
  return 0
}

const getRecommendation = (candidate) => {
  const points = candidate.total_points || 0
  if (points > 2000) {
    return `Кандидат ${candidate.full_name} показывает отличные результаты. Рекомендуется к включению в кадровый резерв высшего уровня.`
  } else if (points > 1000) {
    return `Кандидат ${candidate.full_name} показывает стабильные результаты. Рекомендуется к рассмотрению на стажировку.`
  } else {
    return `Кандидат ${candidate.full_name} имеет потенциал для развития. Рекомендуется дополнительное вовлечение в мероприятия.`
  }
}

const exportToPdf = async () => {
  if (!props.candidates.length) {
    uiStore.addNotification({
      type: 'warning',
      message: 'Нет кандидатов для экспорта'
    })
    return
  }

  isGenerating.value = true

  uiStore.addNotification({
    type: 'info',
    message: 'Генерация PDF отчета...'
  })

  try {
    const printWindow = window.open('', '_blank')
    const today = new Date()

    let html = `
      <!DOCTYPE html>
      <html>
      <head>
        <meta charset="UTF-8">
        <title>Отчет по кандидатам</title>
        <style>
          * { margin: 0; padding: 0; box-sizing: border-box; }
          body {
            font-family: 'Segoe UI', 'Inter', sans-serif;
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
          .candidate {
            margin-bottom: 40px;
            page-break-after: always;
          }
          .candidate:last-child {
            page-break-after: auto;
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
          @media print {
            body { padding: 0; }
            .candidate { page-break-after: always; }
          }
        </style>
      </head>
      <body>
        <div class="report">
          <div class="header">
            <h1>Отчет по кандидатам</h1>
            <div class="date">${formatDate(today)}</div>
          </div>
    `

    props.candidates.forEach(candidate => {
      const age = candidate.age !== undefined ? `${candidate.age} лет` : '—'
      const city = candidate.city || '—'
      const direction = getDirectionLabel(candidate.direction) || '—'
      const totalPoints = candidate.total_points !== undefined ? candidate.total_points.toLocaleString() : '—'
      const eventsCount = candidate.events_count !== undefined ? candidate.events_count : '—'
      const avgPoints = getAveragePoints(candidate)
      const rank = candidate.rank !== undefined ? `#${candidate.rank}` : '—'

      html += `
        <div class="candidate">
          <h2>${candidate.full_name || '—'}</h2>

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
              <div class="stat-value">${avgPoints}</div>
              <div class="stat-label">Средний балл</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">${rank}</div>
              <div class="stat-label">Место в рейтинге</div>
            </div>
          </div>

          <div class="recommendation">
            <strong>Заключение:</strong><br>
            ${getRecommendation(candidate)}
          </div>
        </div>
      `
    })

    html += `
          <div class="footer">
            Платформа активности молодежного парламента
          </div>
        </div>
      </body>
      </html>
    `

    printWindow.document.write(html)
    printWindow.document.close()

    setTimeout(() => {
      printWindow.print()
      printWindow.close()
    }, 500)

    uiStore.addNotification({
      type: 'success',
      message: `PDF отчет успешно сгенерирован (${props.candidates.length} кандидатов)`
    })
  } catch (error) {
    console.error('Ошибка генерации PDF:', error)
    uiStore.addNotification({
      type: 'error',
      message: 'Ошибка при генерации PDF'
    })
  } finally {
    isGenerating.value = false
  }
}
</script>

<style scoped>
.export-container {
  display: inline-block;
}

.btn-export {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}

.btn-export:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.btn-export:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  display: inline-block;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .export-container {
    width: 100%;
  }

  .btn-export {
    width: 100%;
    justify-content: center;
  }
}
</style>