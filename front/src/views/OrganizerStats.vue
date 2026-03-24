<template>
    <div class="organizer-stats">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="Проведено мероприятий">{{ organizer.eventsCount }}</el-descriptions-item>
        <el-descriptions-item label="Рейтинг доверия">
          <el-rate :value="organizer.rating" disabled show-score />
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </template>
  
  <script setup>
  defineProps(['organizer'])
  </script>