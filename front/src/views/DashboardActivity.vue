<template>
    <div class="dashboard">
      <h1>Активность</h1>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <h3>Последние мероприятия</h3>
            <ul>
              <li v-for="ev in dashboardStore.recentEvents" :key="ev.id">
                <router-link :to="`/events/${ev.id}`">{{ ev.title }}</router-link> – {{ ev.organizer_name }}, {{ formatDate(ev.event_date) }}, участников: {{ ev.participants_count }}
              </li>
            </ul>
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card>
            <h3>Облако тегов</h3>
            <div class="tags">
              <el-tag v-for="tag in dashboardStore.trendingTags" :key="tag.tag" size="large" style="margin: 4px">
                {{ tag.tag }} ({{ tag.count }})
              </el-tag>
            </div>
          </el-card>
        </el-col>
      </el-row>
      <el-card style="margin-top: 20px">
        <h3>Динамика рейтинга</h3>
        <LineChart :data="chartData" />
      </el-card>
    </div>
  </template>
  
  <script setup>
  import { onMounted, computed } from 'vue'
  import { useDashboardStore } from '../stores/dashboardStore'
  import LineChart from '../components/LineChart.vue' // компонент с Chart.js
  
  const dashboardStore = useDashboardStore()
  
  onMounted(() => {
    dashboardStore.fetchActivity()
  })
  
  const formatDate = (iso) => new Date(iso).toLocaleDateString()
  
  const chartData = computed(() => ({
    labels: dashboardStore.ratingHistory.map(h => h.date),
    datasets: [{ label: 'Баллы', data: dashboardStore.ratingHistory.map(h => h.points) }]
  }))
  </script>