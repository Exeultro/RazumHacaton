import { defineStore } from 'pinia';
import { ref } from 'vue';
import api from '../api/axios'; // Единый импорт API
import { useUiStore } from './ui';
import {ElMessage} from "element-plus";

export const useAnalyticsStore = defineStore('analytics', () => {
  const recentEvents = ref([])
  const ratingHistory = ref([])
  const trendingTags = ref([])
  const isLoading = ref(false)

  const uiStore = useUiStore()

  const fetchActivity = async () => {
    isLoading.value = true

    try {
      const response = await api.get('/dashboard')
      recentEvents.value = response.data.recent_events
      ratingHistory.value = response.data.rating_history
      trendingTags.value = response.data.trending_tags
    } catch (err) {
      ElMessage.error(err.response?.data?.error?.message || 'Ошибка загрузки активности');
    } finally {
      isLoading.value = false
    }
  }

  const generateMockData = () => {
    // График
    ratingHistory.value = []
    let points = 100
    for (let i = 30; i >= 0; i--) {
      const date = new Date()
      date.setDate(date.getDate() - i)
      points += Math.floor(Math.random() * 50) + 10
      ratingHistory.value.push({
        date: date.toISOString().split('T')[0],
        points: points
      })
    }

    // Облако тегов
    trendingTags.value = [
      { tag: 'IT', count: 156 },
      { tag: 'Хакатон', count: 98 },
      { tag: 'Социальное проектирование', count: 87 },
      { tag: 'Медиа', count: 76 },
      { tag: 'Волонтерство', count: 65 },
      { tag: 'Мастер-класс', count: 54 },
      { tag: 'Лекция', count: 43 },
      { tag: 'Спорт', count: 32 }
    ]

    // Лента
    recentEvents.value = [
      { id: 1, title: 'Хакатон по AI', event_date: new Date().toISOString(), organizer_name: 'Петр Петров', participants_count: 45 },
      { id: 2, title: 'Лекция про карьеру', event_date: new Date(Date.now() - 86400000).toISOString(), organizer_name: 'Анна Смирнова', participants_count: 78 },
      { id: 3, title: 'Спортивный фестиваль', event_date: new Date(Date.now() - 172800000).toISOString(), organizer_name: 'Иван Иванов', participants_count: 120 }
    ]
  }

  return {
    recentEvents,
    ratingHistory,
    trendingTags,
    isLoading,
    fetchActivity
  }
});