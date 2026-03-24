<template>
  <div class="card event-card" @click="goToDetail">
    <h3>{{ event.title }}</h3>
    <p><strong>Дата:</strong> {{ formatDate(event.event_date) }}</p>
    <p><strong>Направление:</strong> {{ event.direction }}</p>
    <p><strong>Баллы:</strong> {{ event.points }}</p>
    <p><strong>Призы:</strong> {{ formatPrizes(event.prizes) }}</p>
    <span class="badge badge-primary">Организатор: {{ event.organizer.full_name }}</span>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps(['event'])
const router = useRouter()

const goToDetail = () => {
  router.push(`/events/${props.event.id}`)
}
const formatDate = (iso) => new Date(iso).toLocaleDateString()
const formatPrizes = (prizes) => prizes?.map(p => p.name).join(', ') || '—'
</script>