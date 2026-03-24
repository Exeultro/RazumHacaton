<template>
  <div class="chart-container">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  }
})

const chartCanvas = ref(null)
let chartInstance = null

const createChart = () => {
  if (!chartCanvas.value || !props.data.length) return
  
  if (chartInstance) {
    chartInstance.destroy()
  }
  
  const ctx = chartCanvas.value.getContext('2d')
  
  chartInstance = new Chart(ctx, {
    type: 'line',
    data: {
      labels: props.data.map(item => {
        const date = new Date(item.date)
        return `${date.getDate()}/${date.getMonth() + 1}`
      }),
      datasets: [
        {
          label: 'Рейтинг',
          data: props.data.map(item => item.points),
          borderColor: '#1d4ed8',
          backgroundColor: 'rgba(29, 78, 216, 0.05)',
          borderWidth: 3,
          pointBackgroundColor: '#1d4ed8',
          pointBorderColor: 'white',
          pointRadius: 4,
          pointHoverRadius: 6,
          tension: 0.3,
          fill: true
        }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        tooltip: {
          callbacks: {
            label: (context) => `Баллы: ${context.raw}`
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          grid: {
            color: '#e2e8f0'
          },
          title: {
            display: true,
            text: 'Количество баллов',
            color: '#64748b'
          }
        },
        x: {
          grid: {
            display: false
          },
          title: {
            display: true,
            text: 'Дата',
            color: '#64748b'
          }
        }
      }
    }
  })
}

onMounted(() => {
  createChart()
})

watch(() => props.data, () => {
  createChart()
}, { deep: true })
</script>

<style scoped>
.chart-container {
  width: 100%;
  height: 350px;
  position: relative;
}
</style>