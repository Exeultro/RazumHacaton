<template>
    <div>
      <el-table :data="adminStore.pendingOrganizers" style="width: 100%">
        <el-table-column prop="name" label="Имя" />
        <el-table-column prop="eventsCount" label="Проведено мероприятий" />
        <el-table-column label="Действия">
          <template #default="{ row }">
            <el-button type="success" @click="approve(row.id)">Одобрить</el-button>
            <el-button type="danger" @click="reject(row.id)">Отклонить</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </template>
  
  <script setup>
  import { useAdminStore } from '../../stores/adminStore'
  
  const adminStore = useAdminStore()
  
  const approve = (id) => {
    adminStore.approveOrganizer(id)
  }
  const reject = (id) => {
    adminStore.rejectOrganizer(id)
  }
  </script>