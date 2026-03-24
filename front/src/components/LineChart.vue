<template>
    <canvas ref="chartCanvas"></canvas>
  </template>
  
  <script setup>
  import { ref, onMounted, watch } from 'vue'
  import { Chart, registerables } from 'chart.js'
  Chart.register(...registerables)
  
  const props = defineProps(['data', 'options'])
  const chartCanvas = ref(null)
  let chartInstance = null
  
  onMounted(() => {
    if (chartCanvas.value) {
      chartInstance = new Chart(chartCanvas.value, {
        type: 'line',
        data: props.data,
        options: props.options || { responsive: true }
      })
    }
  })
  
  watch(() => props.data, (newData) => {
    if (chartInstance) {
      chartInstance.data = newData
      chartInstance.update()
    }
  })
  </script>