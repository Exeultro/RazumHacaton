<template>
    <el-form :inline="true" class="filters">
      <el-form-item label="Поиск">
        <el-input v-model="filters.search" placeholder="Название" clearable />
      </el-form-item>
      <el-form-item label="Направление">
        <el-select v-model="filters.direction" placeholder="Все" clearable>
          <el-option label="IT" value="IT" />
          <el-option label="Социальное проектирование" value="Социальное проектирование" />
          <el-option label="Медиа" value="Медиа" />
        </el-select>
      </el-form-item>
      <el-form-item label="Дата с">
        <el-date-picker v-model="filters.dateFrom" type="date" placeholder="Выберите дату" />
      </el-form-item>
      <el-form-item label="Дата по">
        <el-date-picker v-model="filters.dateTo" type="date" placeholder="Выберите дату" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="applyFilters">Применить</el-button>
        <el-button @click="resetFilters">Сбросить</el-button>
      </el-form-item>
    </el-form>
  </template>
  
  <script setup>
  import { reactive } from 'vue'
  import { useEventStore } from '../stores/eventStore'
  
  const eventStore = useEventStore()
  const filters = reactive({
    search: '',
    direction: '',
    dateFrom: null,
    dateTo: null
  })
  
  const applyFilters = () => {
    eventStore.updateFilters(filters)
  }
  const resetFilters = () => {
    filters.search = ''
    filters.direction = ''
    filters.dateFrom = null
    filters.dateTo = null
    applyFilters()
  }
  </script>