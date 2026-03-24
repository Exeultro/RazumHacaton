<template>
    <el-form :inline="true">
      <el-form-item label="Поиск">
        <el-input v-model="localFilters.search" placeholder="Название" clearable @clear="emitFilters" />
      </el-form-item>
      <el-form-item label="Направление">
        <el-select v-model="localFilters.direction" placeholder="Все" clearable @change="emitFilters">
          <el-option label="IT" value="IT" />
          <el-option label="Социальное" value="social" />
          <el-option label="Медиа" value="media" />
        </el-select>
      </el-form-item>
      <el-form-item label="Формат">
        <el-select v-model="localFilters.format" placeholder="Любой" clearable @change="emitFilters">
          <el-option label="Офлайн" value="offline" />
          <el-option label="Онлайн" value="online" />
          <el-option label="Гибрид" value="hybrid" />
        </el-select>
      </el-form-item>
      <el-form-item label="Дата с">
        <el-date-picker v-model="localFilters.date_from" type="date" @change="emitFilters" />
      </el-form-item>
      <el-form-item label="Дата по">
        <el-date-picker v-model="localFilters.date_to" type="date" @change="emitFilters" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="emitFilters">Применить</el-button>
        <el-button @click="resetFilters">Сбросить</el-button>
      </el-form-item>
    </el-form>
  </template>
  
  <script setup>
  import { reactive } from 'vue'
  
  const emit = defineEmits(['filter'])
  const localFilters = reactive({
    search: '',
    direction: null,
    format: null,
    date_from: null,
    date_to: null
  })
  
  const emitFilters = () => {
    emit('filter', { ...localFilters })
  }
  
  const resetFilters = () => {
    localFilters.search = ''
    localFilters.direction = null
    localFilters.format = null
    localFilters.date_from = null
    localFilters.date_to = null
    emitFilters()
  }
  </script>