<template>
  <div>
    <el-table :data="adminStore.eventTypeWeights" style="width: 100%">
      <el-table-column prop="name" label="Тип мероприятия" />
      <el-table-column label="Вес">
        <template #default="{ row }">
          <el-input-number v-model="row.weight" :min="0" :step="0.1" @change="updateWeight(row.id, row.weight)" />
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { useAdminStore } from '../../stores/adminStore'

const adminStore = useAdminStore()

const updateWeight = (id, weight) => {
  adminStore.updateWeight(id, weight)
}
</script>