<template>
  <div class="table-container">
    <div class="table-wrapper">
      <table class="leaderboard-table">
        <thead>
        <tr>
          <th class="rank-col">Место</th>
          <th class="user-col">Участник</th>
          <th class="direction-col">Направление</th>
          <th class="points-col" @click="toggleSort">
            Баллы
            <span class="sort-icon">{{ sortOrder === 'desc' ? '↓' : '↑' }}</span>
          </th>
        </tr>
        </thead>
        <tbody>
        <tr
            v-for="leader in sortedLeaders"
            :key="leader.user_id"
            :class="{ 'current-user': isCurrentUser(leader) }"
        >
          <td class="rank-col">
            <div class="rank-badge" :class="getRankClass(leader.global_rank)">
              {{ leader.global_rank }}
            </div>
          </td>
          <td class="user-col">
            <div class="user-info">
              <AppAvatar :url="leader.avatar_url" :name="leader.full_name" size="sm" />
              <span class="user-name">{{ leader.full_name || 'Участник' }}</span>
            </div>
          </td>
          <td class="direction-col">
              <span v-if="leader.direction" class="direction-tag" :class="getDirectionClass(leader.direction)">
                {{ getDirectionLabel(leader.direction) }}
              </span>
            <span v-else class="direction-tag">—</span>
          </td>
          <td class="points-col">
            <span class="points-value">{{ formatPoints(leader.total_points) }}</span>
          </td>
        </tr>
        </tbody>
      </table>
    </div>

    <div v-if="userRank" class="user-rank-card">
      <div class="user-rank-info">
        <span class="rank-label">Ваше место:</span>
        <span class="rank-number">#{{ userRank.global_rank }}</span>
        <span class="points-label">Баллов:</span>
        <span class="points-number">{{ formatPoints(userRank.total_points) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AppAvatar from "@/components/ui/AppAvatar.vue";

const props = defineProps({
  leaders: {
    type: Array,
    default: () => []
  },
  userRank: {
    type: Object,
    default: null
  }
})

const authStore = useAuthStore()
const sortOrder = ref('desc')

const sortedLeaders = computed(() => {
  if (sortOrder.value === 'desc') {
    return [...props.leaders]
  }
  return [...props.leaders].reverse()
})

const toggleSort = () => {
  sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc'
}

const isCurrentUser = (leader) => {
  return authStore.user?.id === leader.user_id
}

const formatPoints = (points) => {
  if (!points && points !== 0) return '0'
  return points.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
}

const getRankClass = (rank) => {
  if (rank === 1) return 'gold'
  if (rank === 2) return 'silver'
  if (rank === 3) return 'bronze'
  return ''
}

const getDirectionLabel = (direction) => {
  const labels = {
    'IT': 'IT',
    'social': 'Социальное',
    'media': 'Медиа'
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
</script>

<style scoped>
.table-container {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
}

.table-wrapper {
  overflow-x: auto;
}

.leaderboard-table {
  width: 100%;
  border-collapse: collapse;
}

.leaderboard-table th {
  text-align: left;
  padding: 16px 20px;
  background: #f8fafc;
  font-weight: 600;
  color: #475569;
  font-size: 0.8125rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #e2e8f0;
}

.leaderboard-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.875rem;
}

.leaderboard-table tr:hover {
  background: #f8fafc;
}

.leaderboard-table tr.current-user {
  background: linear-gradient(90deg, rgba(29, 78, 216, 0.05) 0%, rgba(124, 58, 237, 0.05) 100%);
  border-left: 3px solid #1d4ed8;
}

.points-col {
  cursor: pointer;
  user-select: none;
}

.sort-icon {
  display: inline-block;
  margin-left: 6px;
  font-size: 0.75rem;
  color: #94a3b8;
}

.rank-col {
  width: 80px;
}

.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  font-weight: 700;
  background: #f1f5f9;
  color: #475569;
}

.rank-badge.gold {
  background: #fbbf24;
  color: #0f172a;
  box-shadow: 0 2px 6px rgba(251, 191, 36, 0.3);
}

.rank-badge.silver {
  background: #94a3b8;
  color: white;
}

.rank-badge.bronze {
  background: #d97706;
  color: white;
}

.user-col {
  min-width: 220px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name {
  font-weight: 500;
  color: #0f172a;
}

.direction-tag {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 500;
}

.direction-tag.it {
  background: #eef2ff;
  color: #1d4ed8;
}

.direction-tag.social {
  background: #ecfdf5;
  color: #10b981;
}

.direction-tag.media {
  background: #fef3c7;
  color: #f59e0b;
}

.points-value {
  font-weight: 700;
  color: #1d4ed8;
  font-size: 1rem;
}

.user-rank-card {
  padding: 20px 24px;
  background: linear-gradient(135deg, #1d4ed8 0%, #7c3aed 100%);
  margin: 20px;
  border-radius: 12px;
  color: white;
}

.user-rank-info {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;
}

.rank-label, .points-label {
  font-size: 0.875rem;
  opacity: 0.9;
}

.rank-number {
  font-size: 1.75rem;
  font-weight: 700;
}

.points-number {
  font-size: 1.5rem;
  font-weight: 700;
}
</style>